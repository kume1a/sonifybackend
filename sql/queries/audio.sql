-- name: CreateAudio :one 
INSERT INTO audio(
  id, 
  created_at,
  title,
  author,
  duration,
  path,
  size_bytes,
  youtube_video_id,
  thumbnail_path,
  spotify_id,
  thumbnail_url
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING *;

-- name: GetAudiosByUserId :many
SELECT 
  audio.*,
  user_audios.user_id AS user_id
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
  duration = $3, 
  path = $4, 
  thumbnail_path=$5 
WHERE id = $6 
RETURNING *;

-- name: GetUserAudioByVideoId :one
SELECT * FROM user_audios
  INNER JOIN audio ON user_audios.audio_id = audio.id
  WHERE user_audios.user_id = $1 AND audio.youtube_video_id = $2;

-- name: CreateUserAudio :one
INSERT INTO user_audios(user_id, audio_id) VALUES ($1, $2) RETURNING *;