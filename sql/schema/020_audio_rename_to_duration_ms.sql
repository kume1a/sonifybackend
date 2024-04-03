-- +goose Up

ALTER TABLE audio
  RENAME COLUMN duration TO duration_ms;

-- +goose Down
  
ALTER TABLE audio
  RENAME COLUMN duration_ms TO duration;
  