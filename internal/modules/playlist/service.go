package playlist

import (
	"context"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

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

	playlistAudios, err := getPlaylistAudios(ctx, db, database.GetPlaylistAudiosParams{
		PlaylistID: playlistID,
		UserID:     authUserID,
	})
	if err != nil {
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	audiosWithLike := make([]audio.AudioWithAudioLike, len(playlistAudios))
	for index, playlistAudio := range playlistAudios {
		var audioLike *database.AudioLike
		if playlistAudio.AudioLikesAudioID.Valid && playlistAudio.AudioLikesUserID.Valid {
			audioLike = &database.AudioLike{
				AudioID: playlistAudio.AudioLikesAudioID.UUID,
				UserID:  playlistAudio.AudioLikesUserID.UUID,
			}
		}

		audioWithLike := audio.AudioWithAudioLike{
			Audio: &database.Audio{
				ID:             playlistAudio.ID,
				Title:          playlistAudio.Title,
				Author:         playlistAudio.Author,
				DurationMs:     playlistAudio.DurationMs,
				Path:           playlistAudio.Path,
				CreatedAt:      playlistAudio.CreatedAt,
				SizeBytes:      playlistAudio.SizeBytes,
				YoutubeVideoID: playlistAudio.YoutubeVideoID,
				ThumbnailPath:  playlistAudio.ThumbnailPath,
				SpotifyID:      playlistAudio.SpotifyID,
				ThumbnailUrl:   playlistAudio.ThumbnailUrl,
				LocalID:        playlistAudio.LocalID,
			},
			AudioLike: audioLike,
		}

		audiosWithLike[index] = audioWithLike
	}

	return playlist, audiosWithLike, nil
}
