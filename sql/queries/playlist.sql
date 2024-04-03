-- name: CreatePlaylist :one 
INSERT INTO playlists(
  id,
  created_at,
  name,
  thumbnail_path,
  spotify_id,
  thumbnail_url
) VALUES ($1,$2,$3,$4,$5,$6) RETURNING *;

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

-- name: CreateUserPlaylist :one
INSERT INTO user_playlists(
  user_id,
  playlist_id
) VALUES ($1,$2) RETURNING *;

-- name: GetUserPlaylistsBySpotifyIds :many
SELECT 
  playlists.* 
  FROM user_playlists
  INNER JOIN playlists ON user_playlists.playlist_id = playlists.id
  WHERE user_playlists.user_id = sqlc.arg(user_id) AND playlists.spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: DeletePlaylistAudiosByIds :exec
DELETE FROM playlist_audios 
  WHERE playlist_id = sqlc.arg(playlist_id)
  AND audio_id = ANY(sqlc.arg(audio_ids)::uuid[]);