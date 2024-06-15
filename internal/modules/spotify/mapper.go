package spotify

func MapSpotifySearchToSearchSpotifyResult(dto *spotifySearchDTO) *searchSpotifyResDTO {
	res := make([]searchSpotifyResPlaylistDTO, len(dto.Playlists.Items))

	for _, playlist := range dto.Playlists.Items {
		res = append(res, searchSpotifyResPlaylistDTO{
			Title:     playlist.Name,
			ImageUrl:  playlist.Images[0].URL,
			SpotifyID: playlist.ID,
		})
	}

	return &searchSpotifyResDTO{Playlists: res}
}
