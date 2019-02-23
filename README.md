# GoKitty
Small, fast and cute URL router and dispatcher for golang.

[![Go Report Card](https://goreportcard.com/badge/github.com/yaacov/gokitty)](https://goreportcard.com/report/github.com/yaacov/gokitty)
[![Build Status](https://travis-ci.org/yaacov/gokitty.svg?branch=master)](https://travis-ci.org/yaacov/gokitty)
[![GoDoc](https://godoc.org/github.com/yaacov/gokitty/pkg/mux?status.svg)](https://godoc.org/github.com/yaacov/gokitty/pkg/mux)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Because kittens are cuter than gorillas.

Like the standard `http.ServeMux`, `gokitty/pkg/mux` matches incoming requests against a list of registered routes and calls a handler for the route that matches the URL or other conditions. The main features are:
- Precise routes, unlike `http.ServeMux`, kitty does not use patterns, routes must match requested path exectly, if all routes fail, kitty will call the not found handler.
- NotFoundHandler is a handler function called when all routes does not match.
- Route parameters are named URL segments that are used to capture the values specified at their position in the URL.

# Example

``` go
import (
	"net/http"
  ...

	"github.com/yaacov/gokitty/pkg/mux"
)

// notFound handles not found requests.
func notFound(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(404)
  io.WriteString(w, fmt.Sprintf("{\"error\":\"not found\"}"))
}

// getVal handles GET "/val" and GET "/val/:key" requests.
func getVal(w http.ResponseWriter, r *http.Request) {
  // Retrieve the ":key" route parameter.
  key, ok := mux.Var(r, "key")
  ...
}

...

// Register our routes.
myRouter := mux.Router{
  NotFoundHandler: notFound,
}
myRouter.HandleFunc("GET", "/val", getVal)
myRouter.HandleFunc("GET", "/val/:key", getVal)

// Serve on port 8080.
s := &http.Server{
  Addr:           ":8080",
  Handler:        loggingMiddleware(myRouter),
}
log.Fatal(s.ListenAndServe())

```
