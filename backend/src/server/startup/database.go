package startup

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizenheimer/byrd/src/internal/config"
)

func SetupDB(cfg *config.Config) (*pgxpool.Pool, error) {
	connString := prepareConnectionString(cfg)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return pool, nil
}

func prepareConnectionString(cfg *config.Config) string {
	if cfg.Database.ConnectionString != "" {
		return cfg.Database.ConnectionString
	}

	sslMode := "disable"
	if cfg.Environment.EnvProfile == "production" {
		sslMode = "verify-full"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
		sslMode,
	)
}
