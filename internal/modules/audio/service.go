package audio

import (
	"context"
	"database/sql"
	"log"
	"strings"

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

func UpdateAudioByID(
	ctx context.Context,
	db *database.Queries,
	params database.UpdateAudioByIDParams,
) (*database.Audio, error) {
	entity, err := db.UpdateAudioByID(ctx, params)

	if err != nil {
		log.Println("Error updating audio by id:", err)

		if shared.IsDBErrorNotFound(err) {
			return nil, shared.NotFound(shared.ErrAudioNotFound)
		}

		return nil, shared.InternalServerErrorDef()
	}

	return &entity, nil
}

func DeleteAudioByID(
	ctx context.Context,
	db *database.Queries,
	id uuid.UUID,
) error {
	err := db.DeleteAudioByID(ctx, id)

	if err != nil {
		if shared.IsDBErrorNotFound(err) {
			return shared.NotFound(shared.ErrAudioNotFound)
		}

		log.Println("Error deleting unused audios:", err)
		return shared.InternalServerErrorDef()
	}

	return nil
}

func DoesAudioExistByLocalId(
	ctx context.Context,
	db *database.Queries,
	userID uuid.UUID,
	localID string,
) (bool, error) {
	count, err := useraudio.CountUserAudioByLocalID(
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

func GetAudioIDsBySpotifyIDs(
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

func GetAllAudioIDs(
	ctx context.Context,
	db *database.Queries,
) (uuid.UUIDs, error) {
	ids, err := db.GetAllAudioIDs(ctx)

	if err != nil {
		log.Println("Error getting all audio ids: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return ids, err
}

func GetUnusedAudios(
	ctx context.Context,
	db *database.Queries,
) ([]database.Audio, error) {
	audios, err := db.GetUnusedAudios(ctx)

	if err != nil {
		log.Println("Error getting unused audios: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return audios, nil
}

func AudioExistsByYoutubeVideoID(
	ctx context.Context,
	db *database.Queries,
	youtubeVideoID sql.NullString,
) (bool, error) {
	row, err := db.AudioExistsByYoutubeVideoID(ctx, youtubeVideoID)

	if err != nil {
		log.Println("Error checking if audio exists by youtube video id: ", err)
	}

	return row, err
}

func GetAudioByYoutubeVideoID(
	ctx context.Context,
	db *database.Queries,
	youtubeVideoID sql.NullString,
) (*database.Audio, error) {
	audio, err := db.GetAudioByYoutubeVideoID(ctx, youtubeVideoID)

	if err != nil {
		log.Println("Error getting audio by youtube video id: ", err)
	}

	return &audio, nil
}
