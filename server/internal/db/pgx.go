package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXHandler struct {
	con *pgxpool.Pool
}

func NewPGX(url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	// Cap connections per instance so that on autoscaling platforms
	// (e.g. Cloud Run) N instances don't exhaust Supabase's connection limit.
	cfg.MaxConns = 5

	return pgxpool.NewWithConfig(context.Background(), cfg)
}
