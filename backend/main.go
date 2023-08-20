//nolint:all
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listFortuneRe          = regexp.MustCompile(`^/fortunes[/]*$`)
	getFortuneRe           = regexp.MustCompile(`^/fortunes[/](\d+)$`)
	randomFortuneRe        = regexp.MustCompile(`^/fortunes[/]random$`)
	createFortuneRe        = regexp.MustCompile(`^/fortunes[/]*$`)
	customFortunesCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "custom_fortunes_total",
		Help: "The total number of custom fortune cookies created",
	})
	fortunesGiven = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fortunes_granted_total",
		Help: "The total number of custom fortunes granted",
	})
	healtz = regexp.MustCompile(`^/healthz[/]*$`)

)

type fortune struct {
	ID      string `json:"id" redis:"id"`
	Message string `json:"message" redis:"message"`
}

type datastore struct {
	m map[string]fortune
	*sync.RWMutex
}

var datastoreDefault = datastore{m: map[string]fortune{
	"1": {ID: "1", Message: "A new voyage will fill your life with untold memories."},
	"2": {ID: "2", Message: "The measure of time to your next goal is the measure of your discipline."},
	"3": {ID: "3", Message: "The only way to do well is to do better each day."},
	"4": {ID: "4", Message: "It ain't over till it's EOF."},
}, RWMutex: &sync.RWMutex{}}

type fortuneHandler struct {
	store *datastore
}

func (h *fortuneHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listFortuneRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && getFortuneRe.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodGet && randomFortuneRe.MatchString(r.URL.Path):
		h.Random(w, r)
		return
	case r.Method == http.MethodPost && createFortuneRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodGet && healtz.MatchString(r.URL.Path):
		h.Healthz(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

func (h *fortuneHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("healthy"))
	if err != nil {
		log.Fatal(err)
	}
}


func (h *fortuneHandler) List(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	fortunes := make([]fortune, 0, len(h.store.m))
	for _, v := range h.store.m {
		fortunes = append(fortunes, v)
	}
	h.store.RUnlock()
	fortunesGiven.Add(float64(len(h.store.m)))
	jsonBytes, err := json.Marshal(fortunes)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *fortuneHandler) Random(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	fortunes := make([]fortune, 0, len(h.store.m))
	for _, v := range h.store.m {
		fortunes = append(fortunes, v)
	}
	h.store.RUnlock()
	fortunesGiven.Inc()
	if len(fortunes) > 0 {
		u := fortunes[rand.Intn(len(fortunes))]
		r.URL.Path = "/fortunes/" + u.ID
	} else {
		r.URL.Path = "/fortunes/zero"
	}

	h.Get(w, r)
}

func (h *fortuneHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := getFortuneRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}
	
	if usingRedis {
		key := matches[1]
		val, err := dbLink.Do("hget", "fortunes", key)
		if err != nil {
			fmt.Println("redis hget failed", err.Error())
		} else {
			if val != nil {
				msg := fmt.Sprintf("%s", val.([]byte))
				h.store.Lock()
				h.store.m[key] = fortune{ID: key, Message: msg}
				h.store.Unlock()
			}
		}
	}

	h.store.RLock()
	u, ok := h.store.m[matches[1]]
	h.store.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("fortune not found"))
		return
	}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *fortuneHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u fortune
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		internalServerError(w, r)
		return
	}
	customFortunesCreated.Inc()
	h.store.Lock()
	h.store.m[u.ID] = u
	h.store.Unlock()
	fmt.Printf("using redis = %t \n", usingRedis)
	if usingRedis {
		_, err := dbLink.Do("hset", "fortunes", u.ID, u.Message)
		if err != nil {
			fmt.Println("redis hset failed", err.Error())
		}
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("internal server error"))
	if err != nil {
		log.Fatal(err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("not found"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//http.HandleFunc("/healthz", HealthzHandler)
	mux := http.NewServeMux()
	fortuneH := &fortuneHandler{
		store: &datastoreDefault,
	}
	mux.Handle("/fortunes", fortuneH)
	mux.Handle("/fortunes/", fortuneH)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/healthz", fortuneH)
	err := http.ListenAndServe(":9000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
