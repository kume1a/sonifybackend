-- +goose Up

ALTER TABLE audio ADD COLUMN local_id VARCHAR(255);

-- +goose Down

ALTER TABLE audio DROP COLUMN local_id;