package usersync

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func createUserSyncData(
	ctx context.Context,
	db *database.Queries,
	params database.CreateUserSyncDatumParams,
) (*database.UserSyncDatum, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	entity, err := db.CreateUserSyncDatum(ctx, params)

	if err != nil {
		log.Println("Error creating user: ", err)
	}

	return &entity, err
}

func getUserSyncDatumByUserId(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (*database.UserSyncDatum, error) {
	entity, err := db.GetUserSyncDatumByUserID(ctx, userId)

	if err != nil {
		log.Println("Error getting user sync datum: ", err)
	}

	return &entity, err
}

func updateUserSyncDatumByUserId(
	ctx context.Context,
	db *database.Queries,
	params database.UpdateUserSyncDatumByUserIDParams,
) (*database.UserSyncDatum, error) {
	entity, err := db.UpdateUserSyncDatumByUserID(ctx, params)

	if err != nil {
		log.Println("Error updating user sync datum: ", err)
	}

	return &entity, err
}
