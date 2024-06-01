package memory

import (
	"sync"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

type Key struct {
	MType string
	ID    string
}

type Value struct {
	Value *float64
	Delta *int64
}

type MetricStorage struct {
	mux     *sync.Mutex
	metrics map[Key]Value
}

func NewStorage(cfg *Config) *MetricStorage {
	return &MetricStorage{
		mux:     &sync.Mutex{},
		metrics: make(map[Key]Value),
	}
}

func (s *MetricStorage) GetMetric(mType, mName string) (*domain.Metrics, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	value, found := s.metrics[Key{MType: mType, ID: mName}]
	if !found {
		return &domain.Metrics{}, domain.ErrItemNotFound
	}
	return &domain.Metrics{
		ID:    mName,
		MType: mType,
		Value: value.Value,
		Delta: value.Delta,
	}, nil
}

func (s *MetricStorage) SetMetric(m *domain.Metrics) (*domain.Metrics, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	key := Key{MType: m.MType, ID: m.ID}
	if m.MType == domain.Counter {
		value, found := s.metrics[key]
		if found {
			*value.Delta += *m.Delta
			s.metrics[key] = Value{Delta: value.Delta}
			return &domain.Metrics{
				ID:    m.ID,
				MType: m.MType,
				Delta: value.Delta,
			}, nil
		} else {
			s.metrics[key] = Value{Delta: m.Delta}
			return &domain.Metrics{
				ID:    m.ID,
				MType: m.MType,
				Delta: m.Delta,
			}, nil
		}
	} else {
		s.metrics[key] = Value{Value: m.Value}
		return &domain.Metrics{
			ID:    m.ID,
			MType: m.MType,
			Value: m.Value,
		}, nil
	}
}

func (s *MetricStorage) GetAllMetrics() domain.MetricsList {
	s.mux.Lock()
	defer s.mux.Unlock()
	metrics := make(domain.MetricsList, 0)
	for k, v := range s.metrics {
		metrics = append(metrics, domain.Metrics{
			ID:    k.ID,
			MType: k.MType,
			Value: v.Value,
			Delta: v.Delta,
		})
	}
	return metrics
}
