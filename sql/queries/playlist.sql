-- name: CreatePlaylist :one 
INSERT INTO playlists(
  id,
  name,
  thumbnail_path
) VALUES ($1,$2,$3) RETURNING *;

-- name: GetPlaylists :many
SELECT * FROM playlists LIMIT $1;