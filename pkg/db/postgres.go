package db

import (
	"context"
	"database/sql"

	"scm.x5.ru/x5m/go-backend/packages/zlogger"
	"scm.x5.ru/x5m/go-backend/template/internal/config"

	// postgres driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// NewPostgres creates new Postgres connection
func NewPostgres(ctx context.Context, cfg *config.Common) (*sql.DB, error) {
	logger := zlogger.LoadOrCreateFromCtx(ctx)

	logger.Info().Msgf("connect from %q", cfg.Environment)
	conn, err := sql.Open("postgres", cfg.PostgreSQL.BuildDSN())
	if err != nil {
		return nil, errors.Wrap(err, "connect to db")
	}

	if err := conn.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping db")
	}

	return conn, nil
}
