package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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
