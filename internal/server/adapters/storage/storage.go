package storage

import (
	"errors"

	"github.com/agatma/sprint1-http-server/internal/server/adapters/storage/memory"
	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

type MetricStorage interface {
	GetMetric(mType, mName string) (*domain.Metrics, error)
	SetMetric(m *domain.Metrics) (*domain.Metrics, error)
	GetAllMetrics() domain.MetricsList
}

func NewStorage(conf Config) (MetricStorage, error) {
	if conf.Memory != nil {
		return memory.NewStorage(conf.Memory), nil
	}
	return nil, errors.New("no available storage")
}
