package audio

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func (dto downloadYoutubeAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func ValidateImportUserLocalMusicDTO(w http.ResponseWriter, r *http.Request) (*importUserLocalMusicDTO, *shared.HttpError) {
	thumbnailPath, httpErr := shared.HandleUploadFile(
		w, r,
		"thumbnail",
		shared.DirUserLocalAudioThumbnails,
		shared.ImageMimeTypes,
	)
	if httpErr != nil {
		return nil, httpErr
	}

	audioPath, httpErr := shared.HandleUploadFile(
		w, r,
		"audio",
		shared.DirUserLocalAudios,
		shared.AudioMimeTypes,
	)
	if httpErr != nil {
		return nil, httpErr
	}

	localId := r.FormValue("localId")
	title := r.FormValue("title")
	author := r.FormValue("author")
	durationMs := r.FormValue("durationMs")

	if localId == "" {
		return nil, shared.HttpErrBadRequest("localId must be provided")
	}

	if !govalidator.IsByteLength(title, 1, 255) {
		return nil, shared.HttpErrBadRequest("title must be between 1 and 255 characters")
	}

	if !govalidator.IsByteLength(author, 1, 255) {
		return nil, shared.HttpErrBadRequest("author must be between 1 and 255 characters")
	}

	if durationMs == "" {
		return nil, shared.HttpErrBadRequest("durationMs must be provided")
	}

	intDurationMs, err := strconv.Atoi(durationMs)
	if err != nil {
		return nil, shared.HttpErrBadRequest("durationMs must be an integer")
	}

	return &importUserLocalMusicDTO{
		LocalId:       localId,
		Title:         title,
		Author:        author,
		AudioPath:     audioPath,
		ThumbnailPath: thumbnailPath,
		DurationMs:    int32(intDurationMs),
	}, nil
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
	}
}

func UserAudioEntityToDto(e *database.UserAudio) *UserAudioDTO {
	return &UserAudioDTO{
		CreatedAt: e.CreatedAt,
		UserId:    e.UserID,
		AudioId:   e.AudioID,
	}
}
