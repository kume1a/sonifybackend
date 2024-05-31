package playlist

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func CreatePlaylist(
	ctx context.Context,
	db *database.Queries,
	params database.CreatePlaylistParams,
) (*database.Playlist, error) {
	if params.ID == uuid.Nil {
		params.ID = uuid.New()
	}
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	entity, err := db.CreatePlaylist(ctx, params)

	if err != nil {
		log.Println("Error creating playlist:", err)
	}

	return &entity, err
}

func GetPlaylists(
	ctx context.Context,
	db *database.Queries,
	params database.GetPlaylistsParams,
) ([]database.Playlist, error) {
	playlists, err := db.GetPlaylists(ctx, params)

	if err != nil {
		log.Println("Error getting playlists:", err)
	}

	return playlists, err
}

func GetSpotifyUserSavedPlaylistIds(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) (uuid.UUIDs, error) {
	playlistIds, err := db.GetSpotifyUserSavedPlaylistIDs(ctx, userId)

	if err != nil {
		log.Println("Error getting spotify user saved playlist ids:", err)
	}

	return playlistIds, err
}

func DeleteSpotifyUserSavedPlaylistJoins(
	ctx context.Context,
	db *database.Queries,
	userId uuid.UUID,
) error {
	err := db.DeleteSpotifyUserSavedPlaylistJoins(ctx, userId)

	if err != nil {
		log.Println("Error deleting spotify user saved playlist ids:", err)
	}

	return err
}

func DeletePlaylistsByIds(
	ctx context.Context,
	db *database.Queries,
	playlistIds uuid.UUIDs,
) error {
	err := db.DeletePlaylistsByIDs(ctx, playlistIds)

	if err != nil {
		log.Println("Error deleting playlists by ids:", err)
	}

	return err
}

func GetPlaylistById(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
) (*database.Playlist, error) {
	playlist, err := db.GetPlaylistByID(ctx, playlistID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return nil, shared.NotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		log.Println("Error getting playlist by id:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &playlist, nil
}

func GetPlaylistIDBySpotifyID(
	ctx context.Context,
	db *database.Queries,
	spotifyID string,
) (uuid.UUID, error) {
	playlistID, err := db.GetPlaylistIDBySpotifyID(ctx, spotifyID)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return uuid.Nil, shared.NotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		return uuid.Nil, shared.InternalServerErrorDef()
	}

	return playlistID, nil
}

func UpdatePlaylistByID(
	ctx context.Context,
	db *database.Queries,
	params database.UpdatePlaylistByIDParams,
) (*database.Playlist, error) {
	playlist, err := db.UpdatePlaylistByID(ctx, params)

	if err != nil && shared.IsDBErrorNotFound(err) {
		return nil, shared.NotFound(shared.ErrPlaylistNotFound)
	}

	if err != nil {
		log.Println("Error updating playlist by id:", err)
		return nil, shared.InternalServerErrorDef()
	}

	return &playlist, err

}

func GetPlaylistWithAudios(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
	authUserID uuid.UUID,
) (*database.Playlist, []audio.AudioWithAudioLike, error) {
	playlist, err := GetPlaylistById(ctx, db, playlistID)

	if err != nil {
		return nil, nil, err
	}

	playlistAudios, err := playlistaudio.GetPlaylistAudios(
		ctx, db,
		database.GetPlaylistAudiosParams{
			PlaylistID: playlistID,
			UserID:     authUserID,
		},
	)
	if err != nil {
		return nil, nil, shared.InternalServerErrorDef()
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
