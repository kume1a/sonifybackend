package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateUser(db *database.Queries, ctx context.Context, params *database.CreateUserParams) (*database.User, error) {
	createdAt := time.Now().UTC()

	user, err := db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Name:      params.Name,
		Email:     params.Email,
	})

	if err != nil {
		log.Println("Error creating user: ", err)
	}

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
