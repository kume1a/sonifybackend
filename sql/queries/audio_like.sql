-- name: CreateAudioLike :one
INSERT INTO audio_likes (
  audio_id, 
  user_id
) VALUES ($1, $2) 
RETURNING *;

-- name: DeleteAudioLike :exec
DELETE FROM audio_likes 
  WHERE audio_id = $1 AND user_id = $2;