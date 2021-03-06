// Copyright 2019 Yaacov Zamir <kobi.zamir@gmail.com>
// and other contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mux

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Router registers routes to be matched and dispatches a handler.
//
// It implements the http.Handler interface, so it can be registered to serve
// requests:
//
//     var router = mux.Router{
//         NotFoundHandler: notFoundHandler,
//     }
//     router.HandleFunc("GET", "/val", getVaHandler)
//
//     func main() {
//         http.Handle("/", router)
//     }
type Router struct {
	// Configurable custom Handler to be used when no route matches.
	NotFoundHandler func(http.ResponseWriter, *http.Request)

	// List of http routes.
	routes []route
}

// HandleFunc registers a new route with a matcher for the URL path.
func (r *Router) HandleFunc(method string, path string, handler func(http.ResponseWriter, *http.Request)) {
	// Sanity check.
	if len(path) == 0 {
		return
	}

	// Get the path, add `/` at the beginning and remove `/` at the end.
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if path[0] != '/' {
		path = "/" + path
	}

	// Append a new route.
	segments := strings.Split(path, "/")[1:]
	r.routes = append(r.routes, route{
		method:   method,
		segments: segments,
		handler:  handler,
	})
}

// Var returns route variables for the current request using the route
// variable key, ok is true if key is found and value retrieved, o/w ok is false.
func Var(r *http.Request, key string) (string, bool) {
	// Try to get the context variabls.
	vars := r.Context().Value(ctxValsKey)
	if vars == nil {
		return "", false
	}

	// Try to convert the context variabls to a map.
	varsMap, ok := vars.(map[string]string)
	if !ok {
		return "", false
	}

	// Try to get the value we want.
	v, ok := varsMap[key]

	return v, ok
}

// ServeHTTP dispatches the handler registered in the matched route.
//
// When there is a match, route variables can be retrieved calling
// mux.Var(request, key).
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get the path, and clean it.
	path := req.URL.EscapedPath()
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Split path into it's segments.
	segments := strings.Split(path, "/")[1:]

	// Try to match the segments with one of the registered routs.
	for _, route := range r.routes {
		found, vars := r.match(route, req.Method, segments)

		// If found a match, run the handler for this route.
		if found {
			// Add path argv to the context.
			if len(vars) > 0 {
				req = req.WithContext(context.WithValue(req.Context(), ctxValsKey, vars))
			}

			route.handler(w, req)
			return
		}
	}

	// Handle page not found.
	if r.NotFoundHandler != nil {
		r.NotFoundHandler(w, req)
	} else {
		// If no custom "page not found" handler defined,
		// fallback to default 404.4 response.
		pageNotFound(w, req)
	}
}

// Internal context key type.
type ctxKey string

// The context key for the route parameters.
const ctxValsKey = ctxKey("Vals")

// Internal representation of a route.
type route struct {
	method   string
	segments []string
	handler  func(http.ResponseWriter, *http.Request)
}

// pageNotFound no handler configured.
func pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	io.WriteString(w, "404.4 – No handler configured.")
}

// match matches a request to a route, and parse the arguments embedded in the route path.
func (r Router) match(route route, method string, segments []string) (bool, map[string]string) {
	// Check request for method and segments length matching.
	if method != route.method || len(segments) != len(route.segments) {
		return false, nil
	}

	// Set a map for the path args, if found.
	vals := make(map[string]string)

	// Check each segment for a match.
	for i, segment := range route.segments {
		// Check for path argument.
		if segment[0] == ':' {
			// If this is an argument segments, parse it.
			value, _ := url.QueryUnescape(segments[i])
			vals[segment[1:]] = value

			continue
		}

		// Match current segment.
		if segments[i] != segment {
			// This request does not match the route.
			return false, nil
		}
	}

	// Found matching route.
	return true, vals
}
