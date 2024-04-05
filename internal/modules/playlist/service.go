package playlist

import (
	"context"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
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
) (*database.Playlist, []database.Audio, *shared.HttpError) {
	playlist, err := getPlaylistById(ctx, db, playlistID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return nil, nil, shared.HttpErrNotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	audios, err := getPlaylistAudios(ctx, db, playlistID)
	if err != nil {
		return nil, nil, shared.HttpErrInternalServerErrorDef()
	}

	return playlist, audios, nil
}
