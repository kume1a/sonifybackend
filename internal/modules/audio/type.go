package audio

import "github.com/kume1a/sonifybackend/internal/database"

type UserAudioWithAudio struct {
	UserAudio *database.UserAudio
	Audio     *database.Audio
}
