package spotify

type spotifyTrack struct {
	ID              string
	Name            string
	Artist          string
	DurationSeconds int
	ThumbnailUrl    string
}

type spotifyPlaylist struct {
	ID     string
	Name   string
	Tracks []spotifyTrack
}
