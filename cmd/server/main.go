package main

import (
	"github.com/agatma/sprint1-http-server/cmd/handlers"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.JSONFormatter{}
	log.Level = logrus.InfoLevel
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/gauge/`, handlers.GaugeMetricHandler)
	mux.HandleFunc(`/update/counter/`, handlers.CounterMetricHandler)
	mux.HandleFunc(`/update/`, handlers.UndefinedMetricType)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
