package audio

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

func CreateAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreateAudioParams,
) (*database.Audio, error) {
	createdAt := time.Now().UTC()

	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = createdAt
	}

	entity, err := db.CreateAudio(ctx, params)

	if err != nil {
		log.Println("Error creating audio:", err)
	}

	return &entity, err
}

func CreateUserAudio(ctx context.Context, db *database.Queries, params database.CreateUserAudioParams) (*database.UserAudio, error) {
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	entity, err := db.CreateUserAudio(ctx, params)

	if err != nil {
		log.Println("Error creating user audio:", err)
	}

	return &entity, err
}

func GetUserAudioByYoutubeVideoId(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
	youtubeVideoId string,
) (*database.GetUserAudioByVideoIdRow, error) {
	audio, err := db.GetUserAudioByVideoId(ctx, database.GetUserAudioByVideoIdParams{
		UserID:         userId,
		YoutubeVideoID: sql.NullString{String: youtubeVideoId, Valid: true},
	})

	if err != nil {
		log.Println("Error getting user audio by youtube video id: ", err)
	}

	return &audio, err
}

func GetPlaylistAudiosBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistAudiosBySpotifyIdsParams,
) ([]database.Audio, error) {
	audios, err := db.GetPlaylistAudiosBySpotifyIds(ctx, params)

	if err != nil {
		log.Println("Error getting playlist audios by spotify ids: ", err)
	}

	return audios, err
}

func GetAudioSpotifyIdsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	spotifyIds []string,
) ([]database.GetAudioSpotifyIdsBySpotifyIdsRow, error) {
	ids, err := db.GetAudioSpotifyIdsBySpotifyIds(ctx, spotifyIds)

	if err != nil {
		log.Println("Error getting audios spotify ids by spotify ids: ", err)
	}

	return ids, err
}

func GetAudioIdsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	spotifyIds []string,
) (uuid.UUIDs, error) {
	ids, err := db.GetAudioIdsBySpotifyIds(ctx, spotifyIds)

	if err != nil {
		log.Println("Error getting audio ids by spotify ids: ", err)
	}

	return ids, err
}

func CountUserAudioByLocalId(
	ctx context.Context,
	db *database.Queries,
	params database.CountUserAudioByLocalIdParams,
) (int64, error) {
	count, err := db.CountUserAudioByLocalId(ctx, params)

	if err != nil {
		log.Println("Error getting audio by local id: ", err)
	}

	return count, err
}

func GetUserAudios(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) ([]database.Audio, error) {
	audios, err := db.GetAudiosByUserId(ctx, userId)

	if err != nil {
		log.Println("Error getting user audios: ", err)
	}

	return audios, err
}

func GetUserAudioIds(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	ids, err := db.GetUserAudioIds(ctx, userId)

	if err != nil {
		log.Println("Error getting user audio ids: ", err)
	}

	return ids, err
}

func GetAudiosByIds(
	ctx context.Context,
	db *database.Queries,
	ids []uuid.UUID,
) ([]database.Audio, error) {
	audios, err := db.GetAudiosByIds(ctx, ids)

	if err != nil {
		log.Println("Error getting audios by ids: ", err)
	}

	return audios, err
}
