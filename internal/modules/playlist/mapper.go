package playlist

import "github.com/kume1a/sonifybackend/internal/database"

func PlaylistEntityToDto(e database.Playlist) PlaylistDTO {
	return PlaylistDTO{
		ID:                e.ID,
		CreatedAt:         e.CreatedAt,
		Name:              e.Name,
		ThumbnailPath:     e.ThumbnailPath.String,
		ThumbnailUrl:      e.ThumbnailUrl.String,
		SpotifyID:         e.SpotifyID.String,
		AudioImportStatus: e.AudioImportStatus,
		AudioCount:        e.AudioCount,
		TotalAudioCount:   e.TotalAudioCount,
	}
}

func playlistAudioEntityToDto(e *database.PlaylistAudio) playlistAudioDTO {
	return playlistAudioDTO{
		CreatedAt:  e.CreatedAt,
		PlaylistID: e.PlaylistID,
		AudioID:    e.AudioID,
	}
}
