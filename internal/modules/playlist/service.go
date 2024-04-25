package playlist

import (
	"context"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func GetPlaylistAudios(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
) ([]database.Audio, *shared.HttpError) {
	audios, err := getPlaylistAudios(ctx, db, playlistID)
	if err != nil {
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return audios, nil
}

func GetPlaylistById(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
) (*database.Playlist, *shared.HttpError) {
	playlist, err := getPlaylistById(ctx, db, playlistID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return nil, shared.HttpErrNotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return playlist, nil
}

func GetPlaylistWithAudios(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
	authUserID uuid.UUID,
) (*database.Playlist, []audio.AudioWithAudioLike, *shared.HttpError) {
	playlist, err := getPlaylistById(ctx, db, playlistID)

	if err != nil {
		if shared.IsDBErrorNotFound(err) {
			return nil, nil, shared.HttpErrNotFound(shared.ErrPlaylistNotFound)
		}

		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	audios, err := getPlaylistAudios(ctx, db, database.GetPlaylistAudiosParams{
		PlaylistID: playlistID,
		UserID:     authUserID,
	})
	if err != nil {
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	audiosWithLike := make([]audio.AudioWithAudioLike, len(audios))
	for _, audio := range audios {
		var audioLike *database.AudioLike
		if audio.AudioLikesAudioID.Valid && audio.AudioLikesUserID.Valid {
			audioLike = &database.AudioLike{
				AudioID: audio.AudioLikesAudioID.UUID,
				UserID:  audio.AudioLikesUserID.UUID,
			}
		}

		audioWithLike := audio.AudioWithAudioLike{
			Audio: &database.Audio{
				ID:             audio.ID,
				Title:          audio.Title,
				Author:         audio.Author,
				DurationMs:     audio.DurationMs,
				Path:           audio.Path,
				CreatedAt:      audio.CreatedAt,
				SizeBytes:      audio.SizeBytes,
				YoutubeVideoID: audio.YoutubeVideoID,
				ThumbnailPath:  audio.ThumbnailPath,
				SpotifyID:      audio.SpotifyID,
				ThumbnailUrl:   audio.ThumbnailUrl,
				LocalID:        audio.LocalID,
			},
			AudioLike: audioLike,
		}

		audiosWithLike = append(audiosWithLike, audioWithLike)
	}

	return playlist, audiosWithLike, nil
}
