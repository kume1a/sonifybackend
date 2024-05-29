package shared

import (
	"context"
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
)

type DBOverrideOptions struct {
	OverrideID        bool `default:"true"`
	OverrideCreatedAt bool `default:"true"`
	OverrideUpdatedAt bool `default:"true"`
}

func RunDBTransaction[T interface{}](
	ctx context.Context,
	apiCfg *config.ApiConfig,
	f func(queries *database.Queries) (T, error),
) (T, error) {
	tx, err := apiCfg.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return *new(T), err
	}

	defer tx.Rollback()

	qtx := apiCfg.DB.WithTx(tx)

	result, err := f(qtx)
	if err != nil {
		return *new(T), err
	}

	if err := tx.Commit(); err != nil {
		return *new(T), err
	}

	return result, nil
}

func IsDBErrorNotFound(err error) bool {
	return err == sql.ErrNoRows
}
