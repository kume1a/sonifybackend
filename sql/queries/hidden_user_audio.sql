-- name: CreateHiddenUserAudio :one
INSERT INTO hidden_user_audios(
  id,
  audio_id, 
  user_id,
  created_at
) VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: DeleteHiddenUserAudio :exec
DELETE FROM hidden_user_audios
  WHERE audio_id = $1 AND user_id = $2;

-- name: GetHiddenUserAudiosByUserID :many
SELECT * 
FROM hidden_user_audios
WHERE user_id = sqlc.arg(user_id)::uuid;

-- name: GetHiddenUserAudiosByUserIDAndAudioIDs :many
SELECT *
FROM hidden_user_audios
WHERE user_id = sqlc.arg(user_id)::uuid AND 
  audio_id = ANY(sqlc.arg(audio_ids)::uuid[]);