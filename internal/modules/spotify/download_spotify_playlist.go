package spotify

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/modules/userplaylist"
	"github.com/kume1a/sonifybackend/internal/modules/ws"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func DownloadSpotifyPlaylistAudios(
	ctx context.Context,
	resouceConfig *config.ResourceConfig,
	playlistSpotifyID string,
	spotifyAccessToken string,
) error {
	playlistID, err := playlist.GetPlaylistIDBySpotifyID(ctx, resouceConfig.DB, playlistSpotifyID)
	if err != nil {
		setPlaylistImportStatusToFailed(ctx, resouceConfig.DB, playlistID)
		return err
	}

	playlistItems, err := getAllSpotifyPlaylistItems(spotifyAccessToken, playlistSpotifyID)
	if err != nil {
		setPlaylistImportStatusToFailed(ctx, resouceConfig.DB, playlistID)
		return err
	}

	userPlaylistUserIDs, err := userplaylist.GetUserPlaylistUserIDsByPlaylistID(ctx, resouceConfig.DB, playlistID)
	if err != nil {
		setPlaylistImportStatusToFailed(ctx, resouceConfig.DB, playlistID)
		return err
	}

	if err := DownloadWriteSpotifyAudios(
		ctx,
		resouceConfig,
		shared.Map(
			playlistItems,
			func(playlistItem spotifyPlaylistItemDTO) DownloadSpotifyAudioInput {
				artistName := ""
				if len(playlistItem.Track.Artists) > 0 {
					artistName = playlistItem.Track.Artists[0].Name
				}

				thumbnailURL := ""
				if len(playlistItem.Track.Album.Images) > 0 {
					thumbnailURL = playlistItem.Track.Album.Images[0].URL
				}

				return DownloadSpotifyAudioInput{
					SpotifyID:    playlistItem.Track.ID,
					TrackName:    playlistItem.Track.Name,
					ArtistName:   artistName,
					DurationMs:   int32(playlistItem.Track.DurationMS),
					ThumbnailURL: thumbnailURL,
				}
			},
		),
		func(progress, total int) {
			audioImportStatus := database.ProcessStatusPROCESSING
			if progress == total {
				audioImportStatus = database.ProcessStatusCOMPLETED
			}

			updatedPlaylist, err := playlist.UpdatePlaylistByID(
				ctx, resouceConfig.DB,
				database.UpdatePlaylistByIDParams{
					PlaylistID:        playlistID,
					AudioCount:        sql.NullInt32{Int32: int32(progress), Valid: true},
					TotalAudioCount:   sql.NullInt32{Int32: int32(total), Valid: true},
					AudioImportStatus: database.NullProcessStatus{ProcessStatus: audioImportStatus, Valid: true},
				},
			)
			if err != nil {
				log.Println("Error updating playlist by id:", err)
				return
			}

			for _, userPlaylistUserID := range userPlaylistUserIDs {
				updatedPlaylistDTO := sharedmodule.PlaylistEntityToDto(*updatedPlaylist)

				ws.SendWSPayload(ws.SendWSPayloadInput{
					ToSocketId:  userPlaylistUserID.String(),
					MessageType: ws.MessageTypePlaylistUpdated,
					Payload:     updatedPlaylistDTO,
				})
			}
		},
	); err != nil {
		log.Println("Error downloading Spotify playlist audios: ", err)
		setPlaylistImportStatusToFailed(ctx, resouceConfig.DB, playlistID)
		return err
	}

	playlistItemSpotifyIDs := shared.Map(
		playlistItems,
		func(playlistItem spotifyPlaylistItemDTO) string { return playlistItem.Track.ID },
	)

	playlistAudioIDs, err := audio.GetAudioIDsBySpotifyIDs(ctx, resouceConfig.DB, playlistItemSpotifyIDs)
	if err != nil {
		setPlaylistImportStatusToFailed(ctx, resouceConfig.DB, playlistID)
		return err
	}

	createPlaylistAudioParams := []database.CreatePlaylistAudioParams{}
	for _, playlistAudioID := range playlistAudioIDs {
		createPlaylistAudioParams = append(
			createPlaylistAudioParams,
			database.CreatePlaylistAudioParams{
				PlaylistID: playlistID,
				AudioID:    playlistAudioID,
			},
		)
	}

	if _, err := playlistaudio.BulkCreatePlaylistAudiosTx(
		ctx,
		resouceConfig,
		createPlaylistAudioParams,
	); err != nil {
		setPlaylistImportStatusToFailed(ctx, resouceConfig.DB, playlistID)
		return err
	}

	return nil
}

func getAllSpotifyPlaylistItems(
	spotifyAccessToken string,
	playlistSpotifyID string,
) ([]spotifyPlaylistItemDTO, error) {
	allPlaylistItems := []spotifyPlaylistItemDTO{}

	playlistItems, err := SpotifyGetPlaylistItems(spotifyAccessToken, playlistSpotifyID)
	if err != nil {
		return []spotifyPlaylistItemDTO{}, err
	}

	allPlaylistItems = append(allPlaylistItems, playlistItems.Items...)

	nextURL := playlistItems.Next

	for nextURL != "" {
		playlistItems, err = SpotifyGetPlaylistItemsNext(spotifyAccessToken, nextURL)
		if err != nil {
			return nil, err
		}

		allPlaylistItems = append(allPlaylistItems, playlistItems.Items...)
		nextURL = playlistItems.Next
	}

	return allPlaylistItems, nil
}

