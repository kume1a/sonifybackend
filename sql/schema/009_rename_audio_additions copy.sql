-- +goose Up

ALTER TABLE audio 
  RENAME COLUMN sizeInMb TO size_in_mb;

ALTER TABLE audio 
  RENAME COLUMN youtubeVideoId TO youtube_video_id;

-- +goose Down

ALTER TABLE audio 
  RENAME COLUMN size_in_mb TO sizeInMb;

ALTER TABLE audio 
  RENAME COLUMN youtube_video_id TO youtubeVideoId;
