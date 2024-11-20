-- +goose Up
ALTER TABLE audios
  ADD COLUMN playlist_audio_count INT DEFAULT 0 NOT NULL,
  ADD COLUMN user_audio_count INT DEFAULT 0 NOT NULL;
  
ALTER TABLE hidden_user_audios
  DROP CONSTRAINT IF EXISTS hidden_user_audios_user_id_fkey,
  ADD CONSTRAINT hidden_user_audios_user_id_fkey
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE hidden_user_audios
  DROP CONSTRAINT IF EXISTS hidden_user_audios_audio_id_fkey,
  ADD CONSTRAINT hidden_user_audios_audio_id_fkey
  FOREIGN KEY (audio_id) REFERENCES audios(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE hidden_user_audios
  DROP CONSTRAINT IF EXISTS hidden_user_audios_audio_id_fkey;

ALTER TABLE hidden_user_audios
  DROP CONSTRAINT IF EXISTS hidden_user_audios_user_id_fkey;

ALTER TABLE audios
  DROP COLUMN IF EXISTS playlist_audio_count,
  DROP COLUMN IF EXISTS user_audio_count;