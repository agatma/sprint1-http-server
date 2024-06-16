package database

import (
	"context"
	"errors"
	"fmt"
	"metrics/internal/server/config"
	"metrics/internal/server/core/domain"
	"metrics/internal/server/logger"

	"github.com/avast/retry-go"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const retryError = "retry error"

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
	var (
		metric      *domain.Metric
		originalErr error
	)
	err := retry.Do(
		func() error {
			metric, originalErr = px.s.GetMetric(ctx, mType, mName)
			return originalErr
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
	if err != nil {
		logger.Log.Error(retryError, zap.Error(err), zap.Error(originalErr))
		return metric, originalErr
	}
	return metric, nil
}

func (px *ProxyStorage) SetMetric(ctx context.Context, m *domain.Metric) (*domain.Metric, error) {
	var (
		metric      *domain.Metric
		originalErr error
	)
	err := retry.Do(
		func() error {
			metric, originalErr = px.s.SetMetric(ctx, m)
			return originalErr
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
	if err != nil {
		logger.Log.Error(retryError, zap.Error(err), zap.Error(originalErr))
		return metric, originalErr
	}
	return metric, nil
}

func (px *ProxyStorage) SetMetrics(ctx context.Context, metrics domain.MetricsList) (domain.MetricsList, error) {
	var (
		metricsOut  domain.MetricsList
		originalErr error
	)
	err := retry.Do(
		func() error {
			metricsOut, originalErr = px.s.SetMetrics(ctx, metrics)
			return originalErr
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
	if err != nil {
		logger.Log.Error(retryError, zap.Error(err), zap.Error(originalErr))
		return metricsOut, originalErr
	}
	return metricsOut, nil
}

func (px *ProxyStorage) GetAllMetrics(ctx context.Context) (domain.MetricsList, error) {
	var (
		metrics     domain.MetricsList
		originalErr error
	)
	err := retry.Do(
		func() error {
			metrics, originalErr = px.s.GetAllMetrics(ctx)
			return originalErr
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
	if err != nil {
		logger.Log.Error(retryError, zap.Error(err), zap.Error(originalErr))
		return metrics, originalErr
	}
	return metrics, nil
}

func (px *ProxyStorage) Ping(ctx context.Context) error {
	var originalErr error
	err := retry.Do(
		func() error {
			err := px.s.Ping(ctx)
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
	if err != nil {
		logger.Log.Error(retryError, zap.Error(err), zap.Error(originalErr))
		return originalErr
	}
	return nil
}
