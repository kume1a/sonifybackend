package audio

type uploadUserLocalMusicDTO struct {
	LocalID       string
	Title         string
	Author        string
	AudioPath     string
	ThumbnailPath string
	DurationMs    int32
}
