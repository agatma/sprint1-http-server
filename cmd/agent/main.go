package main

import (
	"fmt"
	"log"
	"os"

	"github.com/agatma/sprint1-http-server/internal/agent/adapters/storage"
	"github.com/agatma/sprint1-http-server/internal/agent/adapters/storage/memory"
	"github.com/agatma/sprint1-http-server/internal/agent/adapters/workers"
	"github.com/agatma/sprint1-http-server/internal/agent/core/service"
	"github.com/agatma/sprint1-http-server/internal/agent/logger"
)

func main() {
	f, _ := os.Create("agent_log.txt")
	defer f.Close()
	log.SetOutput(f)
	if err := run(); err != nil {
		log.Println(err)
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := workers.NewConfig()
	if err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}
	if err = logger.Initialize(cfg.LogLevel); err != nil {
		return fmt.Errorf("can't load logger: %w", err)
	}
	gaugeAgentStorage, err := storage.NewAgentStorage(storage.Config{
		Memory: &memory.Config{},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize a storage: %w", err)
	}
	counterAgentStorage, err := storage.NewAgentStorage(storage.Config{
		Memory: &memory.Config{},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize a storage: %w", err)
	}
	agentMetricService := service.NewAgentMetricService(gaugeAgentStorage, counterAgentStorage)
	worker := workers.NewAgentWorker(agentMetricService, cfg)
	if err = worker.Run(); err != nil {
		return fmt.Errorf("server has failed: %w", err)
	}
	return nil
}
