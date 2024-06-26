package config

import (
	"database/sql"

	"github.com/gocraft/work"
	"github.com/kume1a/sonifybackend/internal/database"
)

type ResourceConfig struct {
	DB    *database.Queries
	SqlDB *sql.DB
}

type ApiConfig struct {
	*ResourceConfig
	WorkEnqueuer *work.Enqueuer
}
