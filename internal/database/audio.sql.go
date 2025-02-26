// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: audio.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const audioExistsByYoutubeVideoID = `-- name: AudioExistsByYoutubeVideoID :one
SELECT EXISTS(
  SELECT 1 FROM audios
    WHERE audios.youtube_video_id = $1
)
`

func (q *Queries) AudioExistsByYoutubeVideoID(ctx context.Context, youtubeVideoID sql.NullString) (bool, error) {
	row := q.db.QueryRowContext(ctx, audioExistsByYoutubeVideoID, youtubeVideoID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const countAudioByID = `-- name: CountAudioByID :one
SELECT COUNT(*) FROM audios WHERE id = $1
`

func (q *Queries) CountAudioByID(ctx context.Context, id uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAudioByID, id)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAudio = `-- name: CreateAudio :one
INSERT INTO audios(
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
  local_id,
  playlist_audio_count,
  user_audio_count
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING id, created_at, title, author, duration_ms, path, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id, playlist_audio_count, user_audio_count
`

type CreateAudioParams struct {
	ID                 uuid.UUID
	CreatedAt          time.Time
	Title              sql.NullString
	Author             sql.NullString
	DurationMs         sql.NullInt32
	Path               sql.NullString
	SizeBytes          sql.NullInt64
	YoutubeVideoID     sql.NullString
	ThumbnailPath      sql.NullString
	SpotifyID          sql.NullString
	ThumbnailUrl       sql.NullString
	LocalID            sql.NullString
	PlaylistAudioCount int32
	UserAudioCount     int32
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
		arg.PlaylistAudioCount,
		arg.UserAudioCount,
	)
	var i Audio
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
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
		&i.PlaylistAudioCount,
		&i.UserAudioCount,
	)
	return i, err
}

const decrementPlaylistAudioCountByID = `-- name: DecrementPlaylistAudioCountByID :exec
UPDATE audios SET playlist_audio_count = playlist_audio_count - 1 WHERE id = $1
`

func (q *Queries) DecrementPlaylistAudioCountByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, decrementPlaylistAudioCountByID, id)
	return err
}

const decrementUserAudioCountByID = `-- name: DecrementUserAudioCountByID :exec
UPDATE audios SET user_audio_count = user_audio_count - 1 WHERE id = $1
`

func (q *Queries) DecrementUserAudioCountByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, decrementUserAudioCountByID, id)
	return err
}

const deleteAudioByID = `-- name: DeleteAudioByID :exec
DELETE FROM audios WHERE id = $1
`

func (q *Queries) DeleteAudioByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteAudioByID, id)
	return err
}

const getAllAudioIDs = `-- name: GetAllAudioIDs :many
SELECT id FROM audios
`

func (q *Queries) GetAllAudioIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getAllAudioIDs)
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

const getAudioById = `-- name: GetAudioById :one
SELECT id, created_at, title, author, duration_ms, path, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id, playlist_audio_count, user_audio_count FROM audios WHERE id = $1
`

func (q *Queries) GetAudioById(ctx context.Context, id uuid.UUID) (Audio, error) {
	row := q.db.QueryRowContext(ctx, getAudioById, id)
	var i Audio
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
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
		&i.PlaylistAudioCount,
		&i.UserAudioCount,
	)
	return i, err
}

const getAudioByYoutubeVideoID = `-- name: GetAudioByYoutubeVideoID :one
SELECT id, created_at, title, author, duration_ms, path, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id, playlist_audio_count, user_audio_count FROM audios WHERE youtube_video_id = $1::text
`

func (q *Queries) GetAudioByYoutubeVideoID(ctx context.Context, youtubeVideoID string) (Audio, error) {
	row := q.db.QueryRowContext(ctx, getAudioByYoutubeVideoID, youtubeVideoID)
	var i Audio
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
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
		&i.PlaylistAudioCount,
		&i.UserAudioCount,
	)
	return i, err
}

const getAudioIDsBySpotifyIDs = `-- name: GetAudioIDsBySpotifyIDs :many
SELECT 
  id
FROM audios
WHERE spotify_id = ANY($1::text[])
`

func (q *Queries) GetAudioIDsBySpotifyIDs(ctx context.Context, spotifyIds []string) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getAudioIDsBySpotifyIDs, pq.Array(spotifyIds))
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

const getAudioSpotifyIDsBySpotifyIDs = `-- name: GetAudioSpotifyIDsBySpotifyIDs :many
SELECT 
  id, spotify_id 
FROM audios
WHERE spotify_id = ANY($1::text[])
`

type GetAudioSpotifyIDsBySpotifyIDsRow struct {
	ID        uuid.UUID
	SpotifyID sql.NullString
}

func (q *Queries) GetAudioSpotifyIDsBySpotifyIDs(ctx context.Context, spotifyIds []string) ([]GetAudioSpotifyIDsBySpotifyIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAudioSpotifyIDsBySpotifyIDs, pq.Array(spotifyIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAudioSpotifyIDsBySpotifyIDsRow
	for rows.Next() {
		var i GetAudioSpotifyIDsBySpotifyIDsRow
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

const getAudiosByUserID = `-- name: GetAudiosByUserID :many
SELECT 
  audios.id, audios.created_at, audios.title, audios.author, audios.duration_ms, audios.path, audios.size_bytes, audios.youtube_video_id, audios.thumbnail_path, audios.spotify_id, audios.thumbnail_url, audios.local_id, audios.playlist_audio_count, audios.user_audio_count
  FROM user_audios
  INNER JOIN audios ON user_audios.audio_id = audios.id
  WHERE user_id = $1
`

func (q *Queries) GetAudiosByUserID(ctx context.Context, userID uuid.UUID) ([]Audio, error) {
	rows, err := q.db.QueryContext(ctx, getAudiosByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Audio
	for rows.Next() {
		var i Audio
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
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
			&i.PlaylistAudioCount,
			&i.UserAudioCount,
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

const getUnusedAudios = `-- name: GetUnusedAudios :many
SELECT id, created_at, title, author, duration_ms, path, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id, playlist_audio_count, user_audio_count FROM audios WHERE playlist_audio_count = 0 AND user_audio_count = 0
`

func (q *Queries) GetUnusedAudios(ctx context.Context) ([]Audio, error) {
	rows, err := q.db.QueryContext(ctx, getUnusedAudios)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Audio
	for rows.Next() {
		var i Audio
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
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
			&i.PlaylistAudioCount,
			&i.UserAudioCount,
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

const incrementPlaylistAudioCountByID = `-- name: IncrementPlaylistAudioCountByID :exec
UPDATE audios SET playlist_audio_count = playlist_audio_count + 1 WHERE id = $1
`

func (q *Queries) IncrementPlaylistAudioCountByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, incrementPlaylistAudioCountByID, id)
	return err
}

const incrementUserAudioCountByID = `-- name: IncrementUserAudioCountByID :exec
UPDATE audios SET user_audio_count = user_audio_count + 1 WHERE id = $1
`

func (q *Queries) IncrementUserAudioCountByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, incrementUserAudioCountByID, id)
	return err
}

const updateAudioByID = `-- name: UpdateAudioByID :one
UPDATE audios SET 
  title = COALESCE($1, title),
  author = COALESCE($2, author),
  duration_ms = COALESCE($3, duration_ms),
  path = COALESCE($4, path),
  size_bytes = COALESCE($5, size_bytes),
  youtube_video_id = COALESCE($6, youtube_video_id),
  thumbnail_path = COALESCE($7, thumbnail_path),
  spotify_id = COALESCE($8, spotify_id),
  thumbnail_url = COALESCE($9, thumbnail_url),
  local_id = COALESCE($10, local_id),
  playlist_audio_count = COALESCE($11, playlist_audio_count),
  user_audio_count = COALESCE($12, user_audio_count)
WHERE id = $13
RETURNING id, created_at, title, author, duration_ms, path, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id, playlist_audio_count, user_audio_count
`

type UpdateAudioByIDParams struct {
	Title              sql.NullString
	Author             sql.NullString
	DurationMs         sql.NullInt32
	Path               sql.NullString
	SizeBytes          sql.NullInt64
	YoutubeVideoID     sql.NullString
	ThumbnailPath      sql.NullString
	SpotifyID          sql.NullString
	ThumbnailUrl       sql.NullString
	LocalID            sql.NullString
	PlaylistAudioCount sql.NullInt32
	UserAudioCount     sql.NullInt32
	AudioID            uuid.NullUUID
}

func (q *Queries) UpdateAudioByID(ctx context.Context, arg UpdateAudioByIDParams) (Audio, error) {
	row := q.db.QueryRowContext(ctx, updateAudioByID,
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
		arg.PlaylistAudioCount,
		arg.UserAudioCount,
		arg.AudioID,
	)
	var i Audio
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
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
		&i.PlaylistAudioCount,
		&i.UserAudioCount,
	)
	return i, err
}
