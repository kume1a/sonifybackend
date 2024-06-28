package spotify

import "github.com/kume1a/sonifybackend/internal/database"

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

type spotifySearchPlaylistAndDbPlaylist struct {
	SpotifySearchPlaylist spotifySearchPlaylistItemDTO
	DbPlaylist            *database.Playlist
}
