-- +goose Up

ALTER TABLE playlists 
  ADD COLUMN thumbnail_url VARCHAR(255);
  
ALTER TABLE audio
  DROP COLUMN remote_url;

ALTER TABLE artists
  ADD COLUMN image_url VARCHAR(255);

-- +goose Down
ALTER TABLE playlists 
  DROP COLUMN thumbnail_url;
  
ALTER TABLE audio
  ADD COLUMN remote_url VARCHAR(255);

ALTER TABLE artists
  DROP COLUMN image_url;