func downloadSpotifyPlaylist(
	ctx context.Context,
	apiCfg *config.ApiConfig,
	authUserID uuid.UUID,
	spotifyPlaylistID string,
	spotifyAccessToken string,
) (*database.Playlist, error) {
	spotifyPlaylist, err := SpotifyGetPlaylist(spotifyAccessToken, spotifyPlaylistID)
	if err != nil {
		log.Println("Error getting Spotify playlist: ", err)
		return nil, err
	}

	if err := playlist.DeletePlaylistAndPlaylistAudiosBySpotifyID(
		ctx, apiCfg.ResourceConfig,
		spotifyPlaylistID,
	); err != nil {
		log.Println("Error deleting playlist and playlist audios by spotify id: ", err)
		return nil, err
	}

	playlistEntity, err := shared.RunDBTransaction(
		ctx,
		apiCfg.ResourceConfig,
		func(queries *database.Queries) (*database.Playlist, error) {
			playlistEntity, err := playlist.CreatePlaylist(
				ctx, queries,
				spotifyPlaylistDTOToCreatePlaylistParams(spotifyPlaylist),
			)
			if err != nil {
				return nil, err
			}

			userPlaylistEntity := database.CreateUserPlaylistParams{
				PlaylistID:             playlistEntity.ID,
				UserID:                 authUserID,
				IsSpotifySavedPlaylist: false,
			}

			_, err = userplaylist.CreateUserPlaylist(ctx, queries, userPlaylistEntity)
			if err != nil {
				return nil, err
			}

			return playlistEntity, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if _, err := enqueueDownloadPlaylistAudios(
		apiCfg,
		spotifyPlaylistID,
		spotifyAccessToken,
	); err != nil {
		return nil, err
	}

	return playlistEntity, nil
}

func downloadSpotifyUserSavedPlaylists(
	ctx context.Context,
	apiCfg *config.ApiConfig,
	authUserID uuid.UUID,
	spotifyAccessToken string,
) error {
	spotifyPlaylists, err := SpotifyGetUserSavedPlaylists(spotifyAccessToken)
	if err != nil {
		return err
	}

	if err := playlist.DeleteSpotifyUserSavedPlaylists(
		ctx, apiCfg.ResourceConfig,
		authUserID,
	); err != nil {
		return err
	}

	createPlaylistParams := []database.CreatePlaylistParams{}
	createUserPlaylistParams := []database.CreateUserPlaylistParams{}

	for _, playlist := range spotifyPlaylists.Items {
		playlistEntityCreateParams := spotifyPlaylistDTOToCreatePlaylistParams(&playlist)
		userPlaylistEntityCreateParams := database.CreateUserPlaylistParams{
			PlaylistID:             playlistEntityCreateParams.ID,
			UserID:                 authUserID,
			IsSpotifySavedPlaylist: true,
		}

		createPlaylistParams = append(createPlaylistParams, playlistEntityCreateParams)
		createUserPlaylistParams = append(createUserPlaylistParams, userPlaylistEntityCreateParams)
	}

	if err := shared.RunNoResultDBTransaction(
		ctx,
		apiCfg.ResourceConfig,
		func(queries *database.Queries) error {
			for _, params := range createPlaylistParams {
				_, err := playlist.CreatePlaylist(ctx, queries, params)
				if err != nil {
					return err
				}
			}

			for _, params := range createUserPlaylistParams {
				_, err := userplaylist.CreateUserPlaylist(ctx, queries, params)
				if err != nil {
					return err
				}
			}

			return nil
		},
	); err != nil {
		return err
	}

	for _, spotifyPlaylist := range spotifyPlaylists.Items {
		if _, err := enqueueDownloadPlaylistAudios(
			apiCfg,
			spotifyPlaylist.ID,
			spotifyAccessToken,
		); err != nil {
			return err
		}
	}

	return nil
}

func spotifyPlaylistDTOToCreatePlaylistParams(
	playlist *spotifyPlaylistDTO,
) database.CreatePlaylistParams {
	thumbnailUrl := ""
	if len(playlist.Images) > 0 {
		thumbnailUrl = playlist.Images[0].URL
	}

	return database.CreatePlaylistParams{
		ID:                uuid.New(),
		SpotifyID:         sql.NullString{String: playlist.ID, Valid: true},
		Name:              playlist.Name,
		ThumbnailUrl:      sql.NullString{String: thumbnailUrl, Valid: true},
		AudioImportStatus: database.ProcessStatusPENDING,
		TotalAudioCount:   int32(playlist.Tracks.Total),
		AudioCount:        0,
	}
}

func setPlaylistImportStatusToFailed(
	ctx context.Context,
	db *database.Queries,
	playlistID uuid.UUID,
) error {
	_, err := playlist.UpdatePlaylistByID(
		ctx, db,
		database.UpdatePlaylistByIDParams{
			PlaylistID:        playlistID,
			AudioImportStatus: database.NullProcessStatus{ProcessStatus: database.ProcessStatusFAILED, Valid: true},
		},
	)

	return err
}
