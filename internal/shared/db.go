package shared

import (
	"context"
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
)

func RunDBTransaction[T interface{}](
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	f func(queries *database.Queries) (T, error),
) (T, error) {
	tx, err := resourceConfig.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return *new(T), err
	}

	defer tx.Rollback()

	qtx := resourceConfig.DB.WithTx(tx)

	result, err := f(qtx)
	if err != nil {
		return *new(T), err
	}

	if err := tx.Commit(); err != nil {
		return *new(T), err
	}

	return result, nil
}

func RunNoResultDBTransaction(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	f func(queries *database.Queries) error,
) error {
	tx, err := resourceConfig.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := resourceConfig.DB.WithTx(tx)

	if err := f(qtx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func IsDBErrorNotFound(err error) bool {
	return err == sql.ErrNoRows
}
