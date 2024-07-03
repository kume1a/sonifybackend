package audio

import (
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
)

func GetUserAudiosByAudioIdsRowToUserAudioWithRelDTO(
	e database.GetUserAudiosByAudioIdsRow,
) *sharedmodule.UserAudioWithRelDTO {
	var audioLike *sharedmodule.AudioLikeDTO
	if e.AudioLikesUserID.Valid && e.AudioLikesAudioID.Valid {
		audioLike = &sharedmodule.AudioLikeDTO{
			ID:        e.AudioLikesID.UUID,
			CreatedAt: e.AudioLikesCreatedAt.Time,
			AudioID:   e.AudioLikesAudioID.UUID,
			UserID:    e.AudioLikesUserID.UUID,
		}
	}

	return &sharedmodule.UserAudioWithRelDTO{
		UserAudioDTO: &sharedmodule.UserAudioDTO{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UserId:    e.UserID,
			AudioId:   e.AudioID,
		},
		Audio: &sharedmodule.AudioDTO{
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
