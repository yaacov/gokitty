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
	"net/http"

	"github.com/yaacov/gokitty/pkg/mux"
)

// Write a map[string]interface{} to response writer, or fail.
func writeMap(w http.ResponseWriter, m map[string]interface{}) {
	j, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}
	io.WriteString(w, string(j))
}

// notFound handles no found requests.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	io.WriteString(w, "{\"error\":\"not found\"}")
}

// getVal handles GET "/val" and GET "/val/:key" requests.
func getVal(w http.ResponseWriter, r *http.Request) {
	var m map[string]interface{}

	// Retrieve the ":key" route parameter.
	key, ok := mux.Var(r, "key")

	if ok {
		// Get one value by key:
		val, ok := vals.get(key)
		if ok {
			m = map[string]interface{}{key: val}
		} else {
			// We do not have this key in our store.
			w.WriteHeader(404)
			io.WriteString(w, fmt.Sprintf("{\"error\":\"can't find key %s\"}", key))

			return
		}
	} else {
		// Get all values:
		m = vals.list()
	}

	writeMap(w, m)
}

// postVal handles POST "/val" and PUT "/val" requests.
func postVal(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Read body data as json.
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprintf("{\"error\":\"%s\"}", err))
		return
	}

	// Store new data.
	for k, v := range data {
		vals.upsert(k, v)
	}

	// Write response as json.
	writeMap(w, data)
}

// deleteVal handles DELETE "/val/:key" requests.
func deleteVal(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ":key" route parameter.
	key, ok := mux.Var(r, "key")

	// Get one value by key:
	if ok {
		val, ok := vals.get(key)
		if ok {
			vals.delete(key)
			writeMap(w, map[string]interface{}{key: val})
		} else {
			w.WriteHeader(404)
			io.WriteString(w, fmt.Sprintf("{\"error\":\"can't find key %s\"}", key))
		}
	}
}
