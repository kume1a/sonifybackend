-- name: CreateUserPlaylist :one
INSERT INTO user_playlists(
  id,
  user_id,
  playlist_id,
  is_spotify_saved_playlist,
  created_at
) VALUES ($1,$2,$3,$4,$5) 
RETURNING *;

-- name: GetUserPlaylists :many
SELECT 
  playlists.* 
FROM user_playlists
INNER JOIN playlists ON user_playlists.playlist_id = playlists.id
WHERE user_playlists.user_id = sqlc.arg(user_id) 
  AND (sqlc.arg(ids)::uuid[] IS NULL OR playlists.id = ANY(sqlc.arg(ids)::uuid[]));

-- name: DeleteSpotifyUserSavedPlaylistJoins :exec
DELETE FROM user_playlists
  WHERE user_playlists.user_id = $1
  AND user_playlists.is_spotify_saved_playlist = true;

-- name: GetUserPlaylistIDs :many
SELECT playlist_id FROM user_playlists WHERE user_id = $1;