package usersync

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func GetOrCreateUserSyncDatumByUserId(
	ctx context.Context,
	db *database.Queries,
	userID uuid.UUID,
) (*database.UserSyncDatum, *shared.HttpError) {
	entity, err := db.GetUserSyncDatumByUserID(ctx, userID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		entity, err = db.CreateUserSyncDatum(ctx, database.CreateUserSyncDatumParams{
			ID:                    uuid.New(),
			UserID:                userID,
			SpotifyLastSyncedAt:   sql.NullTime{},
			UserAudioLastSyncedAt: sql.NullTime{},
		})

		if err != nil {
			log.Println("Error creating user sync datum:", err)
			return nil, shared.InternalServerErrorDef()
		}

		return &entity, nil
	} else if err != nil {
		return nil, shared.InternalServerErrorDef()
	}

	return &entity, nil
}

func UpsertUserSyncDatumByUserId(
	ctx context.Context,
	db *database.Queries,
	params *database.UserSyncDatum,
) (*database.UserSyncDatum, *shared.HttpError) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}

	entity, err := db.GetUserSyncDatumByUserID(ctx, params.UserID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		entity, err = db.CreateUserSyncDatum(ctx, database.CreateUserSyncDatumParams{
			ID:                    params.ID,
			UserID:                params.UserID,
			SpotifyLastSyncedAt:   params.SpotifyLastSyncedAt,
			UserAudioLastSyncedAt: params.UserAudioLastSyncedAt,
		})

		if err != nil {
			return nil, shared.InternalServerErrorDef()
		}
	} else if err != nil {
		log.Println("Error fetching user sync datum:", err)
		return nil, shared.InternalServerErrorDef()
	}

	entity, err = db.UpdateUserSyncDatumByUserID(ctx, database.UpdateUserSyncDatumByUserIDParams{
		UserID:                params.UserID,
		SpotifyLastSyncedAt:   params.SpotifyLastSyncedAt,
		UserAudioLastSyncedAt: params.UserAudioLastSyncedAt,
	})
	if err != nil {
		log.Println("Error updating user sync datum:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &entity, nil
}
