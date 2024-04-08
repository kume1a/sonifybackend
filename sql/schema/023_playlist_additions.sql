-- +goose Up

ALTER TABLE user_playlists ADD COLUMN is_spotify_saved_playlist BOOLEAN DEFAULT FALSE;

-- +goose Down

ALTER TABLE user_playlists DROP COLUMN is_spotify_saved_playlist;

