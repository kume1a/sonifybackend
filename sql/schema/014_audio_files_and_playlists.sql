-- +goose Up

ALTER TABLE users 
  ALTER COLUMN password_hash TYPE VARCHAR(255);  

ALTER TABLE audio
  DROP CONSTRAINT IF EXISTS fk_audio_user,
  DROP COLUMN IF EXISTS user_id;

CREATE TABLE playlists (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  thumbnail_path VARCHAR(255)
);

CREATE TABLE artists (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  image_path VARCHAR(255) NOT NULL
);

CREATE TABLE user_audios (
  user_id UUID,
  audio_id UUID,
  PRIMARY KEY (user_id, audio_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (audio_id) REFERENCES audio(id)
);

CREATE TABLE playlist_audios (
  playlist_id UUID,
  audio_id UUID,
  PRIMARY KEY (playlist_id, audio_id),
  FOREIGN KEY (playlist_id) REFERENCES playlists(id),
  FOREIGN KEY (audio_id) REFERENCES audio(id)
);

CREATE TABLE artist_audios (
  artist_id UUID,
  audio_id UUID,
  PRIMARY KEY (artist_id, audio_id),
  FOREIGN KEY (artist_id) REFERENCES artists(id),
  FOREIGN KEY (audio_id) REFERENCES audio(id)
);

CREATE TABLE user_playlists (
  user_id UUID,
  playlist_id UUID,
  PRIMARY KEY (user_id, playlist_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (playlist_id) REFERENCES playlists(id)
);

-- +goose Down
ALTER TABLE users 
  ALTER COLUMN password_hash TYPE TEXT;

ALTER TABLE audio
  ADD COLUMN user_id UUID,
  ADD CONSTRAINT fk_audio_user
  FOREIGN KEY (user_id) REFERENCES users(id);

DROP TABLE IF EXISTS playlists;

DROP TABLE IF EXISTS artists;

DROP TABLE IF EXISTS user_audios;

DROP TABLE IF EXISTS playlist_audios;

DROP TABLE IF EXISTS artist_audios;

DROP TABLE IF EXISTS user_playlists;
