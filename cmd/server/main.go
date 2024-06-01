package main

import (
	"fmt"
	"log"

	"github.com/agatma/sprint1-http-server/internal/server/adapters/api/rest"
	"github.com/agatma/sprint1-http-server/internal/server/adapters/storage"
	"github.com/agatma/sprint1-http-server/internal/server/adapters/storage/memory"
	"github.com/agatma/sprint1-http-server/internal/server/config"
	"github.com/agatma/sprint1-http-server/internal/server/core/service"
	"github.com/agatma/sprint1-http-server/internal/server/logger"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}
	if err = logger.Initialize(cfg.LogLevel); err != nil {
		return fmt.Errorf("can't load logger: %w", err)
	}
	metricStorage, err := storage.NewStorage(storage.Config{
		Memory: &memory.Config{},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize a storage: %w", err)
	}
	metricService := service.NewMetricService(metricStorage)
	api := rest.NewAPI(metricService, cfg)
	if err := api.Run(); err != nil {
		return fmt.Errorf("server has failed: %w", err)
	}
	return nil
}
