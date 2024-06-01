package service

import (
	"fmt"
	"strconv"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

type MetricStorage interface {
	GetMetric(mType, mName string) (*domain.Metric, error)
	SetMetric(m *domain.Metric) (*domain.Metric, error)
	GetAllMetrics() (domain.MetricsList, error)
}

type MetricService struct {
	storage MetricStorage
}

func NewMetricService(storage MetricStorage) *MetricService {
	return &MetricService{
		storage: storage,
	}
}

func (ms *MetricService) GetMetric(mType, mName string) (*domain.Metric, error) {
	metric, err := ms.storage.GetMetric(mType, mName)
	if err != nil {
		return metric, fmt.Errorf("failed to get metric: %w", err)
	}
	return metric, nil
}

func (ms *MetricService) SetMetric(m *domain.Metric) (*domain.Metric, error) {
	switch m.MType {
	case domain.Gauge, domain.Counter:
		metric, err := ms.storage.SetMetric(m)
		if err != nil {
			return metric, fmt.Errorf("%w", err)
		}
		return metric, nil
	default:
		return &domain.Metric{}, domain.ErrIncorrectMetricType
	}
}

func (ms *MetricService) SetMetricValue(req *domain.SetMetricRequest) (*domain.Metric, error) {
	switch req.MType {
	case domain.Gauge:
		value, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			return &domain.Metric{}, domain.ErrIncorrectMetricValue
		}
		metric, err := ms.storage.SetMetric(&domain.Metric{
			ID:    req.ID,
			MType: req.MType,
			Value: &value,
		})
		if err != nil {
			return metric, fmt.Errorf("%w", err)
		}
		return metric, nil
	case domain.Counter:
		value, err := strconv.Atoi(req.Value)
		if err != nil {
			return &domain.Metric{}, domain.ErrIncorrectMetricValue
		}
		valueInt := int64(value)
		metric, err := ms.storage.SetMetric(&domain.Metric{
			ID:    req.ID,
			MType: req.MType,
			Delta: &valueInt,
		})
		if err != nil {
			return metric, fmt.Errorf("%w", err)
		}
		return metric, nil
	default:
		return &domain.Metric{}, domain.ErrIncorrectMetricType
	}
}

func (ms *MetricService) GetMetricValue(mType, mName string) (string, error) {
	metric, err := ms.storage.GetMetric(mType, mName)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	switch mType {
	case domain.Gauge:
		value := strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		return value, nil
	case domain.Counter:
		value := strconv.Itoa(int(*metric.Delta))
		return value, nil
	default:
		return "", domain.ErrIncorrectMetricType
	}
}

func (ms *MetricService) GetAllMetrics() (domain.MetricsList, error) {
	metrics, err := ms.storage.GetAllMetrics()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return metrics, nil
}
