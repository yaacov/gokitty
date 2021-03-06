# GoKitty
Small, fast and cute URL router and dispatcher for golang.

[![Go Report Card](https://goreportcard.com/badge/github.com/yaacov/gokitty)](https://goreportcard.com/report/github.com/yaacov/gokitty)
[![Build Status](https://travis-ci.org/yaacov/gokitty.svg?branch=master)](https://travis-ci.org/yaacov/gokitty)
[![GoDoc](https://godoc.org/github.com/yaacov/gokitty/pkg/mux?status.svg)](https://godoc.org/github.com/yaacov/gokitty/pkg/mux)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Because kittens are cuter than gorillas.

<p align="center">
  <img src="https://raw.githubusercontent.com/yaacov/gokitty/master/img/net.jpg" alt="Random cat image">
</p>

Like the standard `http.ServeMux`, `gokitty/pkg/mux` matches incoming requests against a list of registered routes and calls a handler for the route that matches the URL or other conditions. The main features are using precise routes, implementing a not found handler and using route parameters.

- Precise routes, unlike `http.ServeMux`, kitty does not use patterns, routes must match requested path exectly, if all routes fail, kitty will call the not found handler.
- NotFoundHandler is a handler function called when all routes does not match, if not defined, a default "404" handler is used.
- Route parameters are named URL segments that are used to capture the values specified at their position in the URL.

# Install

``` bash
go get -u github.com/yaacov/gokitty
```

# Example

``` go
import (
	"net/http"
  
  ...

	"github.com/yaacov/gokitty/pkg/mux"
)

...

// notFound handles path not found requests.
func notFound(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(404)
}

// getVal handles GET "/val/:uid" requests.
func getVal(w http.ResponseWriter, r *http.Request) {
  // Retrieve the ":uid" route parameter.
  key, ok := mux.Var(r, "uid")
  if !ok {
    // Route parameter retrieval faild.
    ...
  }
  
  ...
  
}

...

// Create a new router and egister our routes.
router := mux.Router{
  // NotFoundHandler is optional custom "404 - page not found" handler,
  // if not defined a default handler is used.
  NotFoundHandler: notFound,
}
// Routes can have optional route parameters, in this example
// route, ":uid" is a route parameter, once a route is dispatched, 
// the value of ":uid" can be retrieved using the `mux.Var(*http.Request, string)` method.
//
// Example:
// For a request "http://localhost:8080/val/eyfgt654efg7198u",
// the value of ":uid" route parameter will be "eyfgt654efg7198u"
router.HandleFunc("GET", "/val/:uid", getVal)

// Serve on port 8080.
s := &http.Server{
  Addr:           ":8080",
  Handler:        router,
}
s.ListenAndServe()

```
# Gopher image

https://github.com/egonelbre/gophers
