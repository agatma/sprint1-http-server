package main

import (
	"github.com/agatma/sprint1-http-server/internal/agent/collector"
	"github.com/agatma/sprint1-http-server/internal/agent/handlers"
	"github.com/agatma/sprint1-http-server/internal/agent/storage"
	"log"
	"time"
)

func main() {
	collectMetricsTicker := time.NewTicker(2 * time.Second)
	sendMetricsTicker := time.NewTicker(3 * time.Second)
	metricStorage := &storage.MetricsStorage{
		Metrics: make(map[string]float64),
	}
	var PollCount int64
	host := "http://localhost:8080"
	for {
		select {
		case <-collectMetricsTicker.C:
			metrics := collector.CollectMetrics()
			metricStorage.Metrics = metrics
			PollCount++
		case <-sendMetricsTicker.C:
			err := handlers.SendGaugeMetrics(host, metricStorage)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = handlers.SendCounterMetrics(host, "PollCount", PollCount)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
}
