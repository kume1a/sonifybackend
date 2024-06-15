package sharedmodule

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type AudioDTO struct {
	ID             uuid.UUID     `json:"id"`
	CreatedAt      time.Time     `json:"createdAt"`
	Title          string        `json:"title"`
	DurationMs     int32         `json:"durationMs"`
	Path           string        `json:"path"`
	Author         string        `json:"author"`
	SizeBytes      int64         `json:"sizeBytes"`
	YoutubeVideoID string        `json:"youtubeVideoId"`
	ThumbnailPath  string        `json:"thumbnailPath"`
	ThumbnailUrl   string        `json:"thumbnailUrl"`
	SpotifyID      string        `json:"spotifyId"`
	LocalID        string        `json:"localId"`
	AudioLike      *AudioLikeDTO `json:"audioLike"`
}

func AudioEntityToDto(e database.Audio) *AudioDTO {
	return &AudioDTO{
		ID:             e.ID,
		CreatedAt:      e.CreatedAt,
		Title:          e.Title.String,
		DurationMs:     e.DurationMs.Int32,
		Path:           e.Path.String,
		Author:         e.Author.String,
		SizeBytes:      e.SizeBytes.Int64,
		YoutubeVideoID: e.YoutubeVideoID.String,
		ThumbnailPath:  e.ThumbnailPath.String,
		ThumbnailUrl:   e.ThumbnailUrl.String,
		SpotifyID:      e.SpotifyID.String,
		LocalID:        e.LocalID.String,
		AudioLike:      nil,
	}
}

func ValidateAudioExistsByID(
	ctx context.Context,
	db *database.Queries,
	audioID uuid.UUID,
) error {
	count, err := db.CountAudioByID(ctx, audioID)

	if err != nil {
		log.Println("Error validating audio exists by ID: ", err)
		return shared.InternalServerErrorDef()
	}

	if count == 0 {
		return shared.NotFound(shared.ErrAudioNotFound)
	}

	return nil
}
