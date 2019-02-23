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

// Package mux is a small, fast and cute http mux package.
//
// mux.Router supports precise routes, route paramters, and not found handler.
//
// Precise routes, unlike http mux, kitty routes are precise,
// request to path "/hello/world" will not match the route "/hello/".
//
// NotFoundHandler is a handler function called when all routes does not match,
// users should define a nut found handler when using kitty mux router.
//
// Route parameters are named URL segments that are used to capture the values
// specified at their position in the URL. The captured values
// retrieved calling mux.Var(request, key), with the name of the route parameter
// as key.
//
// To define routes with route parameters, simply specify the route parameters
// in the path of the route as shown below.
//
// Example:
//  //  Define a route with "key" route parameter.
//  router.HandleFunc("GET", "/val/:key", getValHandler)
//
// Usage:
//  func getValHandler(w http.ResponseWriter, r *http.Request) {
//      // Retrieve rount variables.
//      key, ok := mux.Var(r, "key")
//      ...
//      action, ok := mux.Var(r, "action")
//      ...
//  }
//
//  ...
//
//  myRouter := mux.Router{
//      NotFoundHandler: notFound,
//  }
//  myRouter.HandleFunc("GET", "/val/:key/:action", getValHandler)
//
//  s := &http.Server{
//      Addr:           ":8080",
//      Handler:        loggingMiddleware(myRouter),
//  }
//  log.Fatal(s.ListenAndServe())
package mux

import (
	"context"
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
	// Configurable Handler to be used when no route matches.
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

	// Get the path, and clean it.
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if path[0] != '/' {
		path = "/" + path
	}

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
	argv := r.Context().Value(ctxKey("argv"))
	if argv == nil {
		return "", false
	}

	argvMap, ok := argv.(map[string]string)
	return argvMap[key], ok
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
		found, argv := r.match(route, req.Method, segments)

		// iIf found a match, run the handler for this route.
		if found {
			// Add path argv to the context.
			if len(argv) > 0 {
				req = req.WithContext(context.WithValue(req.Context(), ctxKey("argv"), argv))
			}

			route.handler(w, req)
			return
		}
	}

	// Handle page not found.
	r.NotFoundHandler(w, req)
}

// Internal context key type.
type ctxKey string

// Internal representation of a route.
type route struct {
	method   string
	segments []string
	handler  func(http.ResponseWriter, *http.Request)
}

// match matches a request to a route, and parse the arguments embedded in the route path.
func (r Router) match(route route, method string, segments []string) (bool, map[string]string) {
	// Check request for method and segments length matching.
	if method != route.method || len(segments) != len(route.segments) {
		return false, nil
	}

	// Set a map for the path args, if found.
	argv := make(map[string]string)

	// Check each segment for a match.
	for i, segment := range route.segments {
		// Check for path argument.
		if segment[0] == ':' {
			// If this is an argument segments, parse it.
			value, _ := url.QueryUnescape(segments[i])
			argv[segment[1:]] = value

			continue
		}

		// Match current segment.
		if segments[i] != segment {
			// This request does not match the route.
			return false, nil
		}
	}

	// Found matching route.
	return true, argv
}
