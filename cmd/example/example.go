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

// Package main
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yaacov/gokitty/pkg/mux"
)

func newRouter() *mux.Router {
	// Create a new handler.
	h := newHandler()

	// Create a new router.
	r := mux.Router{
		NotFoundHandler: notFound,
	}
	r.HandleFunc("GET", "/val", h.getVal)
	r.HandleFunc("GET", "/val/:key", h.getVal)
	r.HandleFunc("POST", "/val", h.postVal)
	r.HandleFunc("PUT", "/val/:key", h.putVal)
	r.HandleFunc("DELETE", "/val/:key", h.deleteVal)

	return &r
}

func main() {
	// Create a logging middleware, it's warm and fuzzy, prrr...
	logger := log.New(os.Stdout, "kitty: ", log.LstdFlags)
	loggingMiddleware := logging(logger)

	// Register our routes.
	router := newRouter()

	// Serve on port 8080.
	s := &http.Server{
		Addr:           ":8080",
		Handler:        loggingMiddleware(router),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Println("Kitty key value server is starting ( try: http://localhost:8080/val ) ...")
	logger.Fatal(s.ListenAndServe())
}
