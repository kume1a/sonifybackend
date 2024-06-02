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
SELECT * 
FROM playlist_audios
WHERE (sqlc.arg(playlist_ids)::uuid[] IS NULL OR playlist_id = ANY(sqlc.arg(playlist_ids)::uuid[])) 
  AND (sqlc.arg(ids)::uuid[] IS NULL OR id = ANY(sqlc.arg(ids)::uuid[]));

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
