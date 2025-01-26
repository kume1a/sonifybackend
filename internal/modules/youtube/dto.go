package youtube

import "github.com/google/uuid"

type downloadYoutubeAudioToUserLibraryDTO struct {
	VideoID string `json:"videoId" valid:"required"`
}

type downloadYoutubeAudioToPlaylistDTO struct {
	VideoID    string    `json:"videoId" valid:"required"`
	PlaylistID uuid.UUID `json:"playlistId" valid:"required"`
}

type youtubeSearchSuggestions struct {
	Query       string   `json:"query"`
	Suggestions []string `json:"suggestions"`
}

type youtubeVideoInfoDTO struct {
	Title           string `json:"title"`
	Uploader        string `json:"uploader"`
	DurationSeconds int    `json:"durationSeconds"`
}
