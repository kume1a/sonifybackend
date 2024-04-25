-- name: CreateUserPlaylist :one
INSERT INTO user_playlists(
  user_id,
  playlist_id,
  is_spotify_saved_playlist,
  created_at
) VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetUserPlaylists :many
SELECT 
  playlists.* 
FROM user_playlists
INNER JOIN playlists ON user_playlists.playlist_id = playlists.id
WHERE user_playlists.user_id = $1;

-- name: DeleteSpotifyUserSavedPlaylistJoins :exec
DELETE FROM user_playlists
  WHERE user_playlists.user_id = $1
  AND user_playlists.is_spotify_saved_playlist = true;
