// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: playlist.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createPlaylist = `-- name: CreatePlaylist :one
INSERT INTO playlists(
  id,
  created_at,
  name,
  thumbnail_path,
  spotify_id,
  thumbnail_url
) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, name, thumbnail_path, created_at, spotify_id, thumbnail_url
`

type CreatePlaylistParams struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	Name          string
	ThumbnailPath sql.NullString
	SpotifyID     sql.NullString
	ThumbnailUrl  sql.NullString
}

func (q *Queries) CreatePlaylist(ctx context.Context, arg CreatePlaylistParams) (Playlist, error) {
	row := q.db.QueryRowContext(ctx, createPlaylist,
		arg.ID,
		arg.CreatedAt,
		arg.Name,
		arg.ThumbnailPath,
		arg.SpotifyID,
		arg.ThumbnailUrl,
	)
	var i Playlist
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ThumbnailPath,
		&i.CreatedAt,
		&i.SpotifyID,
		&i.ThumbnailUrl,
	)
	return i, err
}

const createPlaylistAudio = `-- name: CreatePlaylistAudio :one
INSERT INTO playlist_audios(
  playlist_id,
  audio_id,
  created_at
) VALUES ($1,$2,$3) RETURNING playlist_id, audio_id, created_at
`

type CreatePlaylistAudioParams struct {
	PlaylistID uuid.UUID
	AudioID    uuid.UUID
	CreatedAt  time.Time
}

func (q *Queries) CreatePlaylistAudio(ctx context.Context, arg CreatePlaylistAudioParams) (PlaylistAudio, error) {
	row := q.db.QueryRowContext(ctx, createPlaylistAudio, arg.PlaylistID, arg.AudioID, arg.CreatedAt)
	var i PlaylistAudio
	err := row.Scan(&i.PlaylistID, &i.AudioID, &i.CreatedAt)
	return i, err
}

const createUserPlaylist = `-- name: CreateUserPlaylist :one
INSERT INTO user_playlists(
  user_id,
  playlist_id,
  is_spotify_saved_playlist,
  created_at
) VALUES ($1,$2,$3,$4) RETURNING user_id, playlist_id, created_at, is_spotify_saved_playlist
`

type CreateUserPlaylistParams struct {
	UserID                 uuid.UUID
	PlaylistID             uuid.UUID
	IsSpotifySavedPlaylist bool
	CreatedAt              time.Time
}

func (q *Queries) CreateUserPlaylist(ctx context.Context, arg CreateUserPlaylistParams) (UserPlaylist, error) {
	row := q.db.QueryRowContext(ctx, createUserPlaylist,
		arg.UserID,
		arg.PlaylistID,
		arg.IsSpotifySavedPlaylist,
		arg.CreatedAt,
	)
	var i UserPlaylist
	err := row.Scan(
		&i.UserID,
		&i.PlaylistID,
		&i.CreatedAt,
		&i.IsSpotifySavedPlaylist,
	)
	return i, err
}

const deletePlaylistAudiosByIds = `-- name: DeletePlaylistAudiosByIds :exec
DELETE FROM playlist_audios 
  WHERE playlist_id = $1
  AND audio_id = ANY($2::uuid[])
`

type DeletePlaylistAudiosByIdsParams struct {
	PlaylistID uuid.UUID
	AudioIds   []uuid.UUID
}

func (q *Queries) DeletePlaylistAudiosByIds(ctx context.Context, arg DeletePlaylistAudiosByIdsParams) error {
	_, err := q.db.ExecContext(ctx, deletePlaylistAudiosByIds, arg.PlaylistID, pq.Array(arg.AudioIds))
	return err
}

const deletePlaylistById = `-- name: DeletePlaylistById :exec
DELETE FROM playlists WHERE id = $1
`

func (q *Queries) DeletePlaylistById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePlaylistById, id)
	return err
}

const deletePlaylistsByIds = `-- name: DeletePlaylistsByIds :exec
DELETE FROM playlists WHERE id = ANY($1::uuid[])
`

func (q *Queries) DeletePlaylistsByIds(ctx context.Context, ids []uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePlaylistsByIds, pq.Array(ids))
	return err
}

const deleteSpotifyUserSavedPlaylistJoins = `-- name: DeleteSpotifyUserSavedPlaylistJoins :exec
DELETE FROM user_playlists
  WHERE user_playlists.user_id = $1
  AND user_playlists.is_spotify_saved_playlist = true
`

func (q *Queries) DeleteSpotifyUserSavedPlaylistJoins(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteSpotifyUserSavedPlaylistJoins, userID)
	return err
}

const getPlaylistAudioJoins = `-- name: GetPlaylistAudioJoins :many
SELECT playlist_id, audio_id, playlist_audios.created_at, id, title, author, duration_ms, path, audio.created_at, size_bytes, youtube_video_id, thumbnail_path, spotify_id, thumbnail_url FROM playlist_audios
  INNER JOIN audio ON playlist_audios.audio_id = audio.id
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
	PlaylistID     uuid.UUID
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
			&i.PlaylistID,
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

