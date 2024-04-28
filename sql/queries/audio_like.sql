-- name: CreateAudioLike :one
INSERT INTO audio_likes (
  audio_id, 
  user_id
) VALUES ($1, $2) 
RETURNING *;

-- name: DeleteAudioLike :exec
DELETE FROM audio_likes 
  WHERE audio_id = $1 AND user_id = $2;

-- name: GetAudioLikesByUserId :many
SELECT * FROM audio_likes WHERE user_id = $1;

-- name: DeleteUserAudioLikesByAudioIDs :exec
DELETE FROM audio_likes WHERE user_id = sqlc.arg(user_id) AND audio_id = ANY(sqlc.arg(audio_ids)::uuid[]);

-- name: GetAudioLikesByUserIDAndAudioIDs :many
SELECT * FROM audio_likes 
  WHERE user_id = sqlc.arg(user_id) AND audio_id = ANY(sqlc.arg(audio_ids)::uuid[]);