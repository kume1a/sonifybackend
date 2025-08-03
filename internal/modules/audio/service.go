package audio

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

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

	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	entity, err := db.CreateAudio(ctx, params)

	if err != nil {
		log.Println("Error creating audio:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &entity, nil
}

func UpdateAudioByID(
	ctx context.Context,
	db *database.Queries,
	params database.UpdateAudioByIDParams,
) (*database.Audio, error) {
	entity, err := db.UpdateAudioByID(ctx, params)

	if err != nil {
		log.Println("Error updating audio by ID:", err)

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

func DoesAudioExistByLocalID(
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
		log.Println("Error counting user audio by local ID: ", err)
		return false, shared.InternalServerErrorDef()
	}

	return count > 0, nil
}

func GetAudioSpotifyIDsBySpotifyIDs(
	ctx context.Context,
	db *database.Queries,
	spotifyIDs []string,
) ([]database.GetAudioSpotifyIDsBySpotifyIDsRow, error) {
	ids, err := db.GetAudioSpotifyIDsBySpotifyIDs(ctx, spotifyIDs)

	if err != nil {
		log.Println("Error getting audios spotify IDs by spotify IDs: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return ids, nil
}

func GetAudioIDsBySpotifyIDs(
	ctx context.Context,
	db *database.Queries,
	spotifyIDs []string,
) (uuid.UUIDs, error) {
	ids, err := db.GetAudioIDsBySpotifyIDs(ctx, spotifyIDs)

	if err != nil {
		log.Println("Error getting audio IDs by spotify IDs: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return ids, nil
}

func GetAllAudioIDs(
	ctx context.Context,
	db *database.Queries,
) (uuid.UUIDs, error) {
	ids, err := db.GetAllAudioIDs(ctx)

	if err != nil {
		log.Println("Error getting all audio IDs: ", err)
		return nil, shared.InternalServerErrorDef()
	}

	return ids, nil
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
		log.Println("Error checking if audio exists by youtube video ID: ", err)
	}

	return row, err
}

func GetAudioByYoutubeVideoID(
	ctx context.Context,
	db *database.Queries,
	youtubeVideoID string,
) (*database.Audio, error) {
	audio, err := db.GetAudioByYoutubeVideoID(ctx, youtubeVideoID)

	if err != nil {
		log.Println("Error getting audio by youtube video ID: ", err)
	}

	return &audio, err
}
