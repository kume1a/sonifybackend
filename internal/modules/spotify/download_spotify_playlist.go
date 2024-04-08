package spotify

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	playlistusecase "github.com/kume1a/sonifybackend/internal/modules/playlist/usecase"
	spotifyusecase "github.com/kume1a/sonifybackend/internal/modules/spotify/usecase"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func downloadSpotifyPlaylist(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
	authUserID uuid.UUID,
	spotifyAccessToken string,
) error {
	playlists, err := GetSavedSpotifyPlaylists(spotifyAccessToken)
	if err != nil {
		return err
	}

	playlistusecase.DeleteSpotifyUserSavedPlaylists(ctx, apiCfg, authUserID)

	createPlaylistParams := []database.CreatePlaylistParams{}
	createPlaylistAudioParams := []database.CreatePlaylistAudioParams{}
	createUserPlaylistParams := []database.CreateUserPlaylistParams{}

	for _, playlist := range playlists.Items {
		playlistItems, err := GetSpotifyPlaylistItems(spotifyAccessToken, playlist.ID)
		if err != nil {
			return err
		}

		if err := spotifyusecase.DownloadWriteSpotifyAudios(
			ctx,
			apiCfg,
			shared.Map(
				playlistItems.Items,
				func(playlistItem spotifyPlaylistItemDTO) spotifyusecase.DownloadSpotifyAudioInput {
					artistName := ""
					if len(playlistItem.Track.Artists) > 0 {
						artistName = playlistItem.Track.Artists[0].Name
					}

					thumbnailURL := ""
					if len(playlistItem.Track.Album.Images) > 0 {
						thumbnailURL = playlistItem.Track.Album.Images[0].URL
					}

					return spotifyusecase.DownloadSpotifyAudioInput{
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

		dbPlaylistAudioSpotifyIds, err := audio.GetAudioSpotifyIdsBySpotifyIds(ctx, apiCfg.DB, playlistItemSpotifyIDs)
		if err != nil {
			return err
		}

		playlistEntityCreateParams := spotifyPlaylistDTOToCreatePlaylistParams(&playlist)
		userPlaylistEntityCreateParams := database.CreateUserPlaylistParams{
			PlaylistID: playlistEntityCreateParams.ID,
			UserID:     authUserID,
		}
		playlistAudioCreateParams := shared.Map(
			dbPlaylistAudioSpotifyIds,
			func(e database.GetAudioSpotifyIdsBySpotifyIdsRow) database.CreatePlaylistAudioParams {
				return database.CreatePlaylistAudioParams{
					PlaylistID: playlistEntityCreateParams.ID,
					AudioID:    e.ID,
				}
			},
		)

		createPlaylistParams = append(createPlaylistParams, playlistEntityCreateParams)
		createUserPlaylistParams = append(createUserPlaylistParams, userPlaylistEntityCreateParams)
		createPlaylistAudioParams = append(createPlaylistAudioParams, playlistAudioCreateParams...)
	}

	_, err = shared.RunDbTransaction(
		ctx,
		apiCfg,
		func(queries *database.Queries) (any, error) {
			for _, params := range createPlaylistParams {
				_, err := playlist.CreatePlaylist(ctx, queries, params)
				if err != nil {
					return nil, err
				}
			}

			for _, params := range createPlaylistAudioParams {
				_, err := playlist.CreatePlaylistAudio(ctx, queries, params)
				if err != nil {
					return nil, err
				}
			}

			for _, params := range createUserPlaylistParams {
				_, err := playlist.CreateUserPlaylist(ctx, queries, params)
				if err != nil {
					return nil, err
				}
			}

			return nil, nil
		},
	)

	return err
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
