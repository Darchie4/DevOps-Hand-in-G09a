package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "io/ioutil"
)

func TestHealthz(t *testing.T) {
	mux := http.NewServeMux()
	fortuneH := &fortuneHandler{
		store: &datastoreDefault,
	}

	mux.Handle("/healthz", fortuneH)

	server := httptest.NewServer(mux)

	// Make a GET request to /healthz
	resp, err := http.Get(server.URL + "/healthz")
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	// Check the response body content
	expectedBody := "healthy"
	if string(body) != expectedBody {
		t.Errorf("Expected response body %q, but got %q", expectedBody, string(body))
	}

}
