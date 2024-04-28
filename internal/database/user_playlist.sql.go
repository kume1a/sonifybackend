// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user_playlist.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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

const deleteSpotifyUserSavedPlaylistJoins = `-- name: DeleteSpotifyUserSavedPlaylistJoins :exec
DELETE FROM user_playlists
  WHERE user_playlists.user_id = $1
  AND user_playlists.is_spotify_saved_playlist = true
`

func (q *Queries) DeleteSpotifyUserSavedPlaylistJoins(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteSpotifyUserSavedPlaylistJoins, userID)
	return err
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
