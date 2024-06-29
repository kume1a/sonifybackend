// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user_playlist.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createUserPlaylist = `-- name: CreateUserPlaylist :one
INSERT INTO user_playlists(
  id,
  user_id,
  playlist_id,
  is_spotify_saved_playlist,
  created_at
) VALUES ($1,$2,$3,$4,$5) 
RETURNING id, created_at, user_id, playlist_id, is_spotify_saved_playlist
`

type CreateUserPlaylistParams struct {
	ID                     uuid.UUID
	UserID                 uuid.UUID
	PlaylistID             uuid.UUID
	IsSpotifySavedPlaylist bool
	CreatedAt              time.Time
}

func (q *Queries) CreateUserPlaylist(ctx context.Context, arg CreateUserPlaylistParams) (UserPlaylist, error) {
	row := q.db.QueryRowContext(ctx, createUserPlaylist,
		arg.ID,
		arg.UserID,
		arg.PlaylistID,
		arg.IsSpotifySavedPlaylist,
		arg.CreatedAt,
	)
	var i UserPlaylist
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.PlaylistID,
		&i.IsSpotifySavedPlaylist,
	)
	return i, err
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

const getFullUserPlaylists = `-- name: GetFullUserPlaylists :many
SELECT 
  user_playlists.id as user_playlist_id,
  user_playlists.user_id as user_playlist_user_id,
  user_playlists.playlist_id as user_playlist_playlist_id,
  user_playlists.is_spotify_saved_playlist as user_playlist_is_spotify_saved_playlist,
  user_playlists.created_at as user_playlist_created_at,
  playlists.id as playlist_id,
  playlists.created_at as playlist_created_at,
  playlists.name as playlist_name,
  playlists.thumbnail_path as playlist_thumbnail_path,
  playlists.thumbnail_url as playlist_thumbnail_url,
  playlists.spotify_id as playlist_spotify_id,
  playlists.audio_import_status as playlist_audio_import_status,
  playlists.audio_count as playlist_audio_count,
  playlists.total_audio_count as playlist_total_audio_count
FROM user_playlists
INNER JOIN playlists ON user_playlists.playlist_id = playlists.id
WHERE user_playlists.user_id = $1 
  AND ($2::uuid[] IS NULL OR playlists.id = ANY($2::uuid[]))
ORDER BY user_playlists.created_at DESC
`

type GetFullUserPlaylistsParams struct {
	UserID      uuid.UUID
	PlaylistIds []uuid.UUID
}

type GetFullUserPlaylistsRow struct {
	UserPlaylistID                     uuid.UUID
	UserPlaylistUserID                 uuid.UUID
	UserPlaylistPlaylistID             uuid.UUID
	UserPlaylistIsSpotifySavedPlaylist bool
	UserPlaylistCreatedAt              time.Time
	PlaylistID                         uuid.UUID
	PlaylistCreatedAt                  time.Time
	PlaylistName                       string
	PlaylistThumbnailPath              sql.NullString
	PlaylistThumbnailUrl               sql.NullString
	PlaylistSpotifyID                  sql.NullString
	PlaylistAudioImportStatus          ProcessStatus
	PlaylistAudioCount                 int32
	PlaylistTotalAudioCount            int32
}

func (q *Queries) GetFullUserPlaylists(ctx context.Context, arg GetFullUserPlaylistsParams) ([]GetFullUserPlaylistsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFullUserPlaylists, arg.UserID, pq.Array(arg.PlaylistIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFullUserPlaylistsRow
	for rows.Next() {
		var i GetFullUserPlaylistsRow
		if err := rows.Scan(
			&i.UserPlaylistID,
			&i.UserPlaylistUserID,
			&i.UserPlaylistPlaylistID,
			&i.UserPlaylistIsSpotifySavedPlaylist,
			&i.UserPlaylistCreatedAt,
			&i.PlaylistID,
			&i.PlaylistCreatedAt,
			&i.PlaylistName,
			&i.PlaylistThumbnailPath,
			&i.PlaylistThumbnailUrl,
			&i.PlaylistSpotifyID,
			&i.PlaylistAudioImportStatus,
			&i.PlaylistAudioCount,
			&i.PlaylistTotalAudioCount,
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

const getPlaylistIDsByUserID = `-- name: GetPlaylistIDsByUserID :many
SELECT playlist_id FROM user_playlists WHERE user_id = $1
`

func (q *Queries) GetPlaylistIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylistIDsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var playlist_id uuid.UUID
		if err := rows.Scan(&playlist_id); err != nil {
			return nil, err
		}
		items = append(items, playlist_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserPlaylistIDsByUserID = `-- name: GetUserPlaylistIDsByUserID :many
SELECT id FROM user_playlists WHERE user_id = $1
`

func (q *Queries) GetUserPlaylistIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getUserPlaylistIDsByUserID, userID)
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

const getUserPlaylistUserIDsByPlaylistID = `-- name: GetUserPlaylistUserIDsByPlaylistID :many
SELECT user_id FROM user_playlists WHERE playlist_id = $1
`

func (q *Queries) GetUserPlaylistUserIDsByPlaylistID(ctx context.Context, playlistID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getUserPlaylistUserIDsByPlaylistID, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var user_id uuid.UUID
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		items = append(items, user_id)
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
SELECT id, created_at, user_id, playlist_id, is_spotify_saved_playlist FROM user_playlists
WHERE user_id = $1 
  AND ($2::uuid[] IS NULL OR id = ANY($2::uuid[]))
`

type GetUserPlaylistsParams struct {
	UserID uuid.UUID
	Ids    []uuid.UUID
}

func (q *Queries) GetUserPlaylists(ctx context.Context, arg GetUserPlaylistsParams) ([]UserPlaylist, error) {
	rows, err := q.db.QueryContext(ctx, getUserPlaylists, arg.UserID, pq.Array(arg.Ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserPlaylist
	for rows.Next() {
		var i UserPlaylist
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UserID,
			&i.PlaylistID,
			&i.IsSpotifySavedPlaylist,
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
