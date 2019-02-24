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

// Global key value store.
var vals *Store

func newRouter() *mux.Router {
	r := mux.Router{
		NotFoundHandler: notFound,
	}
	r.HandleFunc("GET", "/val", getVal)
	r.HandleFunc("GET", "/val/:key", getVal)
	r.HandleFunc("POST", "/val", postVal)
	r.HandleFunc("PUT", "/val", postVal)
	r.HandleFunc("DELETE", "/val/:key", deleteVal)

	return &r
}

func main() {
	// Create a logging middleware, it's warm and fuzzy, prrr...
	logger := log.New(os.Stdout, "kitty: ", log.LstdFlags)
	loggingMiddleware := logging(logger)

	// Init our key value data store.
	vals = newStore()

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
