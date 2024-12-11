package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/wizenheimer/iris/src/internal/config"
)

func Init(cfg *config.Config) (*sql.DB, error) {

	connString := prepareConnectionString(cfg)

	// initialize db connection
	return sql.Open(cfg.Database.Driver, connString)
}

func prepareConnectionString(cfg *config.Config) string {
	if cfg.Database.ConnectionString != "" {
		return cfg.Database.ConnectionString
	}

	// Determine the environment
	var sslMode string
	switch cfg.Environment.EnvProfile {
	case "development":
		// Set up development environment
		sslMode = "disable"
	case "production":
		// Set up production environment
		sslMode = "require"
	default:
		// Set up default environment
		sslMode = "disable"
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
