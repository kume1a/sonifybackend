package userplaylist

import (
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
)

func MapUserPlaylistFullEntityToDTO(e database.GetFullUserPlaylistsRow) userPlaylistDTO {
	return userPlaylistDTO{
		ID:                     e.UserPlaylistID,
		PlaylistID:             e.UserPlaylistPlaylistID,
		UserID:                 e.UserPlaylistUserID,
		CreatedAt:              e.UserPlaylistCreatedAt,
		IsSpotifySavedPlaylist: e.UserPlaylistIsSpotifySavedPlaylist,
		Playlist: &sharedmodule.PlaylistDTO{
			ID:                e.PlaylistID,
			CreatedAt:         e.PlaylistCreatedAt,
			Name:              e.PlaylistName,
			ThumbnailPath:     e.PlaylistThumbnailPath.String,
			ThumbnailUrl:      e.PlaylistThumbnailUrl.String,
			SpotifyID:         e.PlaylistSpotifyID.String,
			AudioImportStatus: e.PlaylistAudioImportStatus,
			AudioCount:        e.PlaylistAudioCount,
			TotalAudioCount:   e.PlaylistTotalAudioCount,
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
