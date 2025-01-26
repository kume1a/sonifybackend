package useraudio

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type ValidateUserAudioDoesNotExistParams struct {
	Context        context.Context
	ApiConfig      *config.ApiConfig
	UserID         uuid.UUID
	YoutubeVideoID string
}

func ValidateUserAudioDoesNotExist(params *ValidateUserAudioDoesNotExistParams) error {
	exists, err := UserAudioExistsByYoutubeVideoID(
		params.Context,
		params.ApiConfig.DB,
		database.UserAudioExistsByYoutubeVideoIDParams{
			UserID:         params.UserID,
			YoutubeVideoID: sql.NullString{String: params.YoutubeVideoID, Valid: true},
		},
	)

	if err != nil {
		return shared.InternalServerErrorDef()
	}

	if exists {
		return shared.Conflict(shared.ErrUserAudioAlreadyExists)
	}

	return nil
}
