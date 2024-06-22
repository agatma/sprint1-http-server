package retry

import (
	"context"
	"errors"
	"fmt"
	"metrics/internal/server/logger"
	"time"

	"github.com/avast/retry-go"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const addTime = 2

const Attempts uint = 3

func DelayType(n uint, _ error, config *retry.Config) time.Duration {
	switch n {
	case 0:
		return 1 * time.Second
	case 1:
		return (1 + addTime) * time.Second
	default:
		return (1 + addTime + addTime) * time.Second
	}
}

func OnRetry(n uint, err error) {
	logger.Log.Error(fmt.Sprintf(`%d %s`, n, err.Error()))
}

func ExecContext(db *sqlx.DB, ctx context.Context, query string, args ...any) error {
	var originalErr error
	err := retry.Do(
		func() error {
			_, originalErr := db.ExecContext(
				ctx,
				query,
				args...,
			)
			if originalErr != nil {
				return fmt.Errorf("%w", originalErr)
			}
			return nil
		},
		retry.RetryIf(func(err error) bool {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return true
			}
			return false
		}),
		retry.Attempts(Attempts),
		retry.DelayType(DelayType),
		retry.OnRetry(OnRetry),
	)
	if err != nil {
		logger.Log.Error("retryError", zap.Error(err), zap.Error(originalErr))
		return fmt.Errorf("%w", originalErr)
	}
	return originalErr
}