const getPlaylistAudioJoinsBySpotifyIds = `-- name: GetPlaylistAudioJoinsBySpotifyIds :many
SELECT 
  playlist_audios.playlist_id, playlist_audios.audio_id, playlist_audios.created_at,
  audio.spotify_id AS spotify_id
FROM playlist_audios
INNER JOIN audio ON playlist_audios.audio_id = audio.id
WHERE playlist_audios.playlist_id = $1 AND audio.spotify_id = ANY($2::text[])
`

type GetPlaylistAudioJoinsBySpotifyIdsParams struct {
	PlaylistID uuid.UUID
	SpotifyIds []string
}

type GetPlaylistAudioJoinsBySpotifyIdsRow struct {
	PlaylistID uuid.UUID
	AudioID    uuid.UUID
	CreatedAt  time.Time
	SpotifyID  sql.NullString
}

func (q *Queries) GetPlaylistAudioJoinsBySpotifyIds(ctx context.Context, arg GetPlaylistAudioJoinsBySpotifyIdsParams) ([]GetPlaylistAudioJoinsBySpotifyIdsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylistAudioJoinsBySpotifyIds, arg.PlaylistID, pq.Array(arg.SpotifyIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlaylistAudioJoinsBySpotifyIdsRow
	for rows.Next() {
		var i GetPlaylistAudioJoinsBySpotifyIdsRow
		if err := rows.Scan(
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
SELECT audio.id, audio.title, audio.author, audio.duration_ms, audio.path, audio.created_at, audio.size_bytes, audio.youtube_video_id, audio.thumbnail_path, audio.spotify_id, audio.thumbnail_url 
  FROM playlist_audios 
  INNER JOIN audio ON playlist_audios.audio_id = audio.id
  WHERE playlist_audios.playlist_id = $1
`

func (q *Queries) GetPlaylistAudios(ctx context.Context, playlistID uuid.UUID) ([]Audio, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylistAudios, playlistID)
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

const getPlaylistById = `-- name: GetPlaylistById :one
SELECT id, name, thumbnail_path, created_at, spotify_id, thumbnail_url FROM playlists WHERE id = $1
`

func (q *Queries) GetPlaylistById(ctx context.Context, id uuid.UUID) (Playlist, error) {
	row := q.db.QueryRowContext(ctx, getPlaylistById, id)
	var i Playlist
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ThumbnailPath,
		&i.CreatedAt,
		&i.SpotifyID,
		&i.ThumbnailUrl,
	)
	return i, err
}

const getPlaylists = `-- name: GetPlaylists :many
SELECT id, name, thumbnail_path, created_at, spotify_id, thumbnail_url FROM playlists 
  WHERE created_at > $1
  ORDER BY created_at DESC
  LIMIT $2
`

type GetPlaylistsParams struct {
	CreatedAt time.Time
	Limit     int32
}

func (q *Queries) GetPlaylists(ctx context.Context, arg GetPlaylistsParams) ([]Playlist, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylists, arg.CreatedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Playlist
	for rows.Next() {
		var i Playlist
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ThumbnailPath,
			&i.CreatedAt,
			&i.SpotifyID,
			&i.ThumbnailUrl,
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

const getSpotifyUserSavedPlaylistIds = `-- name: GetSpotifyUserSavedPlaylistIds :many
SELECT id FROM playlists
  INNER JOIN user_playlists ON playlists.id = user_playlists.playlist_id
  WHERE user_playlists.user_id = $1 
  AND user_playlists.is_spotify_saved_playlist = true
`

func (q *Queries) GetSpotifyUserSavedPlaylistIds(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getSpotifyUserSavedPlaylistIds, userID)
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

const getUserPlaylists = `-- name: GetUserPlaylists :many
SELECT 
  playlists.id, playlists.name, playlists.thumbnail_path, playlists.created_at, playlists.spotify_id, playlists.thumbnail_url 
FROM user_playlists
INNER JOIN playlists ON user_playlists.playlist_id = playlists.id
WHERE user_playlists.user_id = $1
`

func (q *Queries) GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]Playlist, error) {
	rows, err := q.db.QueryContext(ctx, getUserPlaylists, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Playlist
	for rows.Next() {
		var i Playlist
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ThumbnailPath,
			&i.CreatedAt,
			&i.SpotifyID,
			&i.ThumbnailUrl,
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

const updatePlaylistById = `-- name: UpdatePlaylistById :one
UPDATE playlists
  SET name = COALESCE($1, name),
      thumbnail_path = COALESCE($2, thumbnail_path)
  WHERE id = $3
  RETURNING id, name, thumbnail_path, created_at, spotify_id, thumbnail_url
`

type UpdatePlaylistByIdParams struct {
	Name          string
	ThumbnailPath sql.NullString
	ID            uuid.UUID
}

func (q *Queries) UpdatePlaylistById(ctx context.Context, arg UpdatePlaylistByIdParams) (Playlist, error) {
	row := q.db.QueryRowContext(ctx, updatePlaylistById, arg.Name, arg.ThumbnailPath, arg.ID)
	var i Playlist
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ThumbnailPath,
		&i.CreatedAt,
		&i.SpotifyID,
		&i.ThumbnailUrl,
	)
	return i, err
}
