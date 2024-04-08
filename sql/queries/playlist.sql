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

-- name: GetPlaylistAudioJoins :many
SELECT * FROM playlist_audios
  INNER JOIN audio ON playlist_audios.audio_id = audio.id
WHERE (playlist_id = $1 or $1 IS NULL) 
  AND playlist_audios.created_at > $2
ORDER BY playlist_audios.created_at DESC
  LIMIT $3;

-- name: CreateUserPlaylist :one
INSERT INTO user_playlists(
  user_id,
  playlist_id,
  is_spotify_saved_playlist
) VALUES ($1,$2,$3) RETURNING *;

-- name: DeletePlaylistAudiosByIds :exec
DELETE FROM playlist_audios 
  WHERE playlist_id = sqlc.arg(playlist_id)
  AND audio_id = ANY(sqlc.arg(audio_ids)::uuid[]);

-- name: GetPlaylistAudioJoinsBySpotifyIds :many
SELECT 
  playlist_audios.*,
  audio.spotify_id AS spotify_id
FROM playlist_audios
INNER JOIN audio ON playlist_audios.audio_id = audio.id
WHERE playlist_audios.playlist_id = sqlc.arg(playlist_id) AND audio.spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: GetUserPlaylists :many
SELECT 
  playlists.* 
FROM user_playlists
INNER JOIN playlists ON user_playlists.playlist_id = playlists.id
WHERE user_playlists.user_id = $1;

-- name: GetPlaylistAudios :many
SELECT audio.* 
  FROM playlist_audios 
  INNER JOIN audio ON playlist_audios.audio_id = audio.id
  WHERE playlist_audios.playlist_id = $1;

-- name: GetSpotifyUserSavedPlaylistIds :many
SELECT id FROM playlists
  INNER JOIN user_playlists ON playlists.id = user_playlists.playlist_id
  WHERE user_playlists.user_id = $1 
  AND user_playlists.is_spotify_saved_playlist = true;

-- name: DeleteSpotifyUserSavedPlaylistJoins :exec
DELETE FROM user_playlists
  WHERE user_playlists.user_id = $1
  AND user_playlists.is_spotify_saved_playlist = true;

-- name: DeletePlaylistsByIds :exec
DELETE FROM playlists WHERE id = ANY(sqlc.arg(ids)::uuid[]);