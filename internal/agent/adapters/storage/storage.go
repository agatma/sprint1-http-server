package storage

import (
	"errors"

	"metrics/internal/agent/adapters/storage/memory"
	"metrics/internal/agent/core/domain"
)

type AgentMetricStorage interface {
	GetMetricValue(request *domain.MetricRequest) *domain.MetricResponse
	SetMetricValue(request *domain.SetMetricRequest) *domain.SetMetricResponse
	GetAllMetrics(request *domain.GetAllMetricsRequest) *domain.GetAllMetricsResponse
}

func NewAgentStorage(conf Config) (AgentMetricStorage, error) {
	if conf.Memory != nil {
		return memory.NewAgentStorage(conf.Memory), nil
	}
	return nil, errors.New("no available agent storage")
}
