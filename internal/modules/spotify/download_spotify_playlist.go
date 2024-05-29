package spotify

import (
	"context"
	"database/sql"

	"github.com/gocraft/work"
	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/modules/userplaylist"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func DownloadSpotifyPlaylistAudios(
	ctx context.Context,
	resouceConfig *config.ResourceConfig,
	playlistSpotifyID string,
	spotifyAccessToken string,
) error {
	createPlaylistAudioParams := []database.CreatePlaylistAudioParams{}

	playlistItems, err := GetSpotifyPlaylistItems(spotifyAccessToken, playlistSpotifyID)
	if err != nil {
		return err
	}

	if err := DownloadWriteSpotifyAudios(
		ctx,
		resouceConfig,
		shared.Map(
			playlistItems.Items,
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
	); err != nil {
		return err
	}

	playlistItemSpotifyIDs := shared.Map(
		playlistItems.Items,
		func(playlistItem spotifyPlaylistItemDTO) string { return playlistItem.Track.ID },
	)

	playlistAudioIDs, err := audio.GetAudioIdsBySpotifyIds(ctx, resouceConfig.DB, playlistItemSpotifyIDs)
	if err != nil {
		return err
	}

	playlistID, err := playlist.GetPlaylistIDBySpotifyID(ctx, resouceConfig.DB, playlistSpotifyID)
	if err != nil {
		return err
	}

	for _, playlistAudioID := range playlistAudioIDs {
		createPlaylistAudioParams = append(
			createPlaylistAudioParams,
			database.CreatePlaylistAudioParams{
				PlaylistID: playlistID,
				AudioID:    playlistAudioID,
			},
		)
	}

	if _, err := playlistaudio.BulkCreatePlaylistAudios(
		ctx,
		resouceConfig,
		createPlaylistAudioParams,
	); err != nil {
		return err
	}

	return nil
}

func downloadSpotifyPlaylist(
	ctx context.Context,
	apiCfg *config.ApiConfig,
	authUserID uuid.UUID,
	spotifyAccessToken string,
) error {
	spotifyPlaylists, err := GetSavedSpotifyPlaylists(spotifyAccessToken)
	if err != nil {
		return err
	}

	playlist.DeleteSpotifyUserSavedPlaylists(ctx, apiCfg.ResourceConfig, authUserID)

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

	if _, err := shared.RunDBTransaction(
		ctx,
		apiCfg.ResourceConfig,
		func(queries *database.Queries) (any, error) {
			for _, params := range createPlaylistParams {
				_, err := playlist.CreatePlaylist(ctx, queries, params)
				if err != nil {
					return nil, err
				}
			}

			for _, params := range createUserPlaylistParams {
				_, err := userplaylist.CreateUserPlaylist(ctx, queries, params)
				if err != nil {
					return nil, err
				}
			}

			return nil, nil
		},
	); err != nil {
		return err
	}

	for _, spotifyPlaylist := range spotifyPlaylists.Items {
		if _, err := apiCfg.WorkEnqueuer.Enqueue(
			shared.BackgroundJobDownloadPlaylistAudios,
			work.Q{
				"playlistSpotifyID":  spotifyPlaylist.ID,
				"spotifyAccessToken": spotifyAccessToken,
			},
		); err != nil {
			return err
		}
	}

	return nil
}

func spotifyPlaylistDTOToCreatePlaylistParams(playlist *spotifyPlaylistDTO) database.CreatePlaylistParams {
	thumbnailUrl := ""
	if len(playlist.Images) > 0 {
		thumbnailUrl = playlist.Images[0].URL
	}

	return database.CreatePlaylistParams{
		ID:           uuid.New(),
		SpotifyID:    sql.NullString{String: playlist.ID, Valid: true},
		Name:         playlist.Name,
		ThumbnailUrl: sql.NullString{String: thumbnailUrl, Valid: true},
	}
}
