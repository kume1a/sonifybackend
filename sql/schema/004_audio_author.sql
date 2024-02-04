-- +goose Up

ALTER TABLE audio RENAME COLUMN artist TO author;

-- +goose Down

ALTER TABLE audio RENAME COLUMN author TO artist;