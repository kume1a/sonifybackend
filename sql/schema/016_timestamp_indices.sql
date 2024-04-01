-- +goose Up

CREATE INDEX idx_playlists_created_at ON playlists (created_at);
CREATE INDEX idx_artists_created_at ON artists (created_at);
CREATE INDEX idx_user_audios_created_at ON user_audios (created_at);
CREATE INDEX idx_playlist_audios_created_at ON playlist_audios (created_at);
CREATE INDEX idx_artist_audios_created_at ON artist_audios (created_at);
CREATE INDEX idx_user_playlists_created_at ON user_playlists (created_at);

-- +goose Down

DROP INDEX idx_playlists_created_at;
DROP INDEX idx_artists_created_at;
DROP INDEX idx_user_audios_created_at;
DROP INDEX idx_playlist_audios_created_at;
DROP INDEX idx_artist_audios_created_at;
DROP INDEX idx_user_playlists_created_at;

