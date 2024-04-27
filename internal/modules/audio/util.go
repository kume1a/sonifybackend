package audio

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func (dto *downloadYoutubeAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *getAudiosByIdsDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *likeAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *unlikeAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *syncAudioLikesDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func ValidateUploadUserLocalMusicDTO(w http.ResponseWriter, r *http.Request) (*uploadUserLocalMusicDTO, *shared.HttpError) {
	thumbnailPath, httpErr := shared.HandleUploadFile(shared.HandleUploadFileArgs{
		ResponseWriter:   w,
		Request:          r,
		FieldName:        "thumbnail",
		Dir:              shared.DirUserLocalAudioThumbnails,
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
		Dir:              shared.DirUserLocalAudios,
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
		return nil, shared.HttpErrBadRequest("localId must be provided")
	}

	if !govalidator.IsByteLength(title, 1, 255) {
		return nil, shared.HttpErrBadRequest("title must be between 1 and 255 characters")
	}

	intDurationMs := 0
	if durationMs != "" {
		intDurationMsValue, err := strconv.Atoi(durationMs)
		if err != nil {
			return nil, shared.HttpErrBadRequest("durationMs must be an integer")
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

func AudioWithAudioLikeToAudioDTO(e AudioWithAudioLike) *AudioDTO {
	return &AudioDTO{
		ID:             e.Audio.ID,
		CreatedAt:      e.Audio.CreatedAt,
		Title:          e.Audio.Title.String,
		DurationMs:     e.Audio.DurationMs.Int32,
		Path:           e.Audio.Path.String,
		Author:         e.Audio.Author.String,
		SizeBytes:      e.Audio.SizeBytes.Int64,
		YoutubeVideoID: e.Audio.YoutubeVideoID.String,
		ThumbnailPath:  e.Audio.ThumbnailPath.String,
		ThumbnailUrl:   e.Audio.ThumbnailUrl.String,
		SpotifyID:      e.Audio.SpotifyID.String,
		LocalID:        e.Audio.LocalID.String,
		AudioLike:      AudioLikeEntityToDTO(e.AudioLike),
	}
}

func UserAudioEntityToDto(e *database.UserAudio) *UserAudioDTO {
	return &UserAudioDTO{
		CreatedAt: e.CreatedAt,
		UserId:    e.UserID,
		AudioId:   e.AudioID,
	}
}

func GetUserAudiosByAudioIdsRowToUserAudioWithRelDTO(
	e database.GetUserAudiosByAudioIdsRow,
) *UserAudioWithRelDTO {
	var audioLike *AudioLikeDTO
	if e.AudioLikesUserID.Valid && e.AudioLikesAudioID.Valid {
		audioLike = &AudioLikeDTO{
			AudioID: e.AudioLikesAudioID.UUID,
			UserID:  e.AudioLikesUserID.UUID,
		}
	}

	return &UserAudioWithRelDTO{
		UserAudioDTO: &UserAudioDTO{
			CreatedAt: e.CreatedAt,
			UserId:    e.UserID,
			AudioId:   e.AudioID,
		},
		Audio: &AudioDTO{
			ID:             e.AudioID,
			CreatedAt:      e.AudioCreatedAt,
			Title:          e.AudioTitle.String,
			DurationMs:     e.AudioDurationMs.Int32,
			Path:           e.AudioPath.String,
			Author:         e.AudioAuthor.String,
			SizeBytes:      e.AudioSizeBytes.Int64,
			YoutubeVideoID: e.AudioYoutubeVideoID.String,
			ThumbnailPath:  e.AudioThumbnailPath.String,
			ThumbnailUrl:   e.AudioThumbnailUrl.String,
			SpotifyID:      e.AudioSpotifyID.String,
			LocalID:        e.AudioLocalID.String,
			AudioLike:      audioLike,
		},
	}
}

func AudioLikeEntityToDTO(e *database.AudioLike) *AudioLikeDTO {
	if e == nil {
		return nil
	}

	return &AudioLikeDTO{
		UserID:  e.UserID,
		AudioID: e.AudioID,
	}
}
