package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/agatma/sprint1-http-server/internal/server/config"
	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
	"github.com/agatma/sprint1-http-server/internal/server/core/files"
	"github.com/agatma/sprint1-http-server/internal/server/logger"
	"go.uber.org/zap"
)

type MetricStorage interface {
	GetMetric(mType, mName string) (*domain.Metric, error)
	SetMetric(m *domain.Metric) (*domain.Metric, error)
	GetAllMetrics() (domain.MetricsList, error)
}

type MetricService struct {
	storage  MetricStorage
	filepath string
}

func NewMetricService(cfg *config.Config, storage MetricStorage) (*MetricService, error) {
	ms := MetricService{
		storage:  storage,
		filepath: cfg.FileStoragePath,
	}
	if cfg.Restore {
		err := ms.loadMetricsFromFile()
		if err != nil {
			return nil, fmt.Errorf("failed to restore data for metric service %w", err)
		}
	}
	if cfg.StoreInterval > 0 {
		timeDuration := time.Duration(cfg.StoreInterval) * time.Second
		time.AfterFunc(timeDuration, func() {
			err := ms.SaveMetricsToFile()
			if err != nil {
				logger.Log.Error("failed to save metrics", zap.Error(err))
			}
			logger.Log.Info("metrics saved to file after timeout", zap.Int("seconds", cfg.StoreInterval))
		})
	}
	return &ms, nil
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

func (ms *MetricService) SaveMetricsToFile() error {
	metricValues := make(domain.MetricValues)
	metrics, err := ms.storage.GetAllMetrics()
	if err != nil {
		return fmt.Errorf("failed to get metrics for saving to file: %w", err)
	}
	for _, v := range metrics {
		metricValues[domain.Key{ID: v.ID, MType: v.MType}] = domain.Value{Value: v.Value, Delta: v.Delta}
	}
	err = files.SaveMetricsToFile(ms.filepath, metricValues)
	if err != nil {
		return fmt.Errorf("failed to save metrics to file: %w", err)
	}
	return nil
}

func (ms *MetricService) loadMetricsFromFile() error {
	metrics, err := files.LoadMetricsFromFile(ms.filepath)
	if err != nil {
		return fmt.Errorf("failed to load metrics for restore: %w", err)
	}
	for k, v := range metrics {
		_, err := ms.storage.SetMetric(&domain.Metric{
			ID:    k.ID,
			MType: k.MType,
			Value: v.Value,
			Delta: v.Delta,
		})
		if err != nil {
			return fmt.Errorf("failed to save metrics in restore: %w", err)
		}
	}
	return nil
}
