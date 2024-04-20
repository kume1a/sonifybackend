-- +goose Up

CREATE TABLE audio_likes (
  user_id UUID NOT NULL,
  audio_id UUID NOT NULL,
  PRIMARY KEY (user_id, audio_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (audio_id) REFERENCES audio(id)
);

-- +goose Down

DROP TABLE audio_likes;