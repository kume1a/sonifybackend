-- name: CreatePlaylist :one 
INSERT INTO playlists(
  id,
  created_at,
  name,
  thumbnail_path,
  spotify_id,
  thumbnail_url
) VALUES ($1,$2,$3,$4,$5,$6) 
RETURNING *;

-- name: GetPlaylists :many
SELECT * FROM playlists 
  WHERE created_at > $1
  ORDER BY created_at DESC
  LIMIT $2;

-- name: UpdatePlaylistByID :one
UPDATE playlists
SET name = COALESCE($1, name),
    thumbnail_path = COALESCE($2, thumbnail_path)
WHERE id = $3
RETURNING *;

-- name: DeletePlaylistByID :exec
DELETE FROM playlists WHERE id = $1;

-- name: GetPlaylistByID :one
SELECT * FROM playlists WHERE id = $1;

-- name: GetSpotifyUserSavedPlaylistIDs :many
SELECT playlists.id FROM playlists
  INNER JOIN user_playlists ON playlists.id = user_playlists.playlist_id
  WHERE user_playlists.user_id = $1 
  AND user_playlists.is_spotify_saved_playlist = true;

-- name: DeletePlaylistsByIDs :exec
DELETE FROM playlists WHERE id = ANY(sqlc.arg(ids)::uuid[]);