package database

import (
	"context"
	"errors"
	"metrics/internal/server/config"
	"metrics/internal/server/logger"

	"github.com/avast/retry-go"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (s *MetricStorage) retryExecRequest(ctx context.Context, query string, args ...any) error {
	var originalErr error
	err := retry.Do(
		func() error {
			_, originalErr := s.db.ExecContext(
				ctx,
				query,
				args...,
			)
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
		logger.Log.Error("retryError", zap.Error(err), zap.Error(originalErr))
		return originalErr
	}
	return nil
}
