package audio

import (
	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/database"
)

func (dto downloadYoutubeAudioDto) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func audioEntityToDto(e *database.Audio) *AudioDto {
	return &AudioDto{
		ID:             e.ID,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		Title:          e.Title.String,
		Duration:       e.Duration.Int32,
		Path:           e.Path.String,
		Author:         e.Author.String,
		SizeBytes:      e.SizeBytes.Int64,
		YoutubeVideoID: e.YoutubeVideoID.String,
		ThumbnailPath:  e.ThumbnailPath.String,
	}
}

func userAudioEntityToDto(e *database.UserAudio) *UserAudioDto {
	return &UserAudioDto{
		UserId:  e.UserID,
		AudioId: e.AudioID,
	}
}
