package storage

import "sync"

type MetricsStorage struct {
	Metrics map[string]float64
	Mutex   sync.RWMutex
}

func NewMetricStorage() *MetricsStorage {
	return &MetricsStorage{
		Metrics: make(map[string]float64),
	}
}
