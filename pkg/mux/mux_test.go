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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func foundHandler(w http.ResponseWriter, r *http.Request) {
	value, _ := Var(r, "key")

	w.WriteHeader(200)
	io.WriteString(w, fmt.Sprintf("{\"key\": \"%s\"}", value))
}

// ExampleRouter tests a route.
func ExampleRouter() {
	// A cat handler for the kitty request.
	catHandler := func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the ":uid" route parameter.
		uid, _ := Var(r, "uid")

		io.WriteString(w, fmt.Sprintf("{\"uid\": \"%s\", \"name\": \"Layla\"}", uid))
	}

	// Create a new router and egister our routes.
	router := Router{
		NotFoundHandler: notFoundHandler,
	}
	// Routes can have optional route parameters, in this example
	// route, ":uid" is a route parameter, once a route is dispatched,
	// the value of ":uid" can be retrieved using the
	// `mux.Var(*http.Request, string)` method.
	//
	// Example:
	// For a request "http://localhost:8080/kitty/eyfgt654efg7198u",
	// the value of ":uid" route parameter will be "eyfgt654efg7198u"
	router.HandleFunc("GET", "/kitty/:uid", catHandler)

	// Start the http server.
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Query server.
	res, err := http.Get(ts.URL + "/kitty/eyfgt654efg7198u")
	if err != nil {
		log.Fatal(err)
	}

	catJason, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", catJason)

	// Output: {"uid": "eyfgt654efg7198u", "name": "Layla"}
}

func benchmarkRouter(i int, b *testing.B) {
	// Create a router with 3 posible routes.
	handler := Router{
		NotFoundHandler: notFoundHandler,
	}
	handler.HandleFunc("GET", "/found", foundHandler)
	handler.HandleFunc("GET", "/found/:key", foundHandler)
	handler.HandleFunc("GET", "/found/:key/info", foundHandler)

	// Handle a request with on route parameter N times.
	for n := 0; n < b.N; n++ {
		req, err := http.NewRequest("GET", "/found/hello", nil)
		if err != nil {
			b.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			b.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `{"key": "hello"}`
		if rr.Body.String() != expected {
			b.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
}

// BenchmarkRouter100 run 10 requests.
func BenchmarkRouter10(b *testing.B) { benchmarkRouter(10, b) }

// BenchmarkRouter200 run 200 requests.
func BenchmarkRouter200(b *testing.B) { benchmarkRouter(200, b) }
