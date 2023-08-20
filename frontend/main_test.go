package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestHealthz(t *testing.T) {

    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
    // pass 'nil' as the third parameter.
    req, err := http.NewRequest("GET", "/healthz", nil)
    if err != nil {
        t.Fatal(err)
    }

    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(HealthzHandler)

    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method
    // directly and pass in our Request and ResponseRecorder.
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Check the response body is what we expect.
    expected := "healthy"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func BenchmarkGetAll(b *testing.B) {
	resp, err := http.Get("http://localhost:9000/fortunes")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	// Convert reponse to json list
	list := []fortune{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		panic(err)
	}

	// Iterate over list and delete all keys
	for _, f := range list {
		var body = []byte(f.ID)
		_, err := http.Post("http://localhost:9000/api/delete", "text/plain", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
	}

	// Add keys
	for i := 0; i < 100; i++ {
		var body = []byte(`{"id":"` + strconv.Itoa(i) + `", "message": "test"}`)
		_, err := http.Post("http://localhost:8081/api/add", "application/json", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
	}

	// Request all keys
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		_, err := http.Get("http://localhost:8081/api/all")
		if err != nil {
			panic(err)
		}
	}

}
