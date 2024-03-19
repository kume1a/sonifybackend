-- +goose Up

ALTER TABLE audio 
  ADD COLUMN thumbnail_path VARCHAR(1023);

-- +goose Down

ALTER TABLE audio 
  DROP COLUMN thumbnail_path;
