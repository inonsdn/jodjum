package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXHandler struct {
	con *pgxpool.Pool
}

func NewPGX(url string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), url)
}
