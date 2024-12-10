package sqlDb

import (
	"database/sql"

	"github.com/wizenheimer/iris/internal/config"
)

func Init(cfg config.DatabaseConfig) (*sql.DB, error) {
	// initialize db connection
	return sql.Open(cfg.Driver, prepareDSN(cfg))
}

func prepareDSN(_ config.DatabaseConfig) string {
	// prepare dsn
	return ""
}
