package file

import (
	"fmt"
	"metrics/internal/server/core/files"
	"sync"

	"metrics/internal/server/core/domain"
)

type MetricStorage struct {
	filepath  string
	syncWrite bool
	mux       *sync.Mutex
	metrics   map[domain.Key]domain.Value
}

func NewStorage(cfg *Config) (*MetricStorage, error) {
	if cfg.StoreInterval == 0 {
		return &MetricStorage{
			filepath:  cfg.Filepath,
			syncWrite: true,
			mux:       &sync.Mutex{},
			metrics:   make(map[domain.Key]domain.Value),
		}, nil
	} else {
		return &MetricStorage{
			filepath:  cfg.Filepath,
			syncWrite: false,
			mux:       &sync.Mutex{},
			metrics:   make(map[domain.Key]domain.Value),
		}, nil
	}
}

func (s *MetricStorage) SetMetric(m *domain.Metric) (*domain.Metric, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	var metric domain.Metric
	key := domain.Key{MType: m.MType, ID: m.ID}
	if m.MType == domain.Counter {
		value, found := s.metrics[key]
		if found {
			*value.Delta += *m.Delta
			s.metrics[key] = domain.Value{Delta: value.Delta}
			metric = domain.Metric{
				ID:    m.ID,
				MType: m.MType,
				Delta: value.Delta,
			}
		} else {
			s.metrics[key] = domain.Value{Delta: m.Delta}
			metric = domain.Metric{
				ID:    m.ID,
				MType: m.MType,
				Delta: m.Delta,
			}
		}
	} else {
		s.metrics[key] = domain.Value{Value: m.Value}
		metric = domain.Metric{
			ID:    m.ID,
			MType: m.MType,
			Value: m.Value,
		}
	}
	if s.syncWrite {
		err := files.SaveMetricsToFile(s.filepath, s.metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to save metrics to file %w", err)
		}
	}
	return &metric, nil
}

func (s *MetricStorage) GetMetric(mType, mName string) (*domain.Metric, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	value, found := s.metrics[domain.Key{MType: mType, ID: mName}]
	if !found {
		return &domain.Metric{}, domain.ErrItemNotFound
	}
	return &domain.Metric{
		ID:    mName,
		MType: mType,
		Value: value.Value,
		Delta: value.Delta,
	}, nil
}

func (s *MetricStorage) GetAllMetrics() (domain.MetricsList, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	metrics := make(domain.MetricsList, 0)
	for k, v := range s.metrics {
		metrics = append(metrics, domain.Metric{
			ID:    k.ID,
			MType: k.MType,
			Value: v.Value,
			Delta: v.Delta,
		})
	}
	return metrics, nil
}
