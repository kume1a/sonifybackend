package spotify

func MapSpotifySearchToSearchSpotifyResult(dto *spotifySearchDTO) searchSpotifyResDTO {
	res := make([]searchSpotifyResPlaylistDTO, len(dto.Playlists.Items))

	for index, playlist := range dto.Playlists.Items {
		imageUrl := ""
		if len(playlist.Images) > 0 {
			imageUrl = playlist.Images[0].URL
		}

		res[index] = searchSpotifyResPlaylistDTO{
			Name:      playlist.Name,
			ImageUrl:  imageUrl,
			SpotifyID: playlist.ID,
		}
	}

	return searchSpotifyResDTO{Playlists: res}
}
