package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func NewDB(ctx context.Context, logger *logrus.Logger, cfg *Config) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, cfg.PgURL)
	if err != nil {
		logger.Fatal(err)
	}
	return pool
}
