-- +goose Up

ALTER TABLE users ALTER COLUMN name DROP NOT NULL;

-- +goose Down

ALTER TABLE users ALTER COLUMN name SET NOT NULL;
