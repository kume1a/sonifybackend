package playlistaudio

import (
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
)

func playlistAudioEntityToDTO(e *database.PlaylistAudio) playlistAudioDTO {
	return playlistAudioDTO{
		ID:         e.ID,
		CreatedAt:  e.CreatedAt,
		PlaylistID: e.PlaylistID,
		AudioID:    e.AudioID,
	}
}

func getPlaylistAudioRowToDTO(e *database.GetPlaylistAudiosRow) playlistAudioDTO {
	return playlistAudioDTO{
		ID:         e.PlaylistAudioID,
		CreatedAt:  e.PlaylistAudioCreatedAt,
		PlaylistID: e.PlaylistAudioPlaylistID,
		AudioID:    e.PlaylistAudioAudioID,
		Audio: &audio.AudioDTO{
			ID:             e.AudioID.UUID,
			Title:          e.AudioTitle.String,
			Author:         e.AudioAuthor.String,
			CreatedAt:      e.AudioCreatedAt.Time,
			DurationMs:     e.AudioDurationMs.Int32,
			Path:           e.AudioPath.String,
			SizeBytes:      e.AudioSizeBytes.Int64,
			YoutubeVideoID: e.AudioYoutubeVideoID.String,
			ThumbnailPath:  e.AudioThumbnailPath.String,
			SpotifyID:      e.AudioSpotifyID.String,
			ThumbnailUrl:   e.AudioThumbnailUrl.String,
			LocalID:        e.AudioLocalID.String,
		},
	}
}

func getPlaylistAudioRowListToDTO(e []database.GetPlaylistAudiosRow) []playlistAudioDTO {
	dto := make([]playlistAudioDTO, 0, len(e))
	for _, v := range e {
		dto = append(dto, getPlaylistAudioRowToDTO(&v))
	}
	return dto
}

func playlistAudioEntityListToDTO(e []database.PlaylistAudio) []playlistAudioDTO {
	dto := make([]playlistAudioDTO, 0, len(e))
	for _, v := range e {
		dto = append(dto, playlistAudioEntityToDTO(&v))
	}
	return dto
}
