package playlistaudio

import "github.com/kume1a/sonifybackend/internal/database"

func playlistAudioEntityToDTO(e *database.PlaylistAudio) playlistAudioDTO {
	return playlistAudioDTO{
		CreatedAt:  e.CreatedAt,
		PlaylistID: e.PlaylistID,
		AudioID:    e.AudioID,
	}
}

func playlistAudioEntityListToDTO(e []database.PlaylistAudio) []playlistAudioDTO {
	dto := make([]playlistAudioDTO, 0, len(e))
	for _, v := range e {
		dto = append(dto, playlistAudioEntityToDTO(&v))
	}
	return dto
}
