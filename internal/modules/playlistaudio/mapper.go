package playlistaudio

import (
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
)

func GetPlaylistAudioRowToDTO(
	e database.GetPlaylistAudiosRow,
) sharedmodule.PlaylistAudioDTO {
	var audioLikeDTO *sharedmodule.AudioLikeDTO
	if e.AudioLikesID.Valid {
		audioLikeDTO = &sharedmodule.AudioLikeDTO{
			ID:        e.AudioLikesID.UUID,
			CreatedAt: e.AudioLikesCreatedAt.Time,
			UserID:    e.AudioLikesUserID.UUID,
			AudioID:   e.AudioLikesAudioID.UUID,
		}
	}

	var audioDTO *sharedmodule.AudioDTO
	if e.AudioID.Valid {
		audioDTO = &sharedmodule.AudioDTO{
			ID:             e.AudioID.UUID,
			Title:          e.AudioTitle.String,
			Author:         e.AudioAuthor.String,
			CreatedAt:      e.AudioCreatedAt.Time,
			DurationMs:     e.AudioDurationMs.Int32,
			Path:           e.AudioPath.String,
			SizeBytes:      e.AudioSizeBytes.Int64,
			YoutubeVideoID: e.AudioYoutubeVideoID.String,
			ThumbnailPath:  e.AudioThumbnailPath.String,
			SpotifyID:      e.AudioSpotifyID.String,
			ThumbnailUrl:   e.AudioThumbnailUrl.String,
			LocalID:        e.AudioLocalID.String,
			AudioLike:      audioLikeDTO,
		}
	}

	return sharedmodule.PlaylistAudioDTO{
		ID:         e.PlaylistAudioID,
		CreatedAt:  e.PlaylistAudioCreatedAt,
		PlaylistID: e.PlaylistAudioPlaylistID,
		AudioID:    e.PlaylistAudioAudioID,
		Audio:      audioDTO,
	}
}

func getPlaylistAudioRowListToDTO(
	e []database.GetPlaylistAudiosRow,
) []sharedmodule.PlaylistAudioDTO {
	dto := make([]sharedmodule.PlaylistAudioDTO, 0, len(e))
	for _, v := range e {
		dto = append(dto, GetPlaylistAudioRowToDTO(v))
	}
	return dto
}
