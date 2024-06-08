package database

import (
	"context"
	"errors"
	"fmt"
	"metrics/internal/server/core/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type MetricStorage struct {
	db *sqlx.DB
}

func NewStorage(cfg *Config) (*MetricStorage, error) {
	db, err := sqlx.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &MetricStorage{db: db}, nil
}

func (s *MetricStorage) GetMetric(mType, mName string) (*domain.Metric, error) {
	return nil, ErrNotImplemented
}

func (s *MetricStorage) SetMetric(m *domain.Metric) (*domain.Metric, error) {
	return nil, ErrNotImplemented
}

func (s *MetricStorage) GetAllMetrics() (domain.MetricsList, error) {
	return nil, ErrNotImplemented
}

func (s *MetricStorage) Ping() error {
	ctx := context.Background()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
