package handlers

import (
	"fmt"
	"github.com/agatma/sprint1-http-server/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

// handleError is a helper function to handle HTTP errors.
func handleError(res http.ResponseWriter, errMsg string, statusCode int) {
	http.Error(res, errMsg, statusCode)
}

func AddMetric(res http.ResponseWriter, req *http.Request) {
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")
	if metricName == "" {
		handleError(res, "empty metric name", http.StatusNotFound)
		return
	}
	var err error
	switch metricType {
	case gauge:
		err = storage.Storage.AddGaugeValues(metricName, chi.URLParam(req, "metricValue"))
	case counter:
		err = storage.Storage.AddCounterValues(metricName, chi.URLParam(req, "metricValue"))
	default:
		handleError(res, "incorrect metric type", http.StatusBadRequest)
		return
	}
	if err != nil {
		handleError(res, "incorrect metric value", http.StatusBadRequest)
		return
	}
}

func GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")
	var v interface{}
	var found bool

	switch metricType {
	case gauge:
		v, found = storage.Storage.GetGaugeValues(metricName)
	case counter:
		v, found = storage.Storage.GetCounterValues(metricName)
	default:
		handleError(res, "incorrect metric type", http.StatusNotFound)
		return
	}
	if !found {
		handleError(res, "metric is not found", http.StatusNotFound)
		return
	}
	res.Write([]byte(fmt.Sprintf("%v", v)))
}

func GetAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	html := "<html><body><ul>"
	for key, value := range storage.Storage.GetAllGaugeValues() {
		html += fmt.Sprintf("<li>%s: %v</li>", key, value)
	}
	for key, value := range storage.Storage.GetAllCounterValues() {
		html += fmt.Sprintf("<li>%s: %v</li>", key, value)
	}
	html += "</ul></body></html>"
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte(html))
}
