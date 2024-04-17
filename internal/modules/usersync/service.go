package usersync

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func GetOrCreateUserSyncDatumByUserId(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (*database.UserSyncDatum, *shared.HttpError) {
	entity, err := getUserSyncDatumByUserId(ctx, db, userId)

	if err != nil && shared.IsDBErrorNotFound(err) {
		entity, err = createUserSyncData(ctx, db, database.CreateUserSyncDatumParams{
			UserID:                userId,
			SpotifyLastSyncedAt:   sql.NullTime{},
			UserAudioLastSyncedAt: sql.NullTime{},
		})
		if err != nil {
			return nil, shared.HttpErrInternalServerErrorDef()
		}

		return entity, nil
	}

	if err != nil {
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return entity, nil
}

func UpdateUserSyncDatumByUserId(
	ctx context.Context,
	db *database.Queries,
	params database.UpdateUserSyncDatumByUserIdParams,
) (*database.UserSyncDatum, *shared.HttpError) {
	entity, err := updateUserSyncDatumByUserId(ctx, db, params)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return nil, shared.HttpErrNotFound(shared.ErrUserSyncDatumNotFound)
	}

	if err != nil {
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return entity, nil
}
