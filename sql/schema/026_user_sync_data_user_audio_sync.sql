-- +goose Up

ALTER TABLE user_sync_data ADD COLUMN user_audio_last_synced_at TIMESTAMPTZ;

-- +goose Down

ALTER TABLE user_sync_data DROP COLUMN user_audio_last_synced_at;