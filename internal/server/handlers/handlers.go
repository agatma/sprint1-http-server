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
	memStorage := storage.GetMemStorage()
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")
	if metricName == "" {
		handleError(res, "empty metric name", http.StatusNotFound)
		return
	}
	var err error
	switch metricType {
	case gauge:
		err = memStorage.AddGaugeValues(metricName, chi.URLParam(req, "metricValue"))
	case counter:
		err = memStorage.AddCounterValues(metricName, chi.URLParam(req, "metricValue"))
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
	memStorage := storage.GetMemStorage()
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")
	var v interface{}
	var found bool

	switch metricType {
	case gauge:
		v, found = memStorage.GetGaugeValues(metricName)
	case counter:
		v, found = memStorage.GetCounterValues(metricName)
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
	memStorage := storage.GetMemStorage()
	html := "<html><body><ul>"
	for key, value := range memStorage.GetAllGaugeValues() {
		html += fmt.Sprintf("<li>%s: %v</li>", key, value)
	}
	for key, value := range memStorage.GetAllCounterValues() {
		html += fmt.Sprintf("<li>%s: %v</li>", key, value)
	}
	html += "</ul></body></html>"
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte(html))
}
