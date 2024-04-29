package handlers

import (
	"fmt"
	"github.com/agatma/sprint1-http-server/internal/agent/storage"
	"log"
	"net/http"
	"sync"
)

type MetricsStorage struct {
	Metrics map[string]float64
	Mutex   sync.RWMutex
}

func SendGaugeMetrics(host string, metricStorage *storage.MetricsStorage) error {
	metricStorage.Mutex.RLock()
	defer metricStorage.Mutex.RUnlock()
	for metricName, metricValue := range metricStorage.Metrics {
		url := fmt.Sprintf("%s/update/gauge/%s/%f", host, metricName, metricValue)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		log.Printf("Making request: %s, %s", req.Method, req.URL)
		if err != nil {
			return err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad request.Status Code %d", resp.StatusCode)
		}
	}
	return nil
}

func SendCounterMetrics(host, metricName string, metricValue int64) error {
	url := fmt.Sprintf("%s/update/counter/%s/%d", host, metricName, metricValue)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	log.Printf("Making request: %s, %s", req.Method, req.URL)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request.Status Code %d", resp.StatusCode)
	}
	return nil
}
