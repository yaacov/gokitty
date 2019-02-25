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

// Handler handle http requests.
type Handler struct {
	store *Store
}

func newHandler() *Handler {
	h := Handler{
		store: newStore(),
	}

	return &h
}

// Write an error.
func writeErr(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	io.WriteString(w, fmt.Sprintf("{\"error\":\"%s\"}", message))
}

// Write a key missing error.
func writeKeyErr(w http.ResponseWriter, key string) {
	writeErr(w, http.StatusNotFound, fmt.Sprintf("can't find key %s", key))
}

// notFound handles no found requests.
func notFound(w http.ResponseWriter, r *http.Request) {
	writeErr(w, http.StatusNotFound, "not found")
}

// Write a map[string]interface{} to response writer, or fail.
func writeMap(w http.ResponseWriter, m map[string]interface{}) {
	j, err := json.Marshal(m)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	io.WriteString(w, string(j))
}

// getVal handles GET "/val" and GET "/val/:key" requests.
func (h Handler) getVal(w http.ResponseWriter, r *http.Request) {
	var m map[string]interface{}

	// Retrieve the ":key" route parameter.
	key, ok := mux.Var(r, "key")

	if ok {
		// Get one value by key:
		val, ok := h.store.get(key)
		if ok {
			m = map[string]interface{}{key: val}
		} else {
			// We do not have this key in our store.
			writeKeyErr(w, key)
			return
		}
	} else {
		// Get all values:
		m = h.store.list()
	}

	writeMap(w, m)
}

// postVal handles POST "/val" requests.
func (h Handler) postVal(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Read body data as json.
	err := decoder.Decode(&data)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check for newly created keys.
	for k := range data {
		if _, ok := h.store.get(k); !ok {
			w.WriteHeader(http.StatusCreated)
			break
		}
	}

	// Create or modify multiple key value pairs.
	for k, v := range data {
		h.store.upsert(k, v)
	}

	// Write response as json.
	writeMap(w, data)
}

// postVal handles PUT "/val/:key" requests.
func (h Handler) putVal(w http.ResponseWriter, r *http.Request) {
	var data interface{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Read body data as json.
	err := decoder.Decode(&data)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	// Retrieve the ":key" route parameter.
	key, ok := mux.Var(r, "key")
	if !ok {
		writeErr(w, http.StatusInternalServerError, "can't get key")
	}

	// Check if this is a new key.
	if val, ok := h.store.get(key); ok {
		// We are modifying an existing value.
		if val == data {
			// Value does not require change.
			w.WriteHeader(http.StatusNotModified)
			writeMap(w, map[string]interface{}{key: data})
			return
		}
	} else {
		// We are creating a new key value pair.
		w.WriteHeader(http.StatusCreated)
	}

	// Create or modify key value pair.
	h.store.upsert(key, data)

	// Write response as json.
	writeMap(w, map[string]interface{}{key: data})
}

// deleteVal handles DELETE "/val/:key" requests.
func (h Handler) deleteVal(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ":key" route parameter.
	key, ok := mux.Var(r, "key")
	if !ok {
		writeErr(w, http.StatusInternalServerError, "can't get key")
	}

	// Get one value by key:
	val, ok := h.store.get(key)
	if ok {
		h.store.delete(key)
		writeMap(w, map[string]interface{}{key: val})
	} else {
		writeKeyErr(w, key)
	}
}
