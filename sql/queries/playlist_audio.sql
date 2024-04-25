-- name: CreatePlaylistAudio :one
INSERT INTO playlist_audios(
  playlist_id,
  audio_id,
  created_at
) VALUES ($1,$2,$3) RETURNING *;

-- name: GetPlaylistAudios :many
SELECT 
  audio.*,
  audio_likes.audio_id AS audio_likes_audio_id,
  audio_likes.user_id AS audio_likes_user_id
FROM playlist_audios 
INNER JOIN audio ON playlist_audios.audio_id = audio.id
LEFT JOIN audio_likes ON 
  playlist_audios.audio_id = audio_likes.audio_id 
  AND audio_likes.user_id = $1 
WHERE playlist_audios.playlist_id = $2;

-- name: GetPlaylistAudioJoins :many
SELECT * FROM playlist_audios
  INNER JOIN audio ON playlist_audios.audio_id = audio.id
WHERE (playlist_id = $1 or $1 IS NULL) 
  AND playlist_audios.created_at > $2
ORDER BY playlist_audios.created_at DESC
  LIMIT $3;

-- name: DeletePlaylistAudiosByIds :exec
DELETE FROM playlist_audios 
  WHERE playlist_id = sqlc.arg(playlist_id)
  AND audio_id = ANY(sqlc.arg(audio_ids)::uuid[]);

-- name: GetPlaylistAudioJoinsBySpotifyIds :many
SELECT 
  playlist_audios.*,
  audio.spotify_id AS spotify_id
FROM playlist_audios
INNER JOIN audio ON playlist_audios.audio_id = audio.id
WHERE playlist_audios.playlist_id = sqlc.arg(playlist_id) AND audio.spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);
