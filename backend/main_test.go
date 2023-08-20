package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthz(t *testing.T) {
    _, err := http.NewRequest("GET", "/healthz", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()

    status := rr.Code;

    if status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := "healthy"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
