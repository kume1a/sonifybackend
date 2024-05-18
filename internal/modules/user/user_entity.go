package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateUser(
	ctx context.Context,
	db *database.Queries,
	params database.CreateUserParams,
) (*database.User, error) {
	createdAt := time.Now().UTC()

	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = createdAt
	}
	user, err := db.CreateUser(ctx, params)

	if err != nil {
		log.Println("Error creating user: ", err)
	}

	return &user, err
}

func GetUserByID(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (*database.User, error) {
	user, err := db.GetUserByID(ctx, userId)

	log.Println("Error getting user by ID: ", err)

	return &user, err
}

func GetUserByEmail(ctx context.Context, db *database.Queries, email string) (*database.User, error) {
	user, err := db.GetUserByEmail(ctx, sql.NullString{
		String: email,
		Valid:  true,
	})

	if err != nil {
		log.Println("Error getting user by email: ", err)
	}

	return &user, err
}

func UpdateUser(
	ctx context.Context,
	db *database.Queries,
	params *database.UpdateUserParams,
) (*database.User, error) {
	user, err := db.UpdateUser(ctx, *params)

	if err != nil {
		log.Println("Error updating user: ", err)
	}

	return &user, err
}

func UserExistsByEmail(
	ctx context.Context,
	db *database.Queries,
	email string,
) (bool, error) {
	count, err := db.CountUsersByEmail(ctx, sql.NullString{String: email, Valid: true})

	if err != nil {
		log.Println("Error counting users by email: ", err)
	}

	return count > 0, err
}
