package spotify

import "github.com/google/uuid"

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

		var playlistID *uuid.UUID
		if playlist.DbPlaylist != nil {
			playlistID = &playlist.DbPlaylist.ID
		}

		res[index] = searchSpotifyResPlaylistDTO{
			Name:       playlist.SpotifySearchPlaylist.Name,
			ImageUrl:   imageUrl,
			SpotifyID:  playlist.SpotifySearchPlaylist.ID,
			PlaylistID: playlistID,
		}
	}

	return searchSpotifyResDTO{Playlists: res}
}
