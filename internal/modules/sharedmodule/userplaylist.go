package sharedmodule

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

type UserPlaylistDTO struct {
	ID                     uuid.UUID    `json:"id"`
	UserID                 uuid.UUID    `json:"userId"`
	PlaylistID             uuid.UUID    `json:"playlistId"`
	CreatedAt              time.Time    `json:"createdAt"`
	IsSpotifySavedPlaylist bool         `json:"isSpotifySavedPlaylist"`
	Playlist               *PlaylistDTO `json:"playlist"`
}

type UserPlaylistWithRelDTO struct {
	*UserPlaylistDTO
	Playlist *PlaylistDTO `json:"playlist"`
}

type UserPlaylistWithRel struct {
	UserPlaylist *database.UserPlaylist
	Playlist     *database.Playlist
}

func UserPlaylistExists(
	ctx context.Context,
	db *database.Queries,
	params database.UserPlaylistExistsByUserIDAndPlaylistIDParams,
) (bool, error) {
	exists, err := db.UserPlaylistExistsByUserIDAndPlaylistID(ctx, params)

	if err != nil {
		log.Println("Error checking user playlist exists by user ID and playlist ID:", err)
	}

	return exists, err
}

func GetPlaylistIDsByUserID(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	playlistIds, err := db.GetPlaylistIDsByUserID(ctx, userId)

	if err != nil {
		log.Println("Error getting user playlist ids:", err)
	}

	return playlistIds, err
}

func MapUserPlaylistFullEntityToDTO(e *database.GetFullUserPlaylistsRow) *UserPlaylistDTO {
	return &UserPlaylistDTO{
		ID:                     e.UserPlaylistID,
		PlaylistID:             e.UserPlaylistPlaylistID,
		UserID:                 e.UserPlaylistUserID,
		CreatedAt:              e.UserPlaylistCreatedAt,
		IsSpotifySavedPlaylist: e.UserPlaylistIsSpotifySavedPlaylist,
		Playlist: &PlaylistDTO{
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

func UserPlaylistEntityToDTO(e *database.UserPlaylist) *UserPlaylistDTO {
	return &UserPlaylistDTO{
		ID:                     e.ID,
		PlaylistID:             e.PlaylistID,
		UserID:                 e.UserID,
		CreatedAt:              e.CreatedAt,
		IsSpotifySavedPlaylist: e.IsSpotifySavedPlaylist,
		Playlist:               nil,
	}
}

func UserPlaylistWithRelToDTO(rel *UserPlaylistWithRel) *UserPlaylistWithRelDTO {
	return &UserPlaylistWithRelDTO{
		UserPlaylistDTO: UserPlaylistEntityToDTO(rel.UserPlaylist),
		Playlist:        PlaylistEntityToDTO(rel.Playlist),
	}
}
