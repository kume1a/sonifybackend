package playlistaudio

import "github.com/kume1a/sonifybackend/internal/database"

func playlistAudioEntityToDto(e *database.PlaylistAudio) playlistAudioDTO {
	return playlistAudioDTO{
		CreatedAt:  e.CreatedAt,
		PlaylistID: e.PlaylistID,
		AudioID:    e.AudioID,
	}
}
