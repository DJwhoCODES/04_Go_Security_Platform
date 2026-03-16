package database

import (
	"context"
	"fmt"
	"time"

	"github.com/djwhocodes/auth-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
