-- name: CreateAudio :one 
INSERT INTO audios(
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

-- name: GetAudiosByUserID :many
SELECT 
  audios.*
  FROM user_audios
  INNER JOIN audios ON user_audios.audio_id = audios.id
  WHERE user_id = $1;

-- name: DeleteAudioById :exec
DELETE FROM audios WHERE id = $1;

-- name: GetAudioById :one
SELECT * FROM audios WHERE id = $1;

-- name: UpdateAudio :one
UPDATE audios SET 
  title = $1, 
  author = $2, 
  duration_ms = $3, 
  path = $4, 
  thumbnail_path=$5 
WHERE id = $6 
RETURNING *;

-- name: GetAudioSpotifyIDsBySpotifyIDs :many
SELECT 
  id, spotify_id 
FROM audios
WHERE spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: GetAudioIDsBySpotifyIDs :many
SELECT 
  id
FROM audios
WHERE spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);
