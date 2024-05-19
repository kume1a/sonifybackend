package audio

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/useraudio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func CreateAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreateAudioParams,
) (*database.Audio, error) {
	entity, err := createAudio(ctx, db, params)

	if err != nil {
		log.Println("Error creating audio:", err)
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return entity, err
}

func BulkCreateAudios(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
	params []database.CreateAudioParams,
) ([]database.Audio, error) {
	return shared.RunDBTransaction(
		ctx,
		apiCfg,
		func(tx *database.Queries) ([]database.Audio, error) {
			audios := make([]database.Audio, 0, len(params))

			for _, param := range params {
				audio, err := createAudio(ctx, tx, param)
				if err != nil {
					log.Println("Error creating audio:", err)
					return nil, shared.HttpErrInternalServerErrorDef()
				}

				audios = append(audios, *audio)
			}

			return audios, nil
		},
	)
}

func DoesAudioExistByLocalId(ctx context.Context, db *database.Queries, userID uuid.UUID, localID string) (bool, error) {
	count, err := useraudio.CountUserAudioByLocalId(ctx, db, database.CountUserAudioByLocalIDParams{
		LocalID: sql.NullString{String: localID, Valid: true},
		UserID:  userID,
	})

	if err != nil {
		log.Println("Error counting user audio by local id: ", err)
		return false, shared.HttpErrInternalServerErrorDef()
	}

	return count > 0, nil
}

func GetAudioSpotifyIdsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	spotifyIds []string,
) ([]database.GetAudioSpotifyIDsBySpotifyIDsRow, error) {
	ids, err := getAudioSpotifyIdsBySpotifyIds(ctx, db, spotifyIds)

	if err != nil {
		log.Println("Error getting audios spotify ids by spotify ids: ", err)
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return ids, err
}

func GetAudioIdsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	spotifyIds []string,
) (uuid.UUIDs, error) {
	ids, err := getAudioIdsBySpotifyIds(ctx, db, spotifyIds)

	if err != nil {
		log.Println("Error getting audio ids by spotify ids: ", err)
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return ids, err
}
