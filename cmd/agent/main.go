package main

import (
	"fmt"
	"github.com/agatma/sprint1-http-server/internal/agent/collector"
	"github.com/agatma/sprint1-http-server/internal/agent/handlers"
	"github.com/agatma/sprint1-http-server/internal/agent/storage"
	"log"
	"strings"
	"time"
)

func run(host string, metricsStorage *storage.MetricsStorage) {
	var PollCount int64
	collectMetricsTicker := time.NewTicker(options.pollInterval)
	sendMetricsTicker := time.NewTicker(options.reportInterval)
	for {
		select {
		case <-collectMetricsTicker.C:
			metrics := collector.CollectMetrics()
			metricsStorage.Metrics = metrics
			PollCount++
		case <-sendMetricsTicker.C:
			err := handlers.SendGaugeMetrics(host, metricsStorage)
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

func main() {
	parseFlags()
	metricStorage := storage.NewMetricStorage()
	port := strings.Split(flagRunAddr, ":")[1]
	host := fmt.Sprintf("http://127.0.0.1:%s", port)
	run(host, metricStorage)
}
