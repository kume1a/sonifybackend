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
		return err
	}

	dbPlaylists, err := playlist.GetUserPlaylistsBySpotifyIds(requestContext, apiCfg.DB, database.GetUserPlaylistsBySpotifyIdsParams{
		SpotifyIds: shared.Map(playlists.Items, func(playlist spotifyPlaylistDTO) string { return playlist.ID }),
		UserID:     authUserID,
	})
	if err != nil {
		return err
	}

	createAudioParams := []database.CreateAudioParams{}
	createPlaylistParams := []database.CreatePlaylistParams{}
	createPlaylistAudioParams := []database.CreatePlaylistAudioParams{}
	createUserPlaylistParams := []database.CreateUserPlaylistParams{}
	deletePlaylistAudioParams := []database.DeletePlaylistAudiosByIdsParams{}

	for _, playlist := range playlists.Items {
		playlistItems, err := GetSpotifyPlaylistItems(spotifyAccessToken, playlist.ID)
		if err != nil {
			return err
		}

		dbPlaylist := shared.FirstOrDefault(dbPlaylists, func(dbPlaylist *database.Playlist) bool {
			if !dbPlaylist.SpotifyID.Valid {
				return false
			}

			return playlist.ID == dbPlaylist.SpotifyID.String
		})

		toDownloadAudios, toAttachPlaylistItems, toDeletePlaylistItems, err := partitionSpotifyPlaylistItems(requestContext, apiCfg, &playlist, &playlistItems.Items, &dbPlaylists)
		if err != nil {
			return err
		}

		spotifyIdToAudioIdMap, err := getSpotifyIdToAudioIdMap(requestContext, apiCfg, createAudioParams, toAttachPlaylistItems)
		if err != nil {
			return err
		}

		trackMetas, err := downloadSpotifyPlaylistItems(toDownloadAudios)
		if err != nil {
			log.Println("Error downloading playlist items: ", err)
			return err
		}

		if dbPlaylist == nil {
			playlistEntityCreateParams := spotifyPlaylistDTOToCreatePlaylistParams(&playlist)
			userPlaylistEntityCreateParams := database.CreateUserPlaylistParams{
				PlaylistID: playlistEntityCreateParams.ID,
				UserID:     authUserID,
			}
			audioEntityCreateParams := playlistItemWithDownloadMetaToCreateAudioParams(trackMetas)
			playlistAudioCreateParams := shared.Map(
				toAttachPlaylistItems,
				func(spotifyPlaylistItem spotifyPlaylistItemDTO) database.CreatePlaylistAudioParams {
					return database.CreatePlaylistAudioParams{
						PlaylistID: playlistEntityCreateParams.ID,
						AudioID:    spotifyIdToAudioIdMap[spotifyPlaylistItem.Track.ID],
					}
				},
			)

			createPlaylistParams = append(createPlaylistParams, playlistEntityCreateParams)
			createUserPlaylistParams = append(createUserPlaylistParams, userPlaylistEntityCreateParams)
			createAudioParams = append(createAudioParams, audioEntityCreateParams...)
			createPlaylistAudioParams = append(createPlaylistAudioParams, playlistAudioCreateParams...)
		} else {
			audioEntityCreateParams := playlistItemWithDownloadMetaToCreateAudioParams(trackMetas)
			playlistAudioCreateParams := shared.Map(
				toAttachPlaylistItems,
				func(spotifyPlaylistItem spotifyPlaylistItemDTO) database.CreatePlaylistAudioParams {
					return database.CreatePlaylistAudioParams{
						PlaylistID: dbPlaylist.ID,
						AudioID:    spotifyIdToAudioIdMap[spotifyPlaylistItem.Track.ID],
					}
				},
			)

			createAudioParams = append(createAudioParams, audioEntityCreateParams...)
			createPlaylistAudioParams = append(createPlaylistAudioParams, playlistAudioCreateParams...)

			if len(toDeletePlaylistItems) > 0 {
				playlistAudioEntityDeleteParams := database.DeletePlaylistAudiosByIdsParams{
					PlaylistID: dbPlaylist.ID,
					AudioIds: shared.Map(
						toDeletePlaylistItems,
						func(item database.GetPlaylistAudioJoinsBySpotifyIdsRow) uuid.UUID { return item.AudioID },
					),
				}

				deletePlaylistAudioParams = append(deletePlaylistAudioParams, playlistAudioEntityDeleteParams)
			}
		}
	}

	shared.RunDbTransaction(
		requestContext,
		apiCfg,
		func(queries *database.Queries) (any, error) {
			for _, params := range deletePlaylistAudioParams {
				if err := playlist.DeletePlaylistAudiosByIds(requestContext, apiCfg.DB, params); err != nil {
					return nil, err
				}
			}

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

func partitionSpotifyPlaylistItems(
	requestContext context.Context,
	apiCfg *shared.ApiConfg,
	spotifyPlaylist *spotifyPlaylistDTO,
	playlistItems *[]spotifyPlaylistItemDTO,
	dbPlaylists *[]database.Playlist,
) (
	toDownloadAudios []spotifyPlaylistItemDTO,
	toAttachPlaylistItems []spotifyPlaylistItemDTO,
	toDeletePlaylistItems []database.GetPlaylistAudioJoinsBySpotifyIdsRow,
	err error,
) {
	playlistTrackSpotifyIds := shared.Map(*playlistItems, func(item spotifyPlaylistItemDTO) string { return item.Track.ID })

	dbAudioSpotifyIds, err := audio.GetAudioSpotifyIdsBySpotifyIds(requestContext, apiCfg.DB, playlistTrackSpotifyIds)
	if err != nil {
		return nil, nil, nil, err
	}

	toDownloadAudios = []spotifyPlaylistItemDTO{}
	for _, item := range *playlistItems {
		spotifyIdFromDB := shared.FirstOrDefault(
			dbAudioSpotifyIds,
			func(e *database.GetAudioSpotifyIdsBySpotifyIdsRow) bool { return e.SpotifyID.String == item.Track.ID },
		)
		if spotifyIdFromDB == nil {
			toDownloadAudios = append(toDownloadAudios, item)
		}
	}

	dbPlaylist := shared.FirstOrDefault(
		*dbPlaylists,
		func(dbPlaylist *database.Playlist) bool {
			return !dbPlaylist.SpotifyID.Valid && spotifyPlaylist.ID == dbPlaylist.SpotifyID.String
		},
	)

	if dbPlaylist == nil {
		return toDownloadAudios, *playlistItems, []database.GetPlaylistAudioJoinsBySpotifyIdsRow{}, nil
	}

	dbPlaylistAudios, err := playlist.GetPlaylistAudioJoinsBySpotifyIds(requestContext, apiCfg.DB, database.GetPlaylistAudioJoinsBySpotifyIdsParams{
		PlaylistID: dbPlaylist.ID,
		SpotifyIds: playlistTrackSpotifyIds,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	dbPlaylistAudiosMap := shared.SliceToMap(dbPlaylistAudios, func(audio *database.GetPlaylistAudioJoinsBySpotifyIdsRow) string { return audio.SpotifyID.String })
	playlistItemsMap := shared.SliceToMap(*playlistItems, func(item *spotifyPlaylistItemDTO) string { return item.Track.ID })

	toAttachPlaylistItems = []spotifyPlaylistItemDTO{}
	toDeletePlaylistItems = []database.GetPlaylistAudioJoinsBySpotifyIdsRow{}

	for _, item := range *playlistItems {
		_, ok := dbPlaylistAudiosMap[item.Track.ID]
		if !ok {
			toAttachPlaylistItems = append(toAttachPlaylistItems, item)
		}
	}

	for _, dbAudio := range dbPlaylistAudios {
		_, ok := playlistItemsMap[dbAudio.SpotifyID.String]
		if !ok {
			toDeletePlaylistItems = append(toDeletePlaylistItems, dbAudio)
		}
	}

	return toDownloadAudios, toAttachPlaylistItems, toDeletePlaylistItems, nil
}

func getSpotifyIdToAudioIdMap(
	requestContext context.Context,
	apiCfg *shared.ApiConfg,
	audioEntityCreateParams []database.CreateAudioParams,
	toAttachPlaylistItems []spotifyPlaylistItemDTO,
) (map[string]uuid.UUID, error) {
	toAttachSpotifyIds := shared.Map(toAttachPlaylistItems, func(item spotifyPlaylistItemDTO) string { return item.Track.ID })

	toAttachIds, err := audio.GetAudioSpotifyIdsBySpotifyIds(requestContext, apiCfg.DB, toAttachSpotifyIds)
	if err != nil {
		return nil, err
	}

	spotifyIdToAudioIds := make(map[string]uuid.UUID)
	for _, toAttachId := range toAttachIds {
		spotifyIdToAudioIds[toAttachId.SpotifyID.String] = toAttachId.ID
	}

	for _, audioEntityCreateParam := range audioEntityCreateParams {
		if _, ok := spotifyIdToAudioIds[audioEntityCreateParam.SpotifyID.String]; !ok {
			spotifyIdToAudioIds[audioEntityCreateParam.SpotifyID.String] = audioEntityCreateParam.ID
		}
	}

	return spotifyIdToAudioIds, nil
}

func downloadSpotifyPlaylistItems(playlistItems []spotifyPlaylistItemDTO) ([]*playlistItemWithDownloadMeta, error) {
	return shared.ExecuteParallel(
		playlistItems,
		func(input *spotifyPlaylistItemDTO) (*playlistItemWithDownloadMeta, error) {
			downloadMeta, err := GetSpotifyAudioDownloadMeta(input.Track.ID)
			if err != nil {
				log.Println("Error getting download meta: ", err)
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

func playlistItemWithDownloadMetaToCreateAudioParams(
	playlistItemWithDownloadMetas []*playlistItemWithDownloadMeta,
) []database.CreateAudioParams {
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
				DurationMs:   sql.NullInt32{Int32: int32(itemWithDownloadMeta.playlistItem.Track.DurationMS), Valid: true},
				Title:        sql.NullString{String: itemWithDownloadMeta.downloadMeta.Metadata.Title, Valid: true},
				SpotifyID:    sql.NullString{String: itemWithDownloadMeta.playlistItem.Track.ID, Valid: true},
				Path:         sql.NullString{String: itemWithDownloadMeta.downloadedAudioPath, Valid: true},
				SizeBytes:    sql.NullInt64{Int64: int64(itemWithDownloadMeta.downloadedAudioSize.Bytes), Valid: true},
				ThumbnailUrl: sql.NullString{String: thumbnailUrl, Valid: true},
			}
		},
	)
}
