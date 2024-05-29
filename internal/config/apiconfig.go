package config

import (
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/database"
)

type ApiConfig struct {
	DB    *database.Queries
	SqlDB *sql.DB
}
