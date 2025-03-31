package playlistaudio

import (
	"context"
	"database/sql"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func (dto *createPlaylistAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *deletePlaylistAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

type ValidatePlaylistAudioDoesNotExistByYoutubeVideoIDParams struct {
	Context        context.Context
	ApiConfig      *config.ApiConfig
	PlaylistID     uuid.UUID
	YoutubeVideoID string
}

func ValidatePlaylistAudioDoesNotExistByYoutubeVideoID(
	params *ValidatePlaylistAudioDoesNotExistByYoutubeVideoIDParams,
) error {
	playlistAudioExists, err := PlaylistAudioExistsByYoutubeVideoID(
		params.Context,
		params.ApiConfig.DB,
		database.PlaylistAudioExistsByYoutubeVideoIDParams{
			PlaylistID:     params.PlaylistID,
			YoutubeVideoID: sql.NullString{String: params.YoutubeVideoID, Valid: true},
		},
	)

	if err != nil {
		return shared.InternalServerErrorDef()
	}

	if playlistAudioExists {
		return shared.Conflict(shared.ErrPlaylistAudioAlreadyExists)
	}

	return nil
}
