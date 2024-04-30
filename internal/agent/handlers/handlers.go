package handlers

import (
	"fmt"
	"github.com/agatma/sprint1-http-server/internal/agent/storage"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type MetricsStorage struct {
	Metrics map[string]float64
	Mutex   sync.RWMutex
}

// sendMetrics sends metrics to a specified host.
func sendMetrics(host string, metricType string, metricName string, metricValue string) error {
	client := resty.New()
	resp, err := client.R().
		SetRawPathParams(map[string]string{
			"metricType":  metricType,
			"metricName":  strings.ToLower(metricName),
			"metricValue": metricValue,
		}).
		Post(fmt.Sprintf("%s/update/{metricType}/{metricName}/{metricValue}", host))

	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad request. Status Code %d", resp.StatusCode())
	}

	log.Printf("made request %s. Got status code %d", resp.Request.URL, resp.StatusCode())
	return nil
}

// SendGaugeMetrics sends gauge metrics to a specified host.
func SendGaugeMetrics(host string, metricStorage *storage.MetricsStorage) error {
	metricStorage.Mutex.RLock()
	defer metricStorage.Mutex.RUnlock()

	for metricName, metricValue := range metricStorage.Metrics {
		err := sendMetrics(host, "gauge", metricName, fmt.Sprintf("%f", metricValue))
		if err != nil {
			return err
		}
	}
	return nil
}

// SendCounterMetrics sends counter metrics to a specified host.
func SendCounterMetrics(host, metricName string, metricValue int64) error {
	return sendMetrics(host, "counter", metricName, strconv.FormatInt(metricValue, 10))
}
