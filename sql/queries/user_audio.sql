-- name: CreateUserAudio :one
INSERT INTO user_audios(
  created_at,
  user_id, 
  audio_id
) VALUES ($1,$2,$3) 
RETURNING *;

-- name: GetUserAudioByVideoID :one
SELECT * FROM user_audios
  INNER JOIN audios ON user_audios.audio_id = audios.id
  WHERE user_audios.user_id = sqlc.arg(user_id) AND audios.youtube_video_id = sqlc.arg(youtube_video_id);

-- name: CountUserAudioByLocalID :one
SELECT COUNT(*) FROM user_audios 
  INNER JOIN audios ON user_audios.audio_id = audios.id 
  WHERE user_audios.user_id = $1 AND audios.local_id = $2;

-- name: GetUserAudioIDs :many
SELECT audio_id FROM user_audios WHERE user_id = $1;

-- name: GetUserAudiosByAudioIds :many
SELECT user_audios.*,
  audio_likes.user_id as audio_likes_user_id,
  audio_likes.audio_id as audio_likes_audio_id,
  audios.id as audio_id,
  audios.created_at as audio_created_at,
  audios.title as audio_title,
  audios.author as audio_author,
  audios.duration_ms as audio_duration_ms,
  audios.path as audio_path,
  audios.size_bytes as audio_size_bytes,
  audios.youtube_video_id as audio_youtube_video_id,
  audios.thumbnail_path as audio_thumbnail_path,
  audios.spotify_id as audio_spotify_id,
  audios.thumbnail_url as audio_thumbnail_url,
  audios.local_id as audio_local_id
FROM user_audios
INNER JOIN audios ON user_audios.audio_id = audios.id
LEFT JOIN audio_likes ON audio_likes.audio_id = audios.id
WHERE user_audios.user_id = sqlc.arg(user_id) AND audios.id = ANY(sqlc.arg(audio_ids)::uuid[]);
