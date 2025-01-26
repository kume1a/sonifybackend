package youtube

type downloadYoutubeAudioDTO struct {
	VideoID string `json:"videoId" valid:"required"`
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
