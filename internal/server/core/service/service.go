package service

import (
	"fmt"
	"strconv"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

type MetricStorage interface {
	GetMetric(mType, mName string) (*domain.Metrics, error)
	SetMetric(m *domain.Metrics) (*domain.Metrics, error)
	GetAllMetrics() domain.MetricsList
}

type MetricService struct {
	storage MetricStorage
}

func NewMetricService(storage MetricStorage) *MetricService {
	return &MetricService{
		storage: storage,
	}
}

func (ms *MetricService) GetMetric(mType, mName string) (*domain.Metrics, error) {
	return ms.storage.GetMetric(mType, mName)
}

func (ms *MetricService) SetMetric(m *domain.Metrics) (*domain.Metrics, error) {
	switch m.MType {
	case domain.Gauge, domain.Counter:
		return ms.storage.SetMetric(m)
	default:
		return &domain.Metrics{}, domain.ErrIncorrectMetricType
	}
}

func (ms *MetricService) SetMetricValue(req *domain.SetMetricRequest) (*domain.Metrics, error) {
	switch req.MType {
	case domain.Gauge:
		value, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			return &domain.Metrics{}, domain.ErrIncorrectMetricValue
		}
		return ms.storage.SetMetric(&domain.Metrics{
			ID:    req.ID,
			MType: req.MType,
			Value: &value,
		})
	case domain.Counter:
		value, err := strconv.Atoi(req.Value)
		if err != nil {
			return &domain.Metrics{}, domain.ErrIncorrectMetricValue
		}
		valueInt := int64(value)
		return ms.storage.SetMetric(&domain.Metrics{
			ID:    req.ID,
			MType: req.MType,
			Delta: &valueInt,
		})
	default:
		return &domain.Metrics{}, domain.ErrIncorrectMetricType
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

func (ms *MetricService) GetAllMetrics() domain.MetricsList {
	return ms.storage.GetAllMetrics()
}
