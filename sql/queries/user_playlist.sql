-- name: CreateUserPlaylist :one
INSERT INTO user_playlists(
  id,
  user_id,
  playlist_id,
  is_spotify_saved_playlist,
  created_at
) VALUES ($1,$2,$3,$4,$5) 
RETURNING *;

-- name: GetFullUserPlaylists :many
SELECT 
  user_playlists.id as user_playlist_id,
  user_playlists.user_id as user_playlist_user_id,
  user_playlists.playlist_id as user_playlist_playlist_id,
  user_playlists.is_spotify_saved_playlist as user_playlist_is_spotify_saved_playlist,
  user_playlists.created_at as user_playlist_created_at,
  playlists.id as playlist_id,
  playlists.created_at as playlist_created_at,
  playlists.name as playlist_name,
  playlists.thumbnail_path as playlist_thumbnail_path,
  playlists.thumbnail_url as playlist_thumbnail_url,
  playlists.spotify_id as playlist_spotify_id,
  playlists.audio_import_status as playlist_audio_import_status,
  playlists.audio_count as playlist_audio_count,
  playlists.total_audio_count as playlist_total_audio_count
FROM user_playlists
INNER JOIN playlists ON user_playlists.playlist_id = playlists.id
WHERE user_playlists.user_id = sqlc.arg(user_id) 
  AND (sqlc.arg(playlist_ids)::uuid[] IS NULL OR playlists.id = ANY(sqlc.arg(playlist_ids)::uuid[]))
ORDER BY user_playlists.created_at DESC;

-- name: GetUserPlaylists :many
SELECT * FROM user_playlists
WHERE user_id = sqlc.arg(user_id) 
  AND (sqlc.arg(ids)::uuid[] IS NULL OR id = ANY(sqlc.arg(ids)::uuid[]));

-- name: DeleteSpotifyUserSavedPlaylistJoins :exec
DELETE FROM user_playlists
  WHERE user_playlists.user_id = $1
  AND user_playlists.is_spotify_saved_playlist = true;

-- name: GetPlaylistIDsByUserID :many
SELECT playlist_id FROM user_playlists WHERE user_id = $1;

-- name: GetUserPlaylistIDsByUserID :many
SELECT id FROM user_playlists WHERE user_id = $1;