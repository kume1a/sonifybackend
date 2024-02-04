-- +goose Up

CREATE TABLE audio(
  id UUID PRIMARY KEY,
  title VARCHAR(255),
  artist VARCHAR(255),
  duration INTEGER,
  path VARCHAR(1023),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down

DROP TABLE audio;
