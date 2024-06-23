package audio

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/useraudio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func CreateAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreateAudioParams,
) (*database.Audio, error) {
	// Trim surrounding quotes from the title
	if params.Title.Valid {
		params.Title.String = strings.Trim(params.Title.String, `"'`)
	}

	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}

	entity, err := db.CreateAudio(ctx, params)

	if err != nil {
		log.Println("Error creating audio:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &entity, err
}

func BulkCreateAudios(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	params []database.CreateAudioParams,
) ([]database.Audio, error) {
	return shared.RunDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) ([]database.Audio, error) {
			audios := make([]database.Audio, 0, len(params))

			for _, param := range params {
				audio, err := CreateAudio(ctx, tx, param)
				if err != nil {
					log.Println("Error creating audio:", err)
					return nil, shared.InternalServerErrorDef()
				}

				audios = append(audios, *audio)
			}

			return audios, nil
		},
	)
}

func DoesAudioExistByLocalId(ctx context.Context, db *database.Queries, userID uuid.UUID, localID string) (bool, error) {
	count, err := useraudio.CountUserAudioByLocalId(
		ctx, db,
		database.CountUserAudioByLocalIDParams{
			LocalID: sql.NullString{String: localID, Valid: true},
			UserID:  userID,
		},
	)

	if err != nil {
		log.Println("Error counting user audio by local id: ", err)
		return false, shared.InternalServerErrorDef()
	}

	return count > 0, nil
}

func GetAudioSpotifyIdsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	spotifyIds []string,
) ([]database.GetAudioSpotifyIDsBySpotifyIDsRow, error) {
	ids, err := db.GetAudioSpotifyIDsBySpotifyIDs(ctx, spotifyIds)

	if err != nil {
		log.Println("Error getting audios spotify ids by spotify ids: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return ids, err
}

func GetAudioIdsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	spotifyIds []string,
) (uuid.UUIDs, error) {
	ids, err := db.GetAudioIDsBySpotifyIDs(ctx, spotifyIds)

	if err != nil {
		log.Println("Error getting audio ids by spotify ids: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return ids, err
}
