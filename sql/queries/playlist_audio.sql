-- name: CreatePlaylistAudio :one
INSERT INTO playlist_audios(
  id,
  playlist_id,
  audio_id,
  created_at
) VALUES ($1,$2,$3,$4) 
RETURNING *;

-- name: GetPlaylistAudiosWithAudioAndAudioLikes :many
SELECT 
  audios.*,
  audio_likes.audio_id AS audio_likes_audio_id,
  audio_likes.user_id AS audio_likes_user_id
FROM playlist_audios 
INNER JOIN audios ON playlist_audios.audio_id = audios.id
LEFT JOIN audio_likes ON 
  playlist_audios.audio_id = audio_likes.audio_id 
  AND audio_likes.user_id = $1 
WHERE playlist_audios.playlist_id = $2;

-- name: GetPlaylistAudios :many
SELECT
  playlist_audios.id AS playlist_audio_id,
  playlist_audios.created_at AS playlist_audio_created_at,
  playlist_audios.playlist_id AS playlist_audio_playlist_id,
  playlist_audios.audio_id AS playlist_audio_audio_id,

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
FROM playlist_audios
LEFT JOIN audios ON playlist_audios.audio_id = audios.id
WHERE (sqlc.arg(playlist_ids)::uuid[] IS NULL OR playlist_audios.playlist_id = ANY(sqlc.arg(playlist_ids)::uuid[])) 
  AND (sqlc.arg(ids)::uuid[] IS NULL OR playlist_audios.id = ANY(sqlc.arg(ids)::uuid[]));

-- name: GetPlaylistAudioIDsByPlaylistIDs :many
SELECT id
FROM playlist_audios
WHERE playlist_id = ANY(sqlc.arg(playlist_ids)::uuid[]);

-- name: DeletePlaylistAudiosByIDs :exec
DELETE FROM playlist_audios 
  WHERE playlist_id = sqlc.arg(playlist_id)
  AND audio_id = ANY(sqlc.arg(audio_ids)::uuid[]);

-- name: GetPlaylistAudioJoinsBySpotifyIDs :many
SELECT 
  playlist_audios.*,
  audios.spotify_id AS spotify_id
FROM playlist_audios
INNER JOIN audios ON playlist_audios.audio_id = audios.id
WHERE playlist_audios.playlist_id = sqlc.arg(playlist_id) AND audios.spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);
