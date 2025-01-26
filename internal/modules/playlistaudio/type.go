package playlistaudio

import "github.com/kume1a/sonifybackend/internal/database"

type PlaylistAudioWithAudio struct {
	PlaylistAudio *database.PlaylistAudio
	Audio         *database.Audio
}
