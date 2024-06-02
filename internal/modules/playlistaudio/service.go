package playlistaudio

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/userplaylist"
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

	entity, err := db.CreatePlaylistAudio(ctx, params)

	if err != nil {
		log.Println("Error creating playlist audio:", err)
	}

	return &entity, err
}

func BulkCreatePlaylistAudios(
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

func DeletePlaylistAudiosByIds(
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

func GetPlaylistAudiosWithAudioAndAudioLikes(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistAudiosWithAudioAndAudioLikesParams,
) ([]database.GetPlaylistAudiosWithAudioAndAudioLikesRow, error) {
	audios, err := db.GetPlaylistAudiosWithAudioAndAudioLikes(ctx, params)

	if err != nil {
		log.Println("Error getting playlist audios:", err)
	}

	return audios, err
}

func GetPlaylistAudiosByUserID(
	ctx context.Context,
	db *database.Queries,
	userID uuid.UUID,
	filterIDs uuid.UUIDs,
) ([]database.PlaylistAudio, error) {
	userPlaylistIDs, err := userplaylist.GetUserPlaylistIDsByUserID(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	playlistAudios, err := db.GetPlaylistAudios(ctx, database.GetPlaylistAudiosParams{
		PlaylistIds: userPlaylistIDs,
		Ids:         filterIDs,
	})
	if err != nil {
		log.Println("Error getting playlist audios by user ID:", err)
		return nil, err
	}

	return playlistAudios, nil
}

func GetPlaylistAudioIDsByUserID(
	ctx context.Context,
	db *database.Queries,
	userID uuid.UUID,
) (uuid.UUIDs, error) {
	userPlaylistIDs, err := userplaylist.GetUserPlaylistIDsByUserID(ctx, db, userID)
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
