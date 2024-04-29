package storage

import "sync"

type MetricsStorage struct {
	Metrics map[string]float64
	Mutex   sync.RWMutex
}
