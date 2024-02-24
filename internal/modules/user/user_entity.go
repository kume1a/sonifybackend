package user

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateUser(db *database.Queries, ctx context.Context, params *database.CreateUserParams) (*database.User, error) {
	user, err := db.CreateUser(ctx, *params)

	return &user, err
}

func GetUserByID(db *database.Queries, ctx context.Context, userId uuid.UUID) (*database.User, error) {
	user, err := db.GetUserById(ctx, userId)

	return &user, err
}

func GetUserByEmail(db *database.Queries, ctx context.Context, email string) (*database.User, error) {
	user, err := db.GetUserByEmail(ctx, sql.NullString{
		String: email,
		Valid:  true,
	})

	return &user, err
}

func UpdateUser(db *database.Queries, ctx context.Context, params *database.UpdateUserParams) (*database.User, error) {
	user, err := db.UpdateUser(ctx, *params)

	return &user, err
}
