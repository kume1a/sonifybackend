package userplaylist

import (
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
)

func MapUserPlaylistFullEntityToDTO(row database.GetFullUserPlaylistsRow) userPlaylistDTO {
	return userPlaylistDTO{
		ID:                     row.UserPlaylistID,
		PlaylistID:             row.UserPlaylistPlaylistID,
		UserID:                 row.UserPlaylistUserID,
		CreatedAt:              row.UserPlaylistCreatedAt,
		IsSpotifySavedPlaylist: row.UserPlaylistIsSpotifySavedPlaylist,
		Playlist: &playlist.PlaylistDTO{
			ID:                row.PlaylistID,
			CreatedAt:         row.PlaylistCreatedAt,
			Name:              row.PlaylistName,
			ThumbnailPath:     row.PlaylistThumbnailPath.String,
			ThumbnailUrl:      row.PlaylistThumbnailUrl.String,
			SpotifyID:         row.PlaylistSpotifyID.String,
			AudioImportStatus: row.PlaylistAudioImportStatus,
			AudioCount:        row.PlaylistAudioCount,
			TotalAudioCount:   row.PlaylistTotalAudioCount,
		},
	}
}

func MapUserPlaylistEntityToDTO(e database.UserPlaylist) userPlaylistDTO {
	return userPlaylistDTO{
		ID:                     e.ID,
		PlaylistID:             e.PlaylistID,
		UserID:                 e.UserID,
		CreatedAt:              e.CreatedAt,
		IsSpotifySavedPlaylist: e.IsSpotifySavedPlaylist,
		Playlist:               nil,
	}
}
