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
	isEmpty             bool
	playlistItem        spotifyPlaylistItemDTO
	downloadedAudioPath string
	downloadedAudioSize *shared.FileSize
}

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

	dbPlaylists, err := playlist.GetUserPlaylistsBySpotifyIds(ctx, apiCfg.DB, database.GetUserPlaylistsBySpotifyIdsParams{
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

		dbPlaylist := shared.FirstOrDefaultP(dbPlaylists, func(dbPlaylist *database.Playlist) bool {
			if !dbPlaylist.SpotifyID.Valid {
				return false
			}

			return playlist.ID == dbPlaylist.SpotifyID.String
		})

		toDownloadAudios, toAttachPlaylistItems, toDeletePlaylistItems, err := partitionSpotifyPlaylistItems(ctx, apiCfg, &playlist, &playlistItems.Items, &dbPlaylists)
		if err != nil {
			return err
		}

		spotifyIdToAudioIdMap, err := getSpotifyIdToAudioIdMap(ctx, apiCfg, createAudioParams, toAttachPlaylistItems)
		if err != nil {
			return err
		}

		trackMetas, err := downloadSpotifyPlaylistItems(toDownloadAudios)
		if err != nil {
			log.Println("Error downloading playlist items: ", err)
			return err
		}

		// remove non-downloaded audios from toAttachPlaylistItems
		toAttachPlaylistItems = shared.Where(toAttachPlaylistItems, func(item spotifyPlaylistItemDTO) bool {
			return shared.ContainsWhere(
				trackMetas,
				func(meta playlistItemWithDownloadMeta) bool {
					return meta.playlistItem.Track.ID == item.Track.ID
				},
			)
		})

		audioEntityCreateParams := playlistItemWithDownloadMetaToCreateAudioParams(trackMetas)
		createAudioParams = append(createAudioParams, audioEntityCreateParams...)

		if dbPlaylist == nil {
			playlistEntityCreateParams := spotifyPlaylistDTOToCreatePlaylistParams(&playlist)
			userPlaylistEntityCreateParams := database.CreateUserPlaylistParams{
				PlaylistID: playlistEntityCreateParams.ID,
				UserID:     authUserID,
			}
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
			createPlaylistAudioParams = append(createPlaylistAudioParams, playlistAudioCreateParams...)
		} else {
			playlistAudioCreateParams := shared.Map(
				toAttachPlaylistItems,
				func(spotifyPlaylistItem spotifyPlaylistItemDTO) database.CreatePlaylistAudioParams {
					return database.CreatePlaylistAudioParams{
						PlaylistID: dbPlaylist.ID,
						AudioID:    spotifyIdToAudioIdMap[spotifyPlaylistItem.Track.ID],
					}
				},
			)

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

	_, err = shared.RunDbTransaction(
		ctx,
		apiCfg,
		func(queries *database.Queries) (any, error) {
			for _, params := range createAudioParams {
				_, err := audio.CreateAudio(ctx, queries, params)
				if err != nil {
					return nil, err
				}
			}

			return nil, nil
		},
	)

	if err != nil {
		return err
	}

	_, err = shared.RunDbTransaction(
		ctx,
		apiCfg,
		func(queries *database.Queries) (any, error) {
			for _, params := range deletePlaylistAudioParams {
				if err := playlist.DeletePlaylistAudiosByIds(ctx, queries, params); err != nil {
					return nil, err
				}
			}

			for _, params := range createPlaylistParams {
				_, err := playlist.CreatePlaylist(ctx, queries, params)
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

			for _, params := range createPlaylistAudioParams {
				_, err := playlist.CreatePlaylistAudio(ctx, queries, params)
				if err != nil {
					return nil, err
				}
			}

			return nil, nil
		},
	)

	return err
}

func partitionSpotifyPlaylistItems(
	ctx context.Context,
	apiCfg *shared.ApiConfig,
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

	dbAudioSpotifyIds, err := audio.GetAudioSpotifyIdsBySpotifyIds(ctx, apiCfg.DB, playlistTrackSpotifyIds)
	if err != nil {
		return nil, nil, nil, err
	}

	toDownloadAudios = []spotifyPlaylistItemDTO{}
	for _, item := range *playlistItems {
		spotifyIdExistsInDB := shared.ContainsWhereP(
			dbAudioSpotifyIds,
			func(e *database.GetAudioSpotifyIdsBySpotifyIdsRow) bool { return e.SpotifyID.String == item.Track.ID },
		)
		if !spotifyIdExistsInDB {
			toDownloadAudios = append(toDownloadAudios, item)
		}
	}

	dbPlaylist := shared.FirstOrDefaultP(
		*dbPlaylists,
		func(dbPlaylist *database.Playlist) bool {
			return !dbPlaylist.SpotifyID.Valid && spotifyPlaylist.ID == dbPlaylist.SpotifyID.String
		},
	)

	if dbPlaylist == nil {
		return toDownloadAudios, *playlistItems, []database.GetPlaylistAudioJoinsBySpotifyIdsRow{}, nil
	}

	dbPlaylistAudios, err := playlist.GetPlaylistAudioJoinsBySpotifyIds(ctx, apiCfg.DB, database.GetPlaylistAudioJoinsBySpotifyIdsParams{
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
	ctx context.Context,
	apiCfg *shared.ApiConfig,
	audioEntityCreateParams []database.CreateAudioParams,
	toAttachPlaylistItems []spotifyPlaylistItemDTO,
) (map[string]uuid.UUID, error) {
	toAttachSpotifyIds := shared.Map(toAttachPlaylistItems, func(item spotifyPlaylistItemDTO) string { return item.Track.ID })

	toAttachIds, err := audio.GetAudioSpotifyIdsBySpotifyIds(ctx, apiCfg.DB, toAttachSpotifyIds)
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

func downloadSpotifyPlaylistItems(playlistItems []spotifyPlaylistItemDTO) ([]playlistItemWithDownloadMeta, error) {
	itemsWithDownloadMeta, err := shared.ExecuteParallel(
		8,
		playlistItems,
		func(playlistItem spotifyPlaylistItemDTO) (playlistItemWithDownloadMeta, error) {
			downloadMeta, err := GetSpotifyAudioDownloadMeta(playlistItem.Track.ID)
			if err != nil {
				log.Println("Error getting download meta: ", err)
				return playlistItemWithDownloadMeta{isEmpty: true}, err
			}

			if !downloadMeta.Success {
				return playlistItemWithDownloadMeta{isEmpty: true}, nil
			}

			downloadedAudioPath, err := shared.NewPublicFileLocation(shared.PublicFileLocationArgs{
				Extension: "mp3",
				Dir:       shared.DirSpotifyAudios,
			})
			if err != nil {
				return playlistItemWithDownloadMeta{isEmpty: true}, err
			}

			err = shared.DownloadFile(downloadedAudioPath, downloadMeta.Link)
			if err != nil {
				log.Println("Error downloading file: ", err, " from: ", downloadMeta.Link, " meta = ", downloadMeta)
				return playlistItemWithDownloadMeta{isEmpty: true}, err
			}

			audioFileSize, err := shared.GetFileSize(downloadedAudioPath)
			if err != nil {
				return playlistItemWithDownloadMeta{isEmpty: true}, err
			}

			return playlistItemWithDownloadMeta{
				isEmpty:             false,
				playlistItem:        playlistItem,
				downloadedAudioPath: downloadedAudioPath,
				downloadedAudioSize: audioFileSize,
			}, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return shared.Where(itemsWithDownloadMeta, func(item playlistItemWithDownloadMeta) bool { return !item.isEmpty }), nil
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
	playlistItemWithDownloadMetas []playlistItemWithDownloadMeta,
) []database.CreateAudioParams {
	return shared.Map(
		playlistItemWithDownloadMetas,
		func(itemWithDownloadMeta playlistItemWithDownloadMeta) database.CreateAudioParams {
			artistName := ""
			if len(itemWithDownloadMeta.playlistItem.Track.Artists) > 0 {
				artistName = itemWithDownloadMeta.playlistItem.Track.Artists[0].Name
			}

			thumbnailUrl := ""
			if len(itemWithDownloadMeta.playlistItem.Track.Album.Images) > 0 {
				thumbnailUrl = itemWithDownloadMeta.playlistItem.Track.Album.Images[0].URL
			}

			log.Println("insert track name: ", itemWithDownloadMeta.playlistItem.Track.Name, ", track id: ", itemWithDownloadMeta.playlistItem.Track.ID)
			a := database.CreateAudioParams{
				ID:           uuid.New(),
				Author:       sql.NullString{String: artistName, Valid: true},
				DurationMs:   sql.NullInt32{Int32: int32(itemWithDownloadMeta.playlistItem.Track.DurationMS), Valid: true},
				Title:        sql.NullString{String: itemWithDownloadMeta.playlistItem.Track.Name, Valid: true},
				SpotifyID:    sql.NullString{String: itemWithDownloadMeta.playlistItem.Track.ID, Valid: true},
				Path:         sql.NullString{String: itemWithDownloadMeta.downloadedAudioPath, Valid: true},
				SizeBytes:    sql.NullInt64{Int64: int64(itemWithDownloadMeta.downloadedAudioSize.Bytes), Valid: true},
				ThumbnailUrl: sql.NullString{String: thumbnailUrl, Valid: true},
			}

			return a
		},
	)
}
