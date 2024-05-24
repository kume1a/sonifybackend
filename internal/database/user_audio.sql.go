// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user_audio.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const countUserAudioByLocalID = `-- name: CountUserAudioByLocalID :one
SELECT COUNT(*) FROM user_audios 
  INNER JOIN audios ON user_audios.audio_id = audios.id 
  WHERE user_audios.user_id = $1 AND audios.local_id = $2
`

type CountUserAudioByLocalIDParams struct {
	UserID  uuid.UUID
	LocalID sql.NullString
}

func (q *Queries) CountUserAudioByLocalID(ctx context.Context, arg CountUserAudioByLocalIDParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUserAudioByLocalID, arg.UserID, arg.LocalID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUserAudio = `-- name: CreateUserAudio :one
INSERT INTO user_audios(
  created_at,
  user_id, 
  audio_id
) VALUES ($1,$2,$3) 
RETURNING id, created_at, user_id, audio_id
`

type CreateUserAudioParams struct {
	CreatedAt time.Time
	UserID    uuid.UUID
	AudioID   uuid.UUID
}

func (q *Queries) CreateUserAudio(ctx context.Context, arg CreateUserAudioParams) (UserAudio, error) {
	row := q.db.QueryRowContext(ctx, createUserAudio, arg.CreatedAt, arg.UserID, arg.AudioID)
	var i UserAudio
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.AudioID,
	)
	return i, err
}

const getUserAudioByVideoID = `-- name: GetUserAudioByVideoID :one
SELECT user_audios.id, user_audios.created_at, user_id, audio_id, audios.id, audios.created_at, title, author, duration_ms, path, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id FROM user_audios
  INNER JOIN audios ON user_audios.audio_id = audios.id
  WHERE user_audios.user_id = $1 AND audios.youtube_video_id = $2
`

type GetUserAudioByVideoIDParams struct {
	UserID         uuid.UUID
	YoutubeVideoID sql.NullString
}

type GetUserAudioByVideoIDRow struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UserID         uuid.UUID
	AudioID        uuid.UUID
	ID_2           uuid.UUID
	CreatedAt_2    time.Time
	Title          sql.NullString
	Author         sql.NullString
	DurationMs     sql.NullInt32
	Path           sql.NullString
	SizeBytes      sql.NullInt64
	YoutubeVideoID sql.NullString
	ThumbnailPath  sql.NullString
	SpotifyID      sql.NullString
	ThumbnailUrl   sql.NullString
	LocalID        sql.NullString
}

func (q *Queries) GetUserAudioByVideoID(ctx context.Context, arg GetUserAudioByVideoIDParams) (GetUserAudioByVideoIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserAudioByVideoID, arg.UserID, arg.YoutubeVideoID)
	var i GetUserAudioByVideoIDRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.AudioID,
		&i.ID_2,
		&i.CreatedAt_2,
		&i.Title,
		&i.Author,
		&i.DurationMs,
		&i.Path,
		&i.SizeBytes,
		&i.YoutubeVideoID,
		&i.ThumbnailPath,
		&i.SpotifyID,
		&i.ThumbnailUrl,
		&i.LocalID,
	)
	return i, err
}

const getUserAudioIDs = `-- name: GetUserAudioIDs :many
SELECT audio_id FROM user_audios WHERE user_id = $1
`

func (q *Queries) GetUserAudioIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getUserAudioIDs, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var audio_id uuid.UUID
		if err := rows.Scan(&audio_id); err != nil {
			return nil, err
		}
		items = append(items, audio_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserAudiosByAudioIds = `-- name: GetUserAudiosByAudioIds :many
SELECT user_audios.id, user_audios.created_at, user_audios.user_id, user_audios.audio_id,
  audio_likes.user_id as audio_likes_user_id,
  audio_likes.audio_id as audio_likes_audio_id,
  audios.id as audio_id,
  audios.created_at as audio_created_at,
  audios.title as audio_title,
  audios.author as audio_author,
  audios.duration_ms as audio_duration_ms,
  audios.path as audio_path,
  audios.size_bytes as audio_size_bytes,
  audios.youtube_video_id as audio_youtube_video_id,
  audios.thumbnail_path as audio_thumbnail_path,
  audios.spotify_id as audio_spotify_id,
  audios.thumbnail_url as audio_thumbnail_url,
  audios.local_id as audio_local_id
FROM user_audios
INNER JOIN audios ON user_audios.audio_id = audios.id
LEFT JOIN audio_likes ON audio_likes.audio_id = audios.id
WHERE user_audios.user_id = $1 AND audios.id = ANY($2::uuid[])
`

type GetUserAudiosByAudioIdsParams struct {
	UserID   uuid.UUID
	AudioIds []uuid.UUID
}

type GetUserAudiosByAudioIdsRow struct {
	ID                  uuid.UUID
	CreatedAt           time.Time
	UserID              uuid.UUID
	AudioID             uuid.UUID
	AudioLikesUserID    uuid.NullUUID
	AudioLikesAudioID   uuid.NullUUID
	AudioID_2           uuid.UUID
	AudioCreatedAt      time.Time
	AudioTitle          sql.NullString
	AudioAuthor         sql.NullString
	AudioDurationMs     sql.NullInt32
	AudioPath           sql.NullString
	AudioSizeBytes      sql.NullInt64
	AudioYoutubeVideoID sql.NullString
	AudioThumbnailPath  sql.NullString
	AudioSpotifyID      sql.NullString
	AudioThumbnailUrl   sql.NullString
	AudioLocalID        sql.NullString
}

func (q *Queries) GetUserAudiosByAudioIds(ctx context.Context, arg GetUserAudiosByAudioIdsParams) ([]GetUserAudiosByAudioIdsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserAudiosByAudioIds, arg.UserID, pq.Array(arg.AudioIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserAudiosByAudioIdsRow
	for rows.Next() {
		var i GetUserAudiosByAudioIdsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UserID,
			&i.AudioID,
			&i.AudioLikesUserID,
			&i.AudioLikesAudioID,
			&i.AudioID_2,
			&i.AudioCreatedAt,
			&i.AudioTitle,
			&i.AudioAuthor,
			&i.AudioDurationMs,
			&i.AudioPath,
			&i.AudioSizeBytes,
			&i.AudioYoutubeVideoID,
			&i.AudioThumbnailPath,
			&i.AudioSpotifyID,
			&i.AudioThumbnailUrl,
			&i.AudioLocalID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
