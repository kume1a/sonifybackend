-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

----------------- ENUMS -----------------
create type auth_provider as enum (
  'EMAIL', 
  'GOOGLE', 
  'FACEBOOK', 
  'APPLE'
);

----------------- USERS -----------------
CREATE TABLE users
(
  id            UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  name          VARCHAR(255),
  email         VARCHAR(255),
  auth_provider auth_provider NOT NULL,
  password_hash VARCHAR(255)
);

CREATE UNIQUE INDEX idx_users_email
  ON users (email);

CREATE INDEX idx_users_created_at
  ON users (created_at);

----------------- ARTISTS -----------------
CREATE TABLE artists
(
    id         UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    name       VARCHAR(255) NOT NULL,
    image_path VARCHAR(255) NOT NULL,
    spotify_id VARCHAR(255),
    image_url  VARCHAR(255)
);

CREATE INDEX idx_artists_created_at
    ON artists (created_at);

----------------- AUDIOS -----------------
CREATE TABLE audios
(
    id               UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at       TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    title            VARCHAR(255),
    author           VARCHAR(255),
    duration_ms      INTEGER,
    path             VARCHAR(1023),
    size_bytes       BIGINT,
    youtube_video_id VARCHAR(255),
    thumbnail_path   VARCHAR(1023),
    spotify_id       VARCHAR(255),
    thumbnail_url    VARCHAR(255),
    local_id         VARCHAR(255)
);

CREATE UNIQUE INDEX idx_audios_spotify_id
    ON audios (spotify_id);

CREATE INDEX idx_audios_created_at
    ON audios (created_at);

----------------- PLAYLISTS -----------------
CREATE TYPE process_status AS ENUM (
  'PENDING',
  'PROCESSING',
  'COMPLETED',
  'FAILED'
);

CREATE TABLE playlists
(
  id                  UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  name                VARCHAR(255) NOT NULL,
  thumbnail_path      VARCHAR(255),
  spotify_id          VARCHAR(255),
  thumbnail_url       VARCHAR(255),
  audio_import_status process_status NOT NULL,
  audio_count         INTEGER NOT NULL,
  total_audio_count   INTEGER NOT NULL
);

CREATE INDEX idx_playlists_created_at
  ON playlists (created_at);

CREATE UNIQUE INDEX idx_playlist_spotify_id
  ON playlists (spotify_id);

----------------- ARTIST AUDIOS -----------------
CREATE TABLE artist_audios
(
  id         UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  artist_id  UUID NOT NULL REFERENCES artists ON DELETE CASCADE,
  audio_id   UUID NOT NULL REFERENCES audios ON DELETE CASCADE,
  UNIQUE (artist_id, audio_id)
);

CREATE INDEX idx_artist_audios_created_at
  ON artist_audios (created_at);

----------------- AUDIO LIKES -----------------
CREATE TABLE audio_likes
(
  id       UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  user_id  UUID NOT NULL REFERENCES users ON DELETE CASCADE,
  audio_id UUID NOT NULL REFERENCES audios ON DELETE CASCADE,
  UNIQUE (user_id, audio_id)
);

CREATE INDEX idx_audio_likes_created_at
  ON audio_likes (created_at);

----------------- PLAYLIST AUDIOS -----------------
CREATE TABLE playlist_audios
(
  id          UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  playlist_id UUID NOT NULL REFERENCES playlists ON DELETE CASCADE,
  audio_id    UUID NOT NULL REFERENCES audios ON DELETE CASCADE,
  UNIQUE (playlist_id, audio_id)
);

CREATE INDEX idx_playlist_audios_created_at
  ON playlist_audios (created_at);

----------------- USER AUDIOS -----------------
CREATE TABLE user_audios
(
  id         UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  user_id    UUID NOT NULL REFERENCES users ON DELETE CASCADE,
  audio_id   UUID NOT NULL REFERENCES audios ON DELETE CASCADE,
  UNIQUE (user_id, audio_id)
);

CREATE INDEX idx_user_audios_created_at
  ON user_audios (created_at);

----------------- USER PLAYLISTS -----------------
CREATE TABLE user_playlists
(
  id                        UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at                TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  
  user_id                   UUID NOT NULL REFERENCES users ON DELETE CASCADE,
  playlist_id               UUID NOT NULL REFERENCES playlists ON DELETE CASCADE,
  is_spotify_saved_playlist BOOLEAN NOT NULL,
  UNIQUE (user_id, playlist_id)
);

CREATE INDEX idx_user_playlists_created_at
  ON user_playlists (created_at);

----------------- USER SYNC DATA -----------------
CREATE TABLE user_sync_data
(
  id                        UUID NOT NULL PRIMARY KEY,
  user_id                   UUID NOT NULL REFERENCES users ON DELETE CASCADE,
  spotify_last_synced_at    TIMESTAMPTZ,
  user_audio_last_synced_at TIMESTAMPTZ,
  UNIQUE (user_id)
);

-- +goose Down

DROP TABLE artist_audios;
DROP TABLE artists;
DROP TABLE audios;
DROP TABLE audio_likes;
DROP TABLE playlist_audios;
DROP TABLE playlists;
DROP TABLE user_audios;
DROP TABLE user_playlists;
DROP TABLE user_sync_data;
DROP TABLE users;

DROP INDEX idx_audios_created_at;
DROP INDEX idx_audios_spotify_id;
DROP INDEX idx_artist_audios_created_at;
DROP INDEX idx_artists_created_at;
DROP INDEX idx_audio_likes_created_at;
DROP INDEX idx_playlist_audios_created_at;
DROP INDEX idx_playlists_created_at;
DROP INDEX idx_user_audios_created_at;
DROP INDEX idx_user_playlists_created_at;
DROP INDEX idx_users_created_at;
DROP INDEX idx_users_email;
