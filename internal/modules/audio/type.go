package audio

import "github.com/kume1a/sonifybackend/internal/database"

type UserAudioWithAudio struct {
	UserAudio *database.UserAudio
	Audio     *database.Audio
}

type AudioWithAudioLike struct {
	Audio     *database.Audio
	AudioLike *database.AudioLike
}
