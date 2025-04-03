package playlistaudio

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func CreatePlaylistAudio(
	ctx context.Context,
	db *database.Queries,
	params database.CreatePlaylistAudioParams,
) (*database.PlaylistAudio, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	exists, err := db.PlaylistAudioExists(ctx, database.PlaylistAudioExistsParams{
		PlaylistID: params.PlaylistID,
		AudioID:    params.AudioID,
	})
	if err != nil {
		log.Println("Error checking if playlist audio exists:", err)
		return nil, shared.InternalServerErrorDef()
	}

	if exists {
		return nil, shared.Conflict(shared.ErrPlaylistAudioAlreadyExists)
	}

	entity, err := db.CreatePlaylistAudio(ctx, params)
	if err != nil {
		log.Println("Error creating playlist audio:", err)
		return nil, shared.InternalServerErrorDef()
	}

	if err := sharedmodule.IncrementAudioPlaylistAudioCountByID(
		ctx, db, params.AudioID,
	); err != nil {
		log.Println("Error incrementing playlist audio count:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &entity, err
}

func CreatePlaylistAudioTx(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	params database.CreatePlaylistAudioParams,
) (*database.PlaylistAudio, error) {
	return shared.RunDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) (*database.PlaylistAudio, error) {
			return CreatePlaylistAudio(ctx, tx, params)
		},
	)
}

func BulkCreatePlaylistAudiosTx(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	params []database.CreatePlaylistAudioParams,
) ([]database.PlaylistAudio, error) {
	return shared.RunDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) ([]database.PlaylistAudio, error) {
			entities := make([]database.PlaylistAudio, 0, len(params))

			for _, param := range params {
				entity, err := CreatePlaylistAudio(ctx, tx, param)
				if err != nil {
					log.Println("Error creating playlist audio:", err)
					return nil, shared.InternalServerErrorDef()
				}

				entities = append(entities, *entity)
			}

			return entities, nil
		},
	)
}

func DeletePlaylistAudiosByIDs(
	ctx context.Context,
	db *database.Queries,
	params database.DeletePlaylistAudiosByIDsParams,
) error {
	err := db.DeletePlaylistAudiosByIDs(ctx, params)

	if err != nil {
		log.Println("Error deleting playlist audios by ids:", err)
	}

	return err
}

func DeletePlaylistAudiosByPlaylistID(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
) error {
	err := db.DeletePlaylistAudiosByPlaylistID(ctx, playlistID)

	if err != nil {
		log.Println("Error deleting playlist audios by playlist ID:", err)
	}

	return err
}

func DeletePlaylistAudioByPlaylistIDAndAudioIDTx(
	ctx context.Context,
	resourceConfig *config.ResourceConfig,
	params database.DeletePlaylistAudioByPlaylistIDAndAudioIDParams,
) error {
	return shared.RunNoResultDBTransaction(
		ctx,
		resourceConfig,
		func(tx *database.Queries) error {
			return DeletePlaylistAudioByPlaylistIDAndAudioID(ctx, tx, params)
		},
	)
}

func DeletePlaylistAudioByPlaylistIDAndAudioID(
	ctx context.Context,
	db *database.Queries,
	params database.DeletePlaylistAudioByPlaylistIDAndAudioIDParams,
) error {
	err := db.DeletePlaylistAudioByPlaylistIDAndAudioID(ctx, params)

	if err != nil {
		log.Println("Error deleting playlist audio by playlist ID and audio ID:", err)
	}

	if err := sharedmodule.DecrementAudioPlaylistAudioCountByID(
		ctx, db, params.AudioID,
	); err != nil {
		log.Println("Error decrementing playlist audio count:", err)
		return shared.InternalServerErrorDef()
	}

	return err
}

func GetPlaylistAudioJoinsBySpotifyIds(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistAudioJoinsBySpotifyIDsParams,
) ([]database.GetPlaylistAudioJoinsBySpotifyIDsRow, error) {
	entities, err := db.GetPlaylistAudioJoinsBySpotifyIDs(ctx, params)

	if err != nil {
		log.Println("Error getting playlist audio joins by spotify ids:", err)
	}

	return entities, err
}

func GetPlaylistAudios(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistAudiosParams,
) ([]database.GetPlaylistAudiosRow, error) {
	playlistAudios, err := db.GetPlaylistAudios(ctx, params)

	if err != nil {
		log.Println("Error getting playlist audios by user ID:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return playlistAudios, nil
}

func GetPlaylistAudioIDsByUserID(
	ctx context.Context,
	db *database.Queries,
	userID uuid.UUID,
) (uuid.UUIDs, error) {
	userPlaylistIDs, err := sharedmodule.GetPlaylistIDsByUserID(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	playlistAudioIds, err := db.GetPlaylistAudioIDsByPlaylistIDs(ctx, userPlaylistIDs)
	if err != nil {
		log.Println("Error getting playlist audio IDs by playlist ID:", err)
		return nil, err
	}

	return playlistAudioIds, nil
}

func CountPlaylistAudiosByAudioID(
	ctx context.Context,
	db *database.Queries,
	audioID uuid.UUID,
) (int64, error) {
	count, err := db.CountPlaylistAudiosByAudioID(ctx, audioID)

	if err != nil {
		log.Println("Error counting playlist audios by audio ID:", err)
		return 0, shared.InternalServerErrorDef()
	}

	return count, err
}

func PlaylistAudioExistsByYoutubeVideoID(
	ctx context.Context,
	db *database.Queries,
	params database.PlaylistAudioExistsByYoutubeVideoIDParams,
) (bool, error) {
	row, err := db.PlaylistAudioExistsByYoutubeVideoID(ctx, params)

	if err != nil {
		log.Println("Error checking if playlist audio exists by youtube video ID:", err)
	}

	return row, err
}
