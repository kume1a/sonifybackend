-- name: CreatePlaylist :one 
INSERT INTO playlists(
  id,
  created_at,
  name,
  thumbnail_path,
  spotify_id,
  thumbnail_url,
  audio_import_status,
  audio_count,
  total_audio_count
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) 
RETURNING *;

-- name: UpdatePlaylistByID :one
UPDATE playlists
SET 
  name = COALESCE(sqlc.narg(name), name),
  thumbnail_path = COALESCE(sqlc.narg(thumbnail_path), thumbnail_path),
  spotify_id = COALESCE(sqlc.narg(spotify_id), spotify_id),
  thumbnail_url = COALESCE(sqlc.narg(thumbnail_url), thumbnail_url),
  audio_import_status = COALESCE(sqlc.narg(audio_import_status), audio_import_status),
  audio_count = COALESCE(sqlc.narg(audio_count), audio_count),
  total_audio_count = COALESCE(sqlc.narg(total_audio_count), total_audio_count)
WHERE id = sqlc.arg(playlist_id)
RETURNING *;

-- name: GetPlaylistsBySpotifyIDs :many
SELECT * FROM playlists WHERE spotify_id = ANY(sqlc.arg(spotify_ids)::text[]);

-- name: GetPlaylistBySpotifyID :one
SELECT * FROM playlists WHERE spotify_id = sqlc.arg(spotify_id)::text;

-- name: GetPlaylistByID :one
SELECT * FROM playlists WHERE id = $1;

-- name: GetPlaylistIDBySpotifyID :one
SELECT id FROM playlists WHERE spotify_id = sqlc.arg(spotify_id)::text;

-- name: GetSpotifyUserSavedPlaylistIDs :many
SELECT playlists.id FROM playlists
  INNER JOIN user_playlists ON playlists.id = user_playlists.playlist_id
  WHERE user_playlists.user_id = $1 
  AND user_playlists.is_spotify_saved_playlist = true;

-- name: DeletePlaylistsByIDs :exec
DELETE FROM playlists WHERE id = ANY(sqlc.arg(ids)::uuid[]);

-- name: DeletePlaylistByID :exec
DELETE FROM playlists WHERE id = $1;
