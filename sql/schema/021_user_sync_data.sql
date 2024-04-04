-- +goose Up

CREATE TABLE user_sync_data (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  spotify_last_synced_at TIMESTAMPTZ NOT NULL
); 

-- +goose Down
  
DROP TABLE user_sync_data;
  