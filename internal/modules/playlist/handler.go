package playlist

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleCreatePlaylist(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ValidateCreatePlaylistDto(w, r)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		playlist, err := CreatePlaylist(r.Context(), apiCfg.DB, database.CreatePlaylistParams{
			Name:              body.Name,
			ThumbnailPath:     sql.NullString{String: body.ThumbnailPath, Valid: true},
			AudioImportStatus: database.ProcessStatusCOMPLETED,
			AudioCount:        0,
			TotalAudioCount:   0,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResCreated(
			w,
			sharedmodule.PlaylistEntityToDto(*playlist),
		)
	}
}

func handleGetPlaylists(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := shared.ValidateRequestQuery[*shared.LastCreatedAtPageParamsDto](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlists, err := GetPlaylists(r.Context(), apiCfg.DB, database.GetPlaylistsParams{
			CreatedAt: query.LastCreatedAt,
			Limit:     query.Limit,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dtos := shared.Map(playlists, sharedmodule.PlaylistEntityToDto)

		shared.ResOK(w, dtos)
	}
}

func handleGetPlaylistFull(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPaylad, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		playlistIDDTO, err := ValidateGetPlaylistByIDVars(r)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		ctx := r.Context()

		playlist, err := GetPlaylistByID(ctx, apiCfg.DB, playlistIDDTO.PlaylistID)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		playlistAudios, err := playlistaudio.GetPlaylistAudios(
			ctx, apiCfg.DB,
			database.GetPlaylistAudiosParams{
				PlaylistIds: []uuid.UUID{playlistIDDTO.PlaylistID},
				UserID:      authPaylad.UserID,
				Ids:         nil,
			},
		)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		dto := PlaylistFullDTO{
			PlaylistDTO:    sharedmodule.PlaylistEntityToDto(*playlist),
			PlaylistAudios: shared.Map(playlistAudios, playlistaudio.GetPlaylistAudioRowToDTO),
		}

		shared.ResOK(w, dto)
	}
}
