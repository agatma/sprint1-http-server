package storage

import (
	"errors"
	"fmt"
	"metrics/internal/server/adapters/storage/database"

	"metrics/internal/server/adapters/storage/file"
	"metrics/internal/server/adapters/storage/memory"
	"metrics/internal/server/core/domain"
)

type MetricStorage interface {
	GetMetric(mType, mName string) (*domain.Metric, error)
	SetMetric(m *domain.Metric) (*domain.Metric, error)
	GetAllMetrics() (domain.MetricsList, error)
	Ping() error
}

func NewStorage(cfg Config) (MetricStorage, error) {
	if cfg.Database != nil {
		storage, err := database.NewStorage(cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		return storage, nil
	}
	if cfg.Memory != nil {
		storage, err := memory.NewStorage(cfg.Memory)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		return storage, nil
	}
	if cfg.File != nil {
		storage, err := file.NewStorage(cfg.File)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		return storage, nil
	}
	return nil, errors.New("no available storage")
}
