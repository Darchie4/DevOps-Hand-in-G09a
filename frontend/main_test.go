package main

import (
	"bytes"
	"encoding/json"
    "net/http"
    "testing"
)

// main_test.go
func TestHealthMarshaller(t *testing.T) {
	data := []byte(`[
        {"id":"2","message":"The measure of time to your next goal is the measure of your discipline."},
        {"id":"3","message":"The only way to do well is to do better each day."},
        {"id":"4","message":"It ain't over till it's EOF."},
        {"id":"1","message":"A new voyage will fill your life with untold memories."}
    ]`)

	reader := bytes.NewReader(data)
	status := HealthMarshaller(reader)

	expected := true
	if status != expected {
		t.Errorf("Marshalling did not result in correct fortunes")
	}
}
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
