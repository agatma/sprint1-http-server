package workers

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

const (
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	LogLevel       string
}

func NewConfig() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "run address")
	flag.IntVar(&cfg.PollInterval, "p", defaultPollInterval, " poll interval ")
	flag.IntVar(&cfg.ReportInterval, "r", defaultReportInterval, " report interval ")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to get config for worker: %w", err)
	}
	return &cfg, nil
}
