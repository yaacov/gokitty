// handlers_test.go
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

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func found(w http.ResponseWriter, r *http.Request) {
	value, _ := Var(r, "key")

	w.WriteHeader(200)
	io.WriteString(w, fmt.Sprintf("{\"key\": \"%s\"}", value))
}

func Example() {
	// A cat handler for the kitty request.
	catHandler := func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the ":uid" route parameter.
		uid, _ := Var(r, "uid")

		io.WriteString(w, fmt.Sprintf("{\"uid\": \"%s\", \"name\": \"Layla\"}", uid))
	}

	// Create a new router and egister our routes.
	router := Router{
		NotFoundHandler: notFound,
	}
	// Routes can have optional route parameters, in this example
	// route, ":uid" is a route parameter, once a route is dispatched,
	// the value of ":uid" can be retrieved using the `mux.Var(*http.Request, string)` method.
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

func TestNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/not-found", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Router{
		NotFoundHandler: notFound,
	}
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/found", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Router{
		NotFoundHandler: notFound,
	}
	handler.HandleFunc("GET", "/found", found)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRouteVars(t *testing.T) {
	req, err := http.NewRequest("GET", "/found/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Router{
		NotFoundHandler: notFound,
	}
	handler.HandleFunc("GET", "/found/:key", found)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"key": "hello"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
