package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/agatma/sprint1-http-server/internal/server/adapters/api/rest"
	"github.com/agatma/sprint1-http-server/internal/server/adapters/storage"
	"github.com/agatma/sprint1-http-server/internal/server/adapters/storage/file"
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
	metricStorage, err := initMetricStorage(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize a storage: %w", err)
	}
	metricService, err := service.NewMetricService(cfg, metricStorage)
	if err != nil {
		return fmt.Errorf("failed to initialize a service: %w", err)
	}
	api := rest.NewAPI(metricService, cfg)
	if err = api.Run(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			err = metricService.SaveMetricsToFile()
			if err != nil {
				return fmt.Errorf("failed to save metrics during shutdown: %w", err)
			}
			logger.Log.Info("metrics are saved to file")
		}

		return fmt.Errorf("server has failed: %w", err)
	}
	return nil
}

func initMetricStorage(cfg *config.Config) (storage.MetricStorage, error) {
	if cfg.StoreInterval == 0 {
		metricStorage, err := storage.NewStorage(storage.Config{
			File: &file.Config{
				Filepath: cfg.FileStoragePath,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to init file storage %w", err)
		}
		logger.Log.Info("initialize file storage")
		return metricStorage, nil
	} else {
		metricStorage, err := storage.NewStorage(storage.Config{
			Memory: &memory.Config{},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to init memory storage %w", err)
		}
		logger.Log.Info("initialize memory storage")
		return metricStorage, nil
	}
}
