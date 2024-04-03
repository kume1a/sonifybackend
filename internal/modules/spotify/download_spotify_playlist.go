package spotify

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type playlistItemWithDownloadMeta struct {
	downloadMeta        *downloadSpotifyTrackMetaDTO
	playlistItem        *spotifyPlaylistItemDTO
	downloadedAudioPath string
	downloadedAudioSize *shared.FileSize
}

func downloadSpotifyPlaylist(
	apiCfg *shared.ApiConfg,
	requestContext context.Context,
	authUserID uuid.UUID,
	spotifyAccessToken string,
) error {
	playlists, err := GetSavedSpotifyPlaylists(spotifyAccessToken)
	if err != nil {
		return shared.HttpErrInternalServerError(shared.ErrFailedToGetSpotifyPlaylists)
	}

	createAudioParams := []database.CreateAudioParams{}
	createPlaylistParams := []database.CreatePlaylistParams{}
	createPlaylistAudioParams := []database.CreatePlaylistAudioParams{}
	createUserPlaylistParams := []database.CreateUserPlaylistParams{}

	for _, playlist := range playlists.Items {
		playlistItems, err := GetSpotifyPlaylistItems(spotifyAccessToken, playlist.ID)
		if err != nil {
			return shared.HttpErrInternalServerError(shared.ErrFailedToGetSpotifyPlaylist)
		}

		trackMetas, err := shared.ExecuteParallel(
			playlistItems.Items,
			func(input *spotifyPlaylistItemDTO) (*playlistItemWithDownloadMeta, error) {
				downloadMeta, err := GetSpotifyAudioDownloadMeta(input.Track.ID)
				if err != nil {
					return nil, err
				}

				downloadedAudioPath, err := shared.NewPublicFileLocation(shared.PublicFileLocationArgs{
					Extension: "mp3",
					Dir:       shared.DirSpotifyAudios,
				})
				if err != nil {
					return nil, err
				}

				err = shared.DownloadFile(downloadedAudioPath, downloadMeta.Link)
				if err != nil {
					return nil, err
				}

				audioFileSize, err := shared.GetFileSize(downloadedAudioPath)
				if err != nil {
					return nil, err
				}

				return &playlistItemWithDownloadMeta{
					downloadMeta:        downloadMeta,
					playlistItem:        input,
					downloadedAudioPath: downloadedAudioPath,
					downloadedAudioSize: audioFileSize,
				}, nil
			},
		)
		if err != nil {
			log.Println("Error downloading playlist items: ", err)
			return err
		}

		playlistEntityCreateParams := mapToCreatePlaylistEntityParams(&playlist)
		userPlaylistEntityCreateParams := database.CreateUserPlaylistParams{
			PlaylistID: playlistEntityCreateParams.ID,
			UserID:     authUserID,
		}
		audioEntityCreateParams := mapToCreateAudioEntityParams(trackMetas)
		playlistAudioCreateParams := shared.Map(
			audioEntityCreateParams,
			func(audioParams database.CreateAudioParams) database.CreatePlaylistAudioParams {
				return database.CreatePlaylistAudioParams{
					PlaylistID: playlistEntityCreateParams.ID,
					AudioID:    audioParams.ID,
				}
			},
		)

		createPlaylistParams = append(createPlaylistParams, playlistEntityCreateParams)
		createUserPlaylistParams = append(createUserPlaylistParams, userPlaylistEntityCreateParams)
		createAudioParams = append(createAudioParams, audioEntityCreateParams...)
		createPlaylistAudioParams = append(createPlaylistAudioParams, playlistAudioCreateParams...)
	}

	shared.RunDbTransaction(
		requestContext,
		apiCfg,
		func(queries *database.Queries) (any, error) {
			for _, params := range createPlaylistParams {
				_, err := playlist.CreatePlaylist(requestContext, apiCfg.DB, params)
				if err != nil {
					return nil, err
				}
			}

			for _, params := range createUserPlaylistParams {
				_, err := playlist.CreateUserPlaylist(requestContext, apiCfg.DB, params)
				if err != nil {
					return nil, err
				}
			}

			for _, params := range createAudioParams {
				_, err := audio.CreateAudio(requestContext, apiCfg.DB, params)
				if err != nil {
					return nil, err
				}
			}

			for _, params := range createPlaylistAudioParams {
				_, err := playlist.CreatePlaylistAudio(requestContext, apiCfg.DB, params)
				if err != nil {
					return nil, err
				}
			}

			return nil, nil
		},
	)

	return nil
}

func mapToCreatePlaylistEntityParams(playlist *spotifyPlaylistDTO) database.CreatePlaylistParams {
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

func mapToCreateAudioEntityParams(playlistItemWithDownloadMetas []*playlistItemWithDownloadMeta) []database.CreateAudioParams {
	return shared.Map(
		playlistItemWithDownloadMetas,
		func(itemWithDownloadMeta *playlistItemWithDownloadMeta) database.CreateAudioParams {
			artistName := ""
			if len(itemWithDownloadMeta.playlistItem.Track.Artists) > 0 {
				artistName = itemWithDownloadMeta.playlistItem.Track.Artists[0].Name
			}

			thumbnailUrl := ""
			if len(itemWithDownloadMeta.playlistItem.Track.Album.Images) > 0 {
				thumbnailUrl = itemWithDownloadMeta.playlistItem.Track.Album.Images[0].URL
			}

			return database.CreateAudioParams{
				ID:           uuid.New(),
				Author:       sql.NullString{String: artistName, Valid: true},
				Duration:     sql.NullInt32{Int32: int32(itemWithDownloadMeta.playlistItem.Track.DurationMS / 1000), Valid: true},
				Title:        sql.NullString{String: itemWithDownloadMeta.downloadMeta.Metadata.Title, Valid: true},
				SpotifyID:    sql.NullString{String: itemWithDownloadMeta.playlistItem.Track.ID, Valid: true},
				Path:         sql.NullString{String: itemWithDownloadMeta.downloadedAudioPath, Valid: true},
				SizeBytes:    sql.NullInt64{Int64: int64(itemWithDownloadMeta.downloadedAudioSize.Bytes), Valid: true},
				ThumbnailUrl: sql.NullString{String: thumbnailUrl, Valid: true},
			}
		},
	)
}
