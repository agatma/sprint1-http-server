package database

import (
	"embed"
	"fmt"
	"metrics/internal/server/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations
var migrations embed.FS

func migrate(db *sqlx.DB, version int64) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres migrate set dialect postgres: %w", err)
	}

	if err := goose.UpTo(db.DB, "migrations", version); err != nil {
		return fmt.Errorf("postgres migrate up: %w", err)
	}
	logger.Log.Info("successful migrations", zap.Int64("version", version))
	return nil
}
