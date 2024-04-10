-- +goose Up


ALTER TABLE user_playlists
  ALTER COLUMN is_spotify_saved_playlist TYPE BOOL using coalesce(is_spotify_saved_playlist, true),
  ALTER COLUMN is_spotify_saved_playlist SET NOT NULL;

-- +goose Down

ALTER TABLE user_playlists ALTER COLUMN is_spotify_saved_playlist DROP NOT NULL;
