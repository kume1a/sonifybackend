-- name: CreateAudioLike :one
INSERT INTO audio_likes(
  id,
  audio_id, 
  user_id,
  created_at
) VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: DeleteAudioLike :exec
DELETE FROM audio_likes 
  WHERE audio_id = $1 AND user_id = $2;

-- name: GetAudioLikesByUserID :many
SELECT * 
FROM audio_likes 
WHERE user_id = sqlc.arg(user_id)::uuid;

-- name: GetAudioLikesByUserIDAndAudioIDs :many
SELECT *
FROM audio_likes
WHERE user_id = sqlc.arg(user_id)::uuid AND 
  audio_id = ANY(sqlc.arg(audio_ids)::uuid[]);