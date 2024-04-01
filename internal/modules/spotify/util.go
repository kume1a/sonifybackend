package spotify

import "github.com/asaskevich/govalidator"

func (dto *downloadSpotifyPlaylist) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func spotifyPlaylistDtoToModel(dto *spotifyPlaylistDTO) *spotifyPlaylist {
	model := &spotifyPlaylist{
		ID:     dto.ID,
		Name:   dto.Name,
		Tracks: []spotifyTrack{},
	}

	for _, track := range dto.Tracks.Items {
		artist := ""
		if len(track.Track.Artists) > 0 {
			artist = track.Track.Artists[0].Name
		}

		thumbnailUrl := ""
		if len(track.Track.Album.Images) > 0 {
			thumbnailUrl = track.Track.Album.Images[0].URL
		}

		model.Tracks = append(model.Tracks, spotifyTrack{
			ID:              track.Track.ID,
			Name:            track.Track.Name,
			Artist:          artist,
			DurationSeconds: track.Track.DurationMS / 1000,
			ThumbnailUrl:    thumbnailUrl,
		})
	}

	return model
}
