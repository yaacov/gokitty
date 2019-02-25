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

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAll(t *testing.T) {
	handler := newRouter()

	// Store new values.
	req, err := http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gorilla\": 123}"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Get all values.
	req, err = http.NewRequest("GET", "/val", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "{\"gorilla\":123,\"kitty\":\"cat\"}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGet(t *testing.T) {
	handler := newRouter()

	// Store new values.
	req, err := http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gorilla\": 123}"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Get gorilla.
	req, err = http.NewRequest("GET", "/val/gorilla", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "{\"gorilla\":123}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetMissing(t *testing.T) {
	handler := newRouter()

	// Store new values.
	req, err := http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gorilla\": 123}"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Get gorilla.
	req, err = http.NewRequest("GET", "/val/dog", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDelete(t *testing.T) {
	handler := newRouter()

	// Store new values.
	req, err := http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gorilla\": 123}"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Get gorilla.
	req, err = http.NewRequest("DELETE", "/val/gorilla", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "{\"gorilla\":123}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteMissing(t *testing.T) {
	handler := newRouter()

	// Store new values.
	req, err := http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gorilla\": 123}"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Get gorilla.
	req, err = http.NewRequest("DELETE", "/val/dog", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestPOST(t *testing.T) {
	handler := newRouter()

	// Store new values.
	req, err := http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gorilla\": 123}"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := "{\"gorilla\":123,\"kitty\":\"cat\"}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPOSTSameKyes(t *testing.T) {
	handler := newRouter()

	// Store new values.
	req, err := http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gorilla\": 123}"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Store same keys again.
	req, err = http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gorilla\": 321}"))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestPOSTBadRequest(t *testing.T) {
	handler := newRouter()

	// Store new values.
	req, err := http.NewRequest("POST", "/val", strings.NewReader("{\"kitty\": \"cat\", \"gori d gga ..^"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPUT(t *testing.T) {
	handler := newRouter()

	// Store new key value pair.
	req, err := http.NewRequest("PUT", "/val/kitty", strings.NewReader("\"cat\""))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := "{\"kitty\":\"cat\"}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Send same value again.
	req, err = http.NewRequest("PUT", "/val/kitty", strings.NewReader("\"cat\""))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotModified {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotModified)
	}

	// Check the response body is what we expect.
	expected = "{\"kitty\":\"cat\"}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Send modified value.
	req, err = http.NewRequest("PUT", "/val/kitty", strings.NewReader("\"sleepy\""))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = "{\"kitty\":\"sleepy\"}"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
