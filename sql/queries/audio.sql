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
  local_id,
  playlist_audio_count,
  user_audio_count
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING *;

-- name: GetAudiosByUserID :many
SELECT 
  audios.*
  FROM user_audios
  INNER JOIN audios ON user_audios.audio_id = audios.id
  WHERE user_id = $1;

-- name: GetAllAudioIDs :many
SELECT id FROM audios;

-- name: GetAudioById :one
SELECT * FROM audios WHERE id = $1;

-- name: UpdateAudioByID :one
UPDATE audios SET 
  title = COALESCE(sqlc.narg(title), title),
  author = COALESCE(sqlc.narg(author), author),
  duration_ms = COALESCE(sqlc.narg(duration_ms), duration_ms),
  path = COALESCE(sqlc.narg(path), path),
  size_bytes = COALESCE(sqlc.narg(size_bytes), size_bytes),
  youtube_video_id = COALESCE(sqlc.narg(youtube_video_id), youtube_video_id),
  thumbnail_path = COALESCE(sqlc.narg(thumbnail_path), thumbnail_path),
  spotify_id = COALESCE(sqlc.narg(spotify_id), spotify_id),
  thumbnail_url = COALESCE(sqlc.narg(thumbnail_url), thumbnail_url),
  local_id = COALESCE(sqlc.narg(local_id), local_id),
  playlist_audio_count = COALESCE(sqlc.narg(playlist_audio_count), playlist_audio_count),
  user_audio_count = COALESCE(sqlc.narg(user_audio_count), user_audio_count)
WHERE id = sqlc.narg(audio_id)
RETURNING *;

-- name: IncrementUserAudioCountByID :exec
UPDATE audios SET user_audio_count = user_audio_count + 1 WHERE id = $1;

-- name: DecrementUserAudioCountByID :exec
UPDATE audios SET user_audio_count = user_audio_count - 1 WHERE id = $1;

-- name: IncrementPlaylistAudioCountByID :exec
UPDATE audios SET playlist_audio_count = playlist_audio_count + 1 WHERE id = $1;

-- name: DecrementPlaylistAudioCountByID :exec
UPDATE audios SET playlist_audio_count = playlist_audio_count - 1 WHERE id = $1;

-- name: GetAudioSpotifyIDsBySpotifyIDs :many
SELECT 
  id, spotify_id 
FROM audios
WHERE spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: GetUnusedAudios :many
SELECT * FROM audios WHERE playlist_audio_count = 0 AND user_audio_count = 0;

-- name: DeleteAudioByID :exec
DELETE FROM audios WHERE id = $1;

-- name: GetAudioIDsBySpotifyIDs :many
SELECT 
  id
FROM audios
WHERE spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: CountAudioByID :one
SELECT COUNT(*) FROM audios WHERE id = $1;