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
// mux.Router supports precise routes, route parameters, and not found handler.
//
// Precise routes, unlike http mux, kitty routes are precise,
// request to path "/hello/world" will not match the route "/hello/".
//
// NotFoundHandler is a custom handler function called when all routes does not match,
// users should define a not found handler when using kitty mux router.
// If NotFoundHandler is not defined a default "404" handler is used.
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
//  // Define a route with "key" route parameter.
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
//      Handler:        myRouter,
//  }
//  log.Fatal(s.ListenAndServe())
package mux
