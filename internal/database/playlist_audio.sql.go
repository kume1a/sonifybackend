// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: playlist_audio.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createPlaylistAudio = `-- name: CreatePlaylistAudio :one
INSERT INTO playlist_audios(
  playlist_id,
  audio_id,
  created_at
) VALUES ($1,$2,$3) 
RETURNING id, playlist_id, audio_id, created_at
`

type CreatePlaylistAudioParams struct {
	PlaylistID uuid.UUID
	AudioID    uuid.UUID
	CreatedAt  time.Time
}

func (q *Queries) CreatePlaylistAudio(ctx context.Context, arg CreatePlaylistAudioParams) (PlaylistAudio, error) {
	row := q.db.QueryRowContext(ctx, createPlaylistAudio, arg.PlaylistID, arg.AudioID, arg.CreatedAt)
	var i PlaylistAudio
	err := row.Scan(
		&i.ID,
		&i.PlaylistID,
		&i.AudioID,
		&i.CreatedAt,
	)
	return i, err
}

const deletePlaylistAudiosByIDs = `-- name: DeletePlaylistAudiosByIDs :exec
DELETE FROM playlist_audios 
  WHERE playlist_id = $1
  AND audio_id = ANY($2::uuid[])
`

type DeletePlaylistAudiosByIDsParams struct {
	PlaylistID uuid.UUID
	AudioIds   []uuid.UUID
}

func (q *Queries) DeletePlaylistAudiosByIDs(ctx context.Context, arg DeletePlaylistAudiosByIDsParams) error {
	_, err := q.db.ExecContext(ctx, deletePlaylistAudiosByIDs, arg.PlaylistID, pq.Array(arg.AudioIds))
	return err
}

const getPlaylistAudioJoins = `-- name: GetPlaylistAudioJoins :many
SELECT playlist_audios.id, playlist_id, audio_id, playlist_audios.created_at, audios.id, title, author, duration_ms, path, audios.created_at, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url, local_id FROM playlist_audios
  INNER JOIN audios ON playlist_audios.audio_id = audios.id
WHERE (playlist_id = $1 or $1 IS NULL) 
  AND playlist_audios.created_at > $2
ORDER BY playlist_audios.created_at DESC
  LIMIT $3
`

type GetPlaylistAudioJoinsParams struct {
	PlaylistID uuid.UUID
	CreatedAt  time.Time
	Limit      int32
}

type GetPlaylistAudioJoinsRow struct {
	ID             uuid.UUID
	PlaylistID     uuid.UUID
	AudioID        uuid.UUID
	CreatedAt      time.Time
	ID_2           uuid.UUID
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

func (q *Queries) GetPlaylistAudioJoins(ctx context.Context, arg GetPlaylistAudioJoinsParams) ([]GetPlaylistAudioJoinsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylistAudioJoins, arg.PlaylistID, arg.CreatedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlaylistAudioJoinsRow
	for rows.Next() {
		var i GetPlaylistAudioJoinsRow
		if err := rows.Scan(
			&i.ID,
			&i.PlaylistID,
			&i.AudioID,
			&i.CreatedAt,
			&i.ID_2,
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

const getPlaylistAudioJoinsBySpotifyIDs = `-- name: GetPlaylistAudioJoinsBySpotifyIDs :many
SELECT 
  playlist_audios.id, playlist_audios.playlist_id, playlist_audios.audio_id, playlist_audios.created_at,
  audios.spotify_id AS spotify_id
FROM playlist_audios
INNER JOIN audios ON playlist_audios.audio_id = audios.id
WHERE playlist_audios.playlist_id = $1 AND audios.spotify_id = ANY($2::text[])
`

type GetPlaylistAudioJoinsBySpotifyIDsParams struct {
	PlaylistID uuid.UUID
	SpotifyIds []string
}

type GetPlaylistAudioJoinsBySpotifyIDsRow struct {
	ID         uuid.UUID
	PlaylistID uuid.UUID
	AudioID    uuid.UUID
	CreatedAt  time.Time
	SpotifyID  sql.NullString
}

func (q *Queries) GetPlaylistAudioJoinsBySpotifyIDs(ctx context.Context, arg GetPlaylistAudioJoinsBySpotifyIDsParams) ([]GetPlaylistAudioJoinsBySpotifyIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylistAudioJoinsBySpotifyIDs, arg.PlaylistID, pq.Array(arg.SpotifyIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlaylistAudioJoinsBySpotifyIDsRow
	for rows.Next() {
		var i GetPlaylistAudioJoinsBySpotifyIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.PlaylistID,
			&i.AudioID,
			&i.CreatedAt,
			&i.SpotifyID,
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

const getPlaylistAudios = `-- name: GetPlaylistAudios :many
SELECT 
  audios.id, audios.title, audios.author, audios.duration_ms, audios.path, audios.created_at, audios.size_bytes, audios.youtube_video_id, audios.thumbnail_path, audios.spotify_id, audios.thumbnail_url, audios.local_id,
  audio_likes.audio_id AS audio_likes_audio_id,
  audio_likes.user_id AS audio_likes_user_id
FROM playlist_audios 
INNER JOIN audios ON playlist_audios.audio_id = audios.id
LEFT JOIN audio_likes ON 
  playlist_audios.audio_id = audio_likes.audio_id 
  AND audio_likes.user_id = $1 
WHERE playlist_audios.playlist_id = $2
`

type GetPlaylistAudiosParams struct {
	UserID     uuid.UUID
	PlaylistID uuid.UUID
}

type GetPlaylistAudiosRow struct {
	ID                uuid.UUID
	Title             sql.NullString
	Author            sql.NullString
	DurationMs        sql.NullInt32
	Path              sql.NullString
	CreatedAt         time.Time
	SizeBytes         sql.NullInt64
	YoutubeVideoID    sql.NullString
	ThumbnailPath     sql.NullString
	SpotifyID         sql.NullString
	ThumbnailUrl      sql.NullString
	LocalID           sql.NullString
	AudioLikesAudioID uuid.NullUUID
	AudioLikesUserID  uuid.NullUUID
}

func (q *Queries) GetPlaylistAudios(ctx context.Context, arg GetPlaylistAudiosParams) ([]GetPlaylistAudiosRow, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylistAudios, arg.UserID, arg.PlaylistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlaylistAudiosRow
	for rows.Next() {
		var i GetPlaylistAudiosRow
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
			&i.AudioLikesAudioID,
			&i.AudioLikesUserID,
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
