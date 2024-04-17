-- name: CreateAudio :one 
INSERT INTO audio(
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
  local_id
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING *;

-- name: GetAudiosByUserId :many
SELECT 
  audio.*
  FROM user_audios
  INNER JOIN audio ON user_audios.audio_id = audio.id
  WHERE user_id = $1;

-- name: DeleteAudioById :exec
DELETE FROM audio WHERE id = $1;

-- name: GetAudioById :one
SELECT * FROM audio WHERE id = $1;

-- name: UpdateAudio :one
UPDATE audio SET 
  title = $1, 
  author = $2, 
  duration_ms = $3, 
  path = $4, 
  thumbnail_path=$5 
WHERE id = $6 
RETURNING *;

-- name: GetUserAudioByVideoId :one
SELECT * FROM user_audios
  INNER JOIN audio ON user_audios.audio_id = audio.id
  WHERE user_audios.user_id = $1 AND audio.youtube_video_id = $2;

-- name: CreateUserAudio :one
INSERT INTO user_audios(
  created_at,
  user_id, 
  audio_id
) VALUES ($1,$2,$3) RETURNING *;

-- name: GetPlaylistAudiosBySpotifyIds :many
SELECT 
  audio.*
  FROM playlist_audios
  INNER JOIN audio ON playlist_audios.audio_id = audio.id
  WHERE playlist_audios.playlist_id = sqlc.arg(playlist_id) AND audio.spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: GetAudioSpotifyIdsBySpotifyIds :many
SELECT 
  id, spotify_id 
FROM audio 
WHERE spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: GetAudioIdsBySpotifyIds :many
SELECT 
  id
FROM audio
WHERE spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: CountUserAudioByLocalId :one
SELECT COUNT(*) FROM user_audios 
  INNER JOIN audio ON user_audios.audio_id = audio.id 
  WHERE user_audios.user_id = $1 AND audio.local_id = $2;

-- name: GetUserAudioIds :many
SELECT audio_id FROM user_audios WHERE user_id = $1;

-- name: GetUserAudiosByAudioIds :many
SELECT user_audios.*,
  audio.id as audio_id,
  audio.created_at as audio_created_at,
  audio.title as audio_title,
  audio.author as audio_author,
  audio.duration_ms as audio_duration_ms,
  audio.path as audio_path,
  audio.size_bytes as audio_size_bytes,
  audio.youtube_video_id as audio_youtube_video_id,
  audio.thumbnail_path as audio_thumbnail_path,
  audio.spotify_id as audio_spotify_id,
  audio.thumbnail_url as audio_thumbnail_url,
  audio.local_id as audio_local_id
FROM user_audios
INNER JOIN audio ON user_audios.audio_id = audio.id
WHERE user_audios.user_id = sqlc.arg(user_id) AND audio.id = ANY(sqlc.arg(audio_ids)::uuid[]);
