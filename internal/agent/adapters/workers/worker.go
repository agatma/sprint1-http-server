package workers

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"metrics/internal/agent/config"

	"metrics/internal/server/logger"
)

type AgentMetricService interface {
	UpdateMetrics(pollCount int) error
	SendMetrics(host string) error
}

type AgentWorker struct {
	agentMetricService AgentMetricService
	config             *config.Config
}

func NewAgentWorker(agentMetricService AgentMetricService, cfg *config.Config) *AgentWorker {
	return &AgentWorker{
		agentMetricService: agentMetricService,
		config:             cfg,
	}
}

func (a *AgentWorker) Run() error {
	address := strings.Split(a.config.Address, ":")
	port := "8080"
	if len(address) > 1 {
		port = address[1]
	}
	host := "http://localhost:" + port
	updateMetricsTicker := time.NewTicker(time.Duration(a.config.PollInterval) * time.Second)
	sendMetricsTicker := time.NewTicker(time.Duration(a.config.ReportInterval) * time.Second)
	pollCount := 0
	for {
		select {
		case <-updateMetricsTicker.C:
			pollCount++
			err := a.agentMetricService.UpdateMetrics(pollCount)
			if err != nil {
				return fmt.Errorf("failed to update metrics %w", err)
			}
		case <-sendMetricsTicker.C:
			err := a.agentMetricService.SendMetrics(host)
			if err != nil {
				logger.Log.Error("failed to send metrics", zap.Error(err))
			}
			pollCount = 0
		}
	}
}
