package audio

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func ValidateUploadUserLocalMusicDTO(w http.ResponseWriter, r *http.Request) (*uploadUserLocalMusicDTO, error) {
	thumbnailPath, err := shared.HandleUploadFile(shared.HandleUploadFileArgs{
		ResponseWriter:   w,
		Request:          r,
		FieldName:        "thumbnail",
		Dir:              config.DirUserLocalAudioThumbnails,
		AllowedMimeTypes: shared.ImageMimeTypes,
		IsOptional:       true,
	})
	if err != nil {
		return nil, err
	}

	audioPath, err := shared.HandleUploadFile(shared.HandleUploadFileArgs{
		ResponseWriter:   w,
		Request:          r,
		FieldName:        "audio",
		Dir:              config.DirUserLocalAudios,
		AllowedMimeTypes: shared.AudioMimeTypes,
		IsOptional:       false,
	})
	if err != nil {
		return nil, err
	}

	localId := r.FormValue("localId")
	title := r.FormValue("title")
	author := r.FormValue("author")
	durationMs := r.FormValue("durationMs")

	if localId == "" {
		return nil, shared.BadRequest("localId must be provided")
	}

	if !govalidator.IsByteLength(title, 1, 255) {
		return nil, shared.BadRequest("title must be between 1 and 255 characters")
	}

	intDurationMs := 0
	if durationMs != "" {
		intDurationMsValue, err := strconv.Atoi(durationMs)
		if err != nil {
			return nil, shared.BadRequest("durationMs must be an integer")
		}

		intDurationMs = intDurationMsValue
	}

	return &uploadUserLocalMusicDTO{
		LocalID:       localId,
		Title:         title,
		Author:        author,
		AudioPath:     audioPath,
		ThumbnailPath: thumbnailPath,
		DurationMs:    int32(intDurationMs),
	}, nil
}

type ValidateAudioDoesNotExistByYoutubeVideoIDParams struct {
	Context        context.Context
	ApiConfig      *config.ApiConfig
	YoutubeVideoID string
}

func ValidateAudioDoesNotExistByYoutubeVideoID(
	params *ValidateAudioDoesNotExistByYoutubeVideoIDParams,
) error {
	exists, err := AudioExistsByYoutubeVideoID(
		params.Context,
		params.ApiConfig.DB,
		sql.NullString{String: params.YoutubeVideoID, Valid: true},
	)

	if err != nil {
		return shared.InternalServerErrorDef()
	}

	if exists {
		return shared.Conflict(shared.ErrAudioAlreadyExists)
	}

	return nil
}
