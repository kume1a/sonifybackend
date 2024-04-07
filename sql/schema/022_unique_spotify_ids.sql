-- +goose Up

DROP INDEX idx_audio_spotify_id;
DROP INDEX idx_playlist_spotify_id;

CREATE UNIQUE INDEX idx_audio_spotify_id ON audio (spotify_id);
CREATE UNIQUE INDEX idx_playlist_spotify_id ON playlists (spotify_id);

-- +goose Down

DROP INDEX idx_audio_spotify_id ON audio;
DROP INDEX idx_playlist_spotify_id ON playlists;

CREATE INDEX idx_audio_spotify_id ON audio (spotify_id);
CREATE INDEX idx_playlist_spotify_id ON playlists (spotify_id);

