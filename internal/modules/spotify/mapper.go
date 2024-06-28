package spotify

func MapMergedSpotifySearchToSearchSpotifyResult(
	dto *spotifySearchDTO,
	spotifyResMergedWithDb []spotifySearchPlaylistAndDbPlaylist,
) searchSpotifyResDTO {
	res := make([]searchSpotifyResPlaylistDTO, len(dto.Playlists.Items))

	for index, playlist := range spotifyResMergedWithDb {
		imageUrl := ""
		if len(playlist.SpotifySearchPlaylist.Images) > 0 {
			imageUrl = playlist.SpotifySearchPlaylist.Images[0].URL
		}

		res[index] = searchSpotifyResPlaylistDTO{
			Name:       playlist.SpotifySearchPlaylist.Name,
			ImageUrl:   imageUrl,
			SpotifyID:  playlist.SpotifySearchPlaylist.ID,
			PlaylistID: playlist.DbPlaylist.ID,
		}
	}

	return searchSpotifyResDTO{Playlists: res}
}
