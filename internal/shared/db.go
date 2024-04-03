package shared

import (
	"context"

	"github.com/kume1a/sonifybackend/internal/database"
)

type DBOverrideOptions struct {
	OverrideID        bool `default:"true"`
	OverrideCreatedAt bool `default:"true"`
	OverrideUpdatedAt bool `default:"true"`
}

func RunDbTransaction[T interface{}](
	ctx context.Context,
	apiCfg *ApiConfg,
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
