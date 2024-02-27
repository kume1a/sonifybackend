-- +goose Up

ALTER TABLE audio 
  RENAME COLUMN size_in_mb TO size_bytes;

-- +goose Down

ALTER TABLE audio 
  RENAME COLUMN size_bytes TO size_in_mb;
