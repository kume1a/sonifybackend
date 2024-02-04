-- +goose Up

ALTER TABLE audio
ADD COLUMN user_id UUID,
ADD CONSTRAINT fk_audio_user
FOREIGN KEY (user_id) REFERENCES users(id);

-- +goose Down

ALTER TABLE audio
DROP CONSTRAINT IF EXISTS fk_audio_user,
DROP COLUMN IF EXISTS user_id;
