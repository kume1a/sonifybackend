-- name: CreateAudio :one 
INSERT INTO audio(
  id, 
  created_at,
  updated_at,
  title,
  author,
  duration,
  path,
  user_id,
  size_bytes,
  youtube_video_id,
  thumbnail_path
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING *;

-- name: GetAudiosByUserId :many
SELECT * FROM audio WHERE user_id = $1;

-- name: DeleteAudioById :exec
DELETE FROM audio WHERE id = $1;

-- name: GetAudioById :one
SELECT * FROM audio WHERE id = $1;

-- name: UpdateAudio :one
UPDATE audio SET title = $1, author = $2, duration = $3, path = $4, thumbnail_path=$5 WHERE id = $6 RETURNING *;

-- name: GetUserAudioByVideoId :one
SELECT * FROM audio WHERE user_id = $1 AND youtube_video_id = $2;