-- name: CreateAudio :one 
INSERT INTO audio(
  id, 
  created_at,
  updated_at,
  title,
  author,
  duration,
  path,
  user_id
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING *;

-- name: GetAudiosByUserId :many
SELECT * FROM audio WHERE user_id = $1;

-- name: DeleteAudioById :exec
DELETE FROM audio WHERE id = $1;

-- name: GetAudioById :one
SELECT * FROM audio WHERE id = $1;

-- name: UpdateAudio :one
UPDATE audio SET title = $1, author = $2, duration = $3, path = $4 WHERE id = $5 RETURNING *;
