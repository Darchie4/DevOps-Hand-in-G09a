package main

import (
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

