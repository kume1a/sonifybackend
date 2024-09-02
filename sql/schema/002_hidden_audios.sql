-- +goose Up
CREATE TABLE IF NOT EXISTS hidden_user_audios(
  id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  user_id UUID NOT NULL,
  audio_id UUID NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (audio_id) REFERENCES audios(id)
);

CREATE INDEX IF NOT EXISTS idx_hidden_user_audios_created_at
  ON hidden_user_audios (created_at);

CREATE INDEX IF NOT EXISTS idx_hidden_user_audios_user_id
  ON hidden_user_audios (user_id);

-- +goose Down

DROP TABLE IF EXISTS hidden_user_audios;

DROP INDEX IF EXISTS idx_hidden_user_audios_created_at;
DROP INDEX IF EXISTS idx_hidden_user_audios_user_id;
