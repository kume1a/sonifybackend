-- +goose Up

ALTER TABLE audio
  ADD COLUMN thumbnail_url VARCHAR(255),
  DROP COLUMN updated_at;

ALTER TABLE users
  DROP COLUMN updated_at;

CREATE INDEX idx_audio_spotify_id ON audio (spotify_id);
CREATE INDEX idx_playlist_spotify_id ON playlists (spotify_id);

-- +goose Down
  
ALTER TABLE audio
  DROP COLUMN thumbnail_url,
  ADD COLUMN updated_at TIMESTAMPTZ NOT NULL;

ALTER TABLE users
  ADD COLUMN updated_at TIMESTAMPTZ NOT NULL;

DROP INDEX idx_audio_spotify_id;
DROP INDEX idx_playlist_spotify_id;
