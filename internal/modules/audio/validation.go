package audio

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func (dto *audioIDsDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func ValidateUploadUserLocalMusicDTO(w http.ResponseWriter, r *http.Request) (*uploadUserLocalMusicDTO, *shared.HttpError) {
	thumbnailPath, httpErr := shared.HandleUploadFile(shared.HandleUploadFileArgs{
		ResponseWriter:   w,
		Request:          r,
		FieldName:        "thumbnail",
		Dir:              config.DirUserLocalAudioThumbnails,
		AllowedMimeTypes: shared.ImageMimeTypes,
		IsOptional:       true,
	})
	if httpErr != nil {
		return nil, httpErr
	}

	audioPath, httpErr := shared.HandleUploadFile(shared.HandleUploadFileArgs{
		ResponseWriter:   w,
		Request:          r,
		FieldName:        "audio",
		Dir:              config.DirUserLocalAudios,
		AllowedMimeTypes: shared.AudioMimeTypes,
		IsOptional:       false,
	})
	if httpErr != nil {
		return nil, httpErr
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
