package database

import (
	"context"
	"errors"
	"fmt"
	"metrics/internal/server/config"
	"metrics/internal/server/core/domain"

	"github.com/avast/retry-go"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProxyStorage struct {
	s *MetricStorage
}

func NewProxyPostgresStorage(cfg *Config) (*ProxyStorage, error) {
	storage, err := NewPostgresStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &ProxyStorage{
		s: storage,
	}, nil
}

func (px *ProxyStorage) GetMetric(ctx context.Context, mType, mName string) (*domain.Metric, error) {
	var metric *domain.Metric
	err := retry.Do(
		func() error {
			var err error
			metric, err = px.s.GetMetric(ctx, mType, mName)
			return err
		},
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return true
			}
			return false
		}),
		retry.Attempts(config.Attempts),
		retry.DelayType(config.DelayType),
		retry.OnRetry(config.OnRetry),
	)
	return metric, fmt.Errorf("%w", err)
}

func (px *ProxyStorage) SetMetric(ctx context.Context, m *domain.Metric) (*domain.Metric, error) {
	var metric *domain.Metric
	err := retry.Do(
		func() error {
			var err error
			metric, err = px.s.SetMetric(ctx, m)
			return err
		},
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return true
			}
			return false
		}),
		retry.Attempts(config.Attempts),
		retry.DelayType(config.DelayType),
		retry.OnRetry(config.OnRetry),
	)
	return metric, fmt.Errorf("%w", err)
}

func (px *ProxyStorage) SetMetrics(ctx context.Context, metrics domain.MetricsList) (domain.MetricsList, error) {
	var metricsOut domain.MetricsList
	err := retry.Do(
		func() error {
			var err error
			metricsOut, err = px.s.SetMetrics(ctx, metrics)
			return err
		},
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return true
			}
			return false
		}),
		retry.Attempts(config.Attempts),
		retry.DelayType(config.DelayType),
		retry.OnRetry(config.OnRetry),
	)
	return metricsOut, fmt.Errorf("%w", err)
}

func (px *ProxyStorage) GetAllMetrics(ctx context.Context) (domain.MetricsList, error) {
	var metrics domain.MetricsList
	err := retry.Do(
		func() error {
			var err error
			metrics, err = px.s.GetAllMetrics(ctx)
			return err
		},
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return true
			}
			return false
		}),
		retry.Attempts(config.Attempts),
		retry.DelayType(config.DelayType),
		retry.OnRetry(config.OnRetry),
	)
	return metrics, fmt.Errorf("%w", err)
}

func (px *ProxyStorage) Ping(ctx context.Context) error {
	return px.s.Ping(ctx)
}
