-- +goose Up

ALTER TABLE audio 
  ADD COLUMN sizeInMb DECIMAL(10,2),
  ADD COLUMN youtubeVideoId VARCHAR(255);

-- +goose Down

ALTER TABLE audio 
  DROP COLUMN sizeInMb,
  DROP COLUMN youtubeVideoId;
