-- +goose Up

ALTER TABLE audio 
  ALTER COLUMN size_bytes TYPE BIGINT;

-- +goose Down

ALTER TABLE audio 
  ALTER COLUMN size_bytes TYPE DECIMAL(10,2);
