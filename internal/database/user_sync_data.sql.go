// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user_sync_data.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUserSyncDatum = `-- name: CreateUserSyncDatum :one
INSERT INTO user_sync_data (
  id,
  user_id, 
  spotify_last_synced_at,
  user_audio_last_synced_at
) VALUES ($1,$2,$3,$4)
RETURNING id, user_id, spotify_last_synced_at, user_audio_last_synced_at
`

type CreateUserSyncDatumParams struct {
	ID                    uuid.UUID
	UserID                uuid.UUID
	SpotifyLastSyncedAt   sql.NullTime
	UserAudioLastSyncedAt sql.NullTime
}

func (q *Queries) CreateUserSyncDatum(ctx context.Context, arg CreateUserSyncDatumParams) (UserSyncDatum, error) {
	row := q.db.QueryRowContext(ctx, createUserSyncDatum,
		arg.ID,
		arg.UserID,
		arg.SpotifyLastSyncedAt,
		arg.UserAudioLastSyncedAt,
	)
	var i UserSyncDatum
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpotifyLastSyncedAt,
		&i.UserAudioLastSyncedAt,
	)
	return i, err
}

const getUserSyncDatumByUserID = `-- name: GetUserSyncDatumByUserID :one
SELECT id, user_id, spotify_last_synced_at, user_audio_last_synced_at FROM user_sync_data WHERE user_id = $1
`

func (q *Queries) GetUserSyncDatumByUserID(ctx context.Context, userID uuid.UUID) (UserSyncDatum, error) {
	row := q.db.QueryRowContext(ctx, getUserSyncDatumByUserID, userID)
	var i UserSyncDatum
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpotifyLastSyncedAt,
		&i.UserAudioLastSyncedAt,
	)
	return i, err
}

const updateUserSyncDatumByUserID = `-- name: UpdateUserSyncDatumByUserID :one
UPDATE user_sync_data
  SET spotify_last_synced_at = COALESCE($1, spotify_last_synced_at),
      user_audio_last_synced_at = COALESCE($2, user_audio_last_synced_at)
WHERE user_id = $3
RETURNING id, user_id, spotify_last_synced_at, user_audio_last_synced_at
`

type UpdateUserSyncDatumByUserIDParams struct {
	SpotifyLastSyncedAt   sql.NullTime
	UserAudioLastSyncedAt sql.NullTime
	UserID                uuid.UUID
}

func (q *Queries) UpdateUserSyncDatumByUserID(ctx context.Context, arg UpdateUserSyncDatumByUserIDParams) (UserSyncDatum, error) {
	row := q.db.QueryRowContext(ctx, updateUserSyncDatumByUserID, arg.SpotifyLastSyncedAt, arg.UserAudioLastSyncedAt, arg.UserID)
	var i UserSyncDatum
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpotifyLastSyncedAt,
		&i.UserAudioLastSyncedAt,
	)
	return i, err
}
