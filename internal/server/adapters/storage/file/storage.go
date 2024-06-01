package file

import (
	"fmt"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
	"github.com/agatma/sprint1-http-server/internal/server/core/files"
)

type MetricStorage struct {
	filepath string
}

func NewStorage(cfg *Config) (*MetricStorage, error) {
	return &MetricStorage{
		filepath: cfg.Filepath,
	}, nil
}

func (s *MetricStorage) SetMetric(m *domain.Metric) (*domain.Metric, error) {
	var metric domain.Metric
	metricValues, err := files.LoadMetricsFromFile(s.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load metrics from file %w", err)
	}
	key := domain.Key{MType: m.MType, ID: m.ID}
	if m.MType == domain.Counter {
		value, found := metricValues[key]
		if found {
			*value.Delta += *m.Delta
			metricValues[key] = domain.Value{Delta: value.Delta}

			metric = domain.Metric{
				ID:    m.ID,
				MType: m.MType,
				Delta: value.Delta,
			}
		} else {
			metricValues[key] = domain.Value{Delta: m.Delta}
			metric = domain.Metric{
				ID:    m.ID,
				MType: m.MType,
				Delta: m.Delta,
			}
		}
	} else {
		metricValues[key] = domain.Value{Value: m.Value}
		metric = domain.Metric{
			ID:    m.ID,
			MType: m.MType,
			Value: m.Value,
		}
	}
	err = files.SaveMetricsToFile(s.filepath, metricValues)
	if err != nil {
		return nil, fmt.Errorf("failed to save metrics to file %w", err)
	}
	return &metric, nil
}

func (s *MetricStorage) GetMetric(mType, mName string) (*domain.Metric, error) {
	metricValues, err := files.LoadMetricsFromFile(s.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load metrics from file %w", err)
	}
	key := domain.Key{MType: mType, ID: mName}
	value, found := metricValues[key]
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
	metricValues, err := files.LoadMetricsFromFile(s.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load metrics from file %w", err)
	}
	metrics := make(domain.MetricsList, 0)
	for k, v := range metricValues {
		metrics = append(metrics, domain.Metric{
			ID:    k.ID,
			MType: k.MType,
			Value: v.Value,
			Delta: v.Delta,
		})
	}
	return metrics, nil
}
