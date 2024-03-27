package youtube

type YoutubeVideoInfo struct {
	Title             string `json:"title"`
	Uploader          string `json:"uploader"`
	DurationInSeconds int    `json:"duration"`
}
