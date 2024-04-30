package main

import (
	"github.com/agatma/sprint1-http-server/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", handlers.AddMetric)
	})
	r.Get("/value/{metricType}/{metricName}", handlers.GetMetric)
	r.Get("/", handlers.GetAllMetricsHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
