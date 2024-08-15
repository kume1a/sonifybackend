-- name: CreateUserAudio :one
INSERT INTO user_audios(
  id,
  created_at,
  user_id, 
  audio_id
) VALUES ($1,$2,$3,$4) 
RETURNING *;

-- name: GetUserAudioByVideoID :one
SELECT * FROM user_audios
  INNER JOIN audios ON user_audios.audio_id = audios.id
  WHERE user_audios.user_id = sqlc.arg(user_id) AND audios.youtube_video_id = sqlc.arg(youtube_video_id);

-- name: CountUserAudioByLocalID :one
SELECT COUNT(*) FROM user_audios 
  INNER JOIN audios ON user_audios.audio_id = audios.id 
  WHERE user_audios.user_id = $1 AND audios.local_id = $2;

-- name: CountUserAudio :one
SELECT COUNT(1) 
FROM user_audios 
WHERE user_id = $1
  AND audio_id = $2;

-- name: GetUserAudioIDs :many
SELECT audio_id FROM user_audios WHERE user_id = $1;

-- name: GetUserAudiosByAudioIds :many
SELECT user_audios.*,
  audio_likes.id AS audio_likes_id,
  audio_likes.created_at AS audio_likes_created_at,
  audio_likes.user_id AS audio_likes_user_id,
  audio_likes.audio_id AS audio_likes_audio_id,

  audios.id AS audio_id,
  audios.created_at AS audio_created_at,
  audios.title AS audio_title,
  audios.author AS audio_author,
  audios.duration_ms AS audio_duration_ms,
  audios.path AS audio_path,
  audios.size_bytes AS audio_size_bytes,
  audios.youtube_video_id AS audio_youtube_video_id,
  audios.thumbnail_path AS audio_thumbnail_path,
  audios.spotify_id AS audio_spotify_id,
  audios.thumbnail_url AS audio_thumbnail_url,
  audios.local_id AS audio_local_id
FROM user_audios
INNER JOIN audios ON user_audios.audio_id = audios.id
LEFT JOIN audio_likes ON audio_likes.audio_id = audios.id
WHERE user_audios.user_id = sqlc.arg(user_id) AND audios.id = ANY(sqlc.arg(audio_ids)::uuid[]);

-- name: DeleteUserAudio :exec
DELETE FROM user_audios WHERE user_id = $1 AND audio_id = $2;