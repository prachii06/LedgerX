package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prachii06/LedgerX/internal/config"
)

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN())
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
