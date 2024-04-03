-- +goose Up

ALTER TABLE playlists 
  ADD COLUMN spotify_id VARCHAR(255);
  
ALTER TABLE audio
  ADD COLUMN spotify_id VARCHAR(255),
  ADD COLUMN remote_url VARCHAR(255);

ALTER TABLE artists
  ADD COLUMN spotify_id VARCHAR(255);

-- +goose Down

ALTER TABLE playlists
  DROP COLUMN spotify_id;

ALTER TABLE audio
  DROP COLUMN spotify_id,
  DROP COLUMN remote_url;

ALTER TABLE artists
  DROP COLUMN spotify_id;
