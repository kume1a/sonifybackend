// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: audio.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const countUserAudioByLocalId = `-- name: CountUserAudioByLocalId :one
SELECT COUNT(*) FROM user_audios 
  INNER JOIN audio ON user_audios.audio_id = audio.id 
  WHERE user_audios.user_id = $1 AND audio.local_id = $2
`

type CountUserAudioByLocalIdParams struct {
	UserID  uuid.UUID
	LocalID sql.NullString
}

func (q *Queries) CountUserAudioByLocalId(ctx context.Context, arg CountUserAudioByLocalIdParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUserAudioByLocalId, arg.UserID, arg.LocalID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAudio = `-- name: CreateAudio :one
INSERT INTO audio(
  id, 
  created_at,
  title,
  author,
  duration_ms,
  path,
  size_bytes,
  youtube_video_id,
  thumbnail_path,
  spotify_id,
  thumbnail_url,
  local_id
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id, title, author, duration_ms, path, created_at, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id
`

type CreateAudioParams struct {
	ID             uuid.UUID
	CreatedAt      time.Time
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

func (q *Queries) CreateAudio(ctx context.Context, arg CreateAudioParams) (Audio, error) {
	row := q.db.QueryRowContext(ctx, createAudio,
		arg.ID,
		arg.CreatedAt,
		arg.Title,
		arg.Author,
		arg.DurationMs,
		arg.Path,
		arg.SizeBytes,
		arg.YoutubeVideoID,
		arg.ThumbnailPath,
		arg.SpotifyID,
		arg.ThumbnailUrl,
		arg.LocalID,
	)
	var i Audio
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.DurationMs,
		&i.Path,
		&i.CreatedAt,
		&i.SizeBytes,
		&i.YoutubeVideoID,
		&i.ThumbnailPath,
		&i.SpotifyID,
		&i.ThumbnailUrl,
		&i.LocalID,
	)
	return i, err
}

const createUserAudio = `-- name: CreateUserAudio :one
INSERT INTO user_audios(
  created_at,
  user_id, 
  audio_id
) VALUES ($1,$2,$3) RETURNING user_id, audio_id, created_at
`

type CreateUserAudioParams struct {
	CreatedAt time.Time
	UserID    uuid.UUID
	AudioID   uuid.UUID
}

func (q *Queries) CreateUserAudio(ctx context.Context, arg CreateUserAudioParams) (UserAudio, error) {
	row := q.db.QueryRowContext(ctx, createUserAudio, arg.CreatedAt, arg.UserID, arg.AudioID)
	var i UserAudio
	err := row.Scan(&i.UserID, &i.AudioID, &i.CreatedAt)
	return i, err
}

const deleteAudioById = `-- name: DeleteAudioById :exec
DELETE FROM audio WHERE id = $1
`

func (q *Queries) DeleteAudioById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteAudioById, id)
	return err
}

const getAudioById = `-- name: GetAudioById :one
SELECT id, title, author, duration_ms, path, created_at, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id FROM audio WHERE id = $1
`

func (q *Queries) GetAudioById(ctx context.Context, id uuid.UUID) (Audio, error) {
	row := q.db.QueryRowContext(ctx, getAudioById, id)
	var i Audio
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.DurationMs,
		&i.Path,
		&i.CreatedAt,
		&i.SizeBytes,
		&i.YoutubeVideoID,
		&i.ThumbnailPath,
		&i.SpotifyID,
		&i.ThumbnailUrl,
		&i.LocalID,
	)
	return i, err
}

const getAudioIdsBySpotifyIds = `-- name: GetAudioIdsBySpotifyIds :many
SELECT 
  id
FROM audio
WHERE spotify_id = ANY($1::text[])
`

func (q *Queries) GetAudioIdsBySpotifyIds(ctx context.Context, spotifyIds []string) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getAudioIdsBySpotifyIds, pq.Array(spotifyIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAudioSpotifyIdsBySpotifyIds = `-- name: GetAudioSpotifyIdsBySpotifyIds :many
SELECT 
  id, spotify_id 
FROM audio 
WHERE spotify_id = ANY($1::text[])
`

type GetAudioSpotifyIdsBySpotifyIdsRow struct {
	ID        uuid.UUID
	SpotifyID sql.NullString
}

func (q *Queries) GetAudioSpotifyIdsBySpotifyIds(ctx context.Context, spotifyIds []string) ([]GetAudioSpotifyIdsBySpotifyIdsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAudioSpotifyIdsBySpotifyIds, pq.Array(spotifyIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAudioSpotifyIdsBySpotifyIdsRow
	for rows.Next() {
		var i GetAudioSpotifyIdsBySpotifyIdsRow
		if err := rows.Scan(&i.ID, &i.SpotifyID); err != nil {
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

const getAudiosByIds = `-- name: GetAudiosByIds :many
SELECT id, title, author, duration_ms, path, created_at, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id FROM audio WHERE id = ANY($1::uuid[])
`

func (q *Queries) GetAudiosByIds(ctx context.Context, ids []uuid.UUID) ([]Audio, error) {
	rows, err := q.db.QueryContext(ctx, getAudiosByIds, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Audio
	for rows.Next() {
		var i Audio
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.DurationMs,
			&i.Path,
			&i.CreatedAt,
			&i.SizeBytes,
			&i.YoutubeVideoID,
			&i.ThumbnailPath,
			&i.SpotifyID,
			&i.ThumbnailUrl,
			&i.LocalID,
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

const getAudiosByUserId = `-- name: GetAudiosByUserId :many
SELECT 
  audio.id, audio.title, audio.author, audio.duration_ms, audio.path, audio.created_at, audio.size_bytes, audio.youtube_video_id, audio.thumbnail_path, audio.spotify_id, audio.thumbnail_url, audio.local_id
  FROM user_audios
  INNER JOIN audio ON user_audios.audio_id = audio.id
  WHERE user_id = $1
`

func (q *Queries) GetAudiosByUserId(ctx context.Context, userID uuid.UUID) ([]Audio, error) {
	rows, err := q.db.QueryContext(ctx, getAudiosByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Audio
	for rows.Next() {
		var i Audio
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.DurationMs,
			&i.Path,
			&i.CreatedAt,
			&i.SizeBytes,
			&i.YoutubeVideoID,
			&i.ThumbnailPath,
			&i.SpotifyID,
			&i.ThumbnailUrl,
			&i.LocalID,
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

const getPlaylistAudiosBySpotifyIds = `-- name: GetPlaylistAudiosBySpotifyIds :many
SELECT 
  audio.id, audio.title, audio.author, audio.duration_ms, audio.path, audio.created_at, audio.size_bytes, audio.youtube_video_id, audio.thumbnail_path, audio.spotify_id, audio.thumbnail_url, audio.local_id
  FROM playlist_audios
  INNER JOIN audio ON playlist_audios.audio_id = audio.id
  WHERE playlist_audios.playlist_id = $1 AND audio.spotify_id = ANY($2::text[])
`

type GetPlaylistAudiosBySpotifyIdsParams struct {
	PlaylistID uuid.UUID
	SpotifyIds []string
}

func (q *Queries) GetPlaylistAudiosBySpotifyIds(ctx context.Context, arg GetPlaylistAudiosBySpotifyIdsParams) ([]Audio, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylistAudiosBySpotifyIds, arg.PlaylistID, pq.Array(arg.SpotifyIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Audio
	for rows.Next() {
		var i Audio
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.DurationMs,
			&i.Path,
			&i.CreatedAt,
			&i.SizeBytes,
			&i.YoutubeVideoID,
			&i.ThumbnailPath,
			&i.SpotifyID,
			&i.ThumbnailUrl,
			&i.LocalID,
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

const getUserAudioByVideoId = `-- name: GetUserAudioByVideoId :one
SELECT user_id, audio_id, user_audios.created_at, id, title, author, duration_ms, path, audio.created_at, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id FROM user_audios
  INNER JOIN audio ON user_audios.audio_id = audio.id
  WHERE user_audios.user_id = $1 AND audio.youtube_video_id = $2
`

type GetUserAudioByVideoIdParams struct {
	UserID         uuid.UUID
	YoutubeVideoID sql.NullString
}

type GetUserAudioByVideoIdRow struct {
	UserID         uuid.UUID
	AudioID        uuid.UUID
	CreatedAt      time.Time
	ID             uuid.UUID
	Title          sql.NullString
	Author         sql.NullString
	DurationMs     sql.NullInt32
	Path           sql.NullString
	CreatedAt_2    time.Time
	SizeBytes      sql.NullInt64
	YoutubeVideoID sql.NullString
	ThumbnailPath  sql.NullString
	SpotifyID      sql.NullString
	ThumbnailUrl   sql.NullString
	LocalID        sql.NullString
}

func (q *Queries) GetUserAudioByVideoId(ctx context.Context, arg GetUserAudioByVideoIdParams) (GetUserAudioByVideoIdRow, error) {
	row := q.db.QueryRowContext(ctx, getUserAudioByVideoId, arg.UserID, arg.YoutubeVideoID)
	var i GetUserAudioByVideoIdRow
	err := row.Scan(
		&i.UserID,
		&i.AudioID,
		&i.CreatedAt,
		&i.ID,
		&i.Title,
		&i.Author,
		&i.DurationMs,
		&i.Path,
		&i.CreatedAt_2,
		&i.SizeBytes,
		&i.YoutubeVideoID,
		&i.ThumbnailPath,
		&i.SpotifyID,
		&i.ThumbnailUrl,
		&i.LocalID,
	)
	return i, err
}

const getUserAudioIds = `-- name: GetUserAudioIds :many
SELECT audio_id FROM user_audios WHERE user_id = $1
`

func (q *Queries) GetUserAudioIds(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getUserAudioIds, userID)
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

const updateAudio = `-- name: UpdateAudio :one
UPDATE audio SET 
  title = $1, 
  author = $2, 
  duration_ms = $3, 
  path = $4, 
  thumbnail_path=$5 
WHERE id = $6 
RETURNING id, title, author, duration_ms, path, created_at, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id
`

type UpdateAudioParams struct {
	Title         sql.NullString
	Author        sql.NullString
	DurationMs    sql.NullInt32
	Path          sql.NullString
	ThumbnailPath sql.NullString
	ID            uuid.UUID
}

func (q *Queries) UpdateAudio(ctx context.Context, arg UpdateAudioParams) (Audio, error) {
	row := q.db.QueryRowContext(ctx, updateAudio,
		arg.Title,
		arg.Author,
		arg.DurationMs,
		arg.Path,
		arg.ThumbnailPath,
		arg.ID,
	)
	var i Audio
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.DurationMs,
		&i.Path,
		&i.CreatedAt,
		&i.SizeBytes,
		&i.YoutubeVideoID,
		&i.ThumbnailPath,
		&i.SpotifyID,
		&i.ThumbnailUrl,
		&i.LocalID,
	)
	return i, err
}
