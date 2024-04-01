package shared

import (
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/database"
)

type ApiConfg struct {
	DB    *database.Queries
	SqlDB *sql.DB
}
