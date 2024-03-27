-- +goose Up

CREATE TYPE auth_provider AS ENUM ('EMAIL', 'GOOGLE', 'FACEBOOK', 'APPLE');

ALTER TABLE users 
  ADD COLUMN auth_provider auth_provider NOT NULL DEFAULT 'GOOGLE',
  ADD COLUMN password_hash TEXT;

-- +goose Down

ALTER TABLE users 
  DROP COLUMN auth_provider, password_hash;

DROP TYPE auth_provider;
