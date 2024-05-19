package useraudio

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateUserAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreateUserAudioParams,
) (*database.UserAudio, error) {
	entity, err := db.CreateUserAudio(ctx, params)

	if err != nil {
		log.Println("Error creating user audio:", err)
	}

	return &entity, err
}

func GetUserAudioByYoutubeVideoId(
	ctx context.Context,
	db *database.Queries,
	params database.GetUserAudioByVideoIDParams,
) (*database.GetUserAudioByVideoIDRow, error) {
	audio, err := db.GetUserAudioByVideoID(ctx, params)

	if err != nil {
		log.Println("Error getting user audio by youtube video id: ", err)
	}

	return &audio, err
}

func CountUserAudioByLocalId(
	ctx context.Context,
	db *database.Queries,
	params database.CountUserAudioByLocalIDParams,
) (int64, error) {
	count, err := db.CountUserAudioByLocalID(ctx, params)

	if err != nil {
		log.Println("Error getting audio by local id: ", err)
	}

	return count, err
}

func GetUserAudioIds(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	ids, err := db.GetUserAudioIDs(ctx, userId)

	if err != nil {
		log.Println("Error getting user audio ids: ", err)
	}

	return ids, err
}

func GetUserAudiosByAudioIds(
	ctx context.Context,
	db *database.Queries,
	params database.GetUserAudiosByAudioIdsParams,
) ([]database.GetUserAudiosByAudioIdsRow, error) {
	if len(params.AudioIds) == 0 {
		return []database.GetUserAudiosByAudioIdsRow{}, nil
	}

	audios, err := db.GetUserAudiosByAudioIds(ctx, params)

	if err != nil {
		log.Println("Error getting audios by ids: ", err)
	}

	return audios, err
}
