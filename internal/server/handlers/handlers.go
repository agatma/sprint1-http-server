package handlers

import (
	"github.com/agatma/sprint1-http-server/internal/server/storage"
	"log"
	"net/http"
	"strings"
)

func GaugeMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", req.Method)
		return
	}
	urlPath := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(urlPath) < 4 {
		http.Error(res, "incorrect input path", http.StatusNotFound)
		return
	}
	metricName, metricValue := urlPath[2], urlPath[3]
	if metricName == "" {
		http.Error(res, "empty metric name", http.StatusNotFound)
		return
	}
	err := storage.Storage.AddGaugeValues(metricName, metricValue)
	if err != nil {
		http.Error(res, "incorrect metric value", http.StatusBadRequest)
		return
	}
	log.Printf("current storage state: %v", storage.Storage)
}

func UndefinedMetricType(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "incorrect metric type", http.StatusBadRequest)
}

func CounterMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	urlPath := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(urlPath) < 4 {
		http.Error(res, "incorrect input path", http.StatusNotFound)
		return
	}
	metricName, metricValue := urlPath[2], urlPath[3]
	if metricName == "" {
		http.Error(res, "empty metric name", http.StatusNotFound)
		return
	}
	err := storage.Storage.AddCounterValues(metricName, metricValue)
	if err != nil {
		http.Error(res, "incorrect metric value", http.StatusBadRequest)
		return
	}
	log.Printf("current storage state: %v", storage.Storage)
}
