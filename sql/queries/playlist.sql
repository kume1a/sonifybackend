-- name: CreatePlaylist :one 
INSERT INTO playlists(
  id,
  created_at,
  name,
  thumbnail_path
) VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetPlaylists :many
SELECT * FROM playlists 
  WHERE created_at > $1
  ORDER BY created_at DESC
  LIMIT $2;

-- name: UpdatePlaylistById :one
UPDATE playlists
  SET name = COALESCE($1, name),
      thumbnail_path = COALESCE($2, thumbnail_path)
  WHERE id = $3
  RETURNING *;

-- name: DeletePlaylistById :exec
DELETE FROM playlists WHERE id = $1;

-- name: GetPlaylistById :one
SELECT * FROM playlists WHERE id = $1;

-- name: CreatePlaylistAudio :one
INSERT INTO playlist_audios(
  playlist_id,
  audio_id
) VALUES ($1,$2) RETURNING *;

-- name: GetPlaylistAudios :many
SELECT * FROM playlist_audios
  INNER JOIN audio ON playlist_audios.audio_id = audio.id
WHERE (playlist_id = $1 or $1 IS NULL) 
  AND playlist_audios.created_at > $2
ORDER BY playlist_audios.created_at DESC
  LIMIT $3;
