package usersync

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateUserSyncData(
	ctx context.Context,
	db *database.Queries,
	params database.CreateUserSyncDataParams,
) (*database.UserSyncDatum, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	entity, err := db.CreateUserSyncData(ctx, params)

	if err != nil {
		log.Println("Error creating user: ", err)
	}

	return &entity, err
}

func GetUserSyncDatumByUserId(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (*database.UserSyncDatum, error) {
	entity, err := db.GetUserSyncDatumByUserId(ctx, userId)

	if err != nil {
		log.Println("Error getting user sync datum: ", err)
	}

	return &entity, err
}

func UpdateUserSyncDatumById(
	ctx context.Context,
	db *database.Queries,
	params database.UpdateUserSyncDatumByUserIdParams,
) (*database.UserSyncDatum, error) {
	entity, err := db.UpdateUserSyncDatumByUserId(ctx, params)

	if err != nil {
		log.Println("Error updating user sync datum: ", err)
	}

	return &entity, err
}
