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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yaacov/gokitty/pkg/mux"
)

var vals map[string]string

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			next.ServeHTTP(w, r)
		})
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	io.WriteString(w, fmt.Sprintf("{\"error\":\"not found\"}"))
}

func getVal(w http.ResponseWriter, r *http.Request) {
	key, ok := mux.GetPathArg(r, "key")

	// Get one value by key:
	if ok {
		val, ok := vals[key]
		if ok {
			io.WriteString(w, fmt.Sprintf("{\"%s\":\"%s\"}", key, val))
		} else {
			w.WriteHeader(404)
			io.WriteString(w, fmt.Sprintf("{\"error\":\"can't find key %s\"}", key))
		}
		return
	}

	// Get all values:
	j, err := json.Marshal(vals)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}
	io.WriteString(w, string(j))
}

func postVal(w http.ResponseWriter, r *http.Request) {
	var j []byte
	var data map[string]string

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Read body data as json.
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}

	// Write response as json.
	j, err = json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}

	for k, v := range data {
		vals[k] = v
	}
	io.WriteString(w, string(j))
}

func deleteVal(w http.ResponseWriter, r *http.Request) {
	key, ok := mux.GetPathArg(r, "key")

	// Get one value by key:
	if ok {
		val, ok := vals[key]
		if ok {
			delete(vals, key)
			io.WriteString(w, fmt.Sprintf("{\"%s\":\"%s\"}", key, val))
		} else {
			w.WriteHeader(404)
			io.WriteString(w, fmt.Sprintf("{\"error\":\"can't find key %s\"}", key))
		}
	}
}

func main() {
	logger := log.New(os.Stdout, "kitty: ", log.LstdFlags)
	loggingMiddleware := logging(logger)

	vals = make(map[string]string)

	myRouter := mux.Router{
		NotFoundHandler: notFound,
	}
	myRouter.HandleFunc("GET", "/val", getVal)
	myRouter.HandleFunc("GET", "/val/:key", getVal)
	myRouter.HandleFunc("POST", "/val", postVal)
	myRouter.HandleFunc("PUT", "/val", postVal)
	myRouter.HandleFunc("DELETE", "/val/:key", deleteVal)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        loggingMiddleware(myRouter),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Println("Kitty key value server is starting ( try: http://localhost:8080/val ) ...")
	log.Fatal(s.ListenAndServe())
}
