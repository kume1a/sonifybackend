package audio

type uploadUserLocalMusicDTO struct {
	LocalID       string
	Title         string
	Author        string
	AudioPath     string
	ThumbnailPath string
	DurationMs    int32
}

type deleteUnusedAudioResultDTO struct {
	DeletedCount int      `json:"deletedCount"`
	AudioNames   []string `json:"audioNames"`
}
