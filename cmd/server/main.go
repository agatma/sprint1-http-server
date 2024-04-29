package main

import (
	"github.com/agatma/sprint1-http-server/internal/server/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/gauge/`, handlers.GaugeMetricHandler)
	mux.HandleFunc(`/update/counter/`, handlers.CounterMetricHandler)
	mux.HandleFunc(`/update/`, handlers.UndefinedMetricType)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
