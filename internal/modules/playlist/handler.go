package playlist

import (
	"database/sql"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleCreatePlaylist(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, httpErr := ValidateCreatePlaylistDto(w, r)
		if httpErr != nil {
			shared.ResHttpError(w, *httpErr)
			return
		}

		playlist, err := CreatePlaylist(r.Context(), apiCfg.DB, database.CreatePlaylistParams{
			Name:          body.Name,
			ThumbnailPath: sql.NullString{String: body.ThumbnailPath, Valid: true},
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := playlistEntityToDto(*playlist)

		shared.ResCreated(w, dto)
	}
}

func handleGetPlaylists(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := shared.ValidateRequestQuery[*shared.LastCreatedAtPageParamsDto](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlists, err := GetPlaylists(r.Context(), apiCfg.DB, query.LastCreatedAt, query.Limit)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dtos := shared.Map(playlists, playlistEntityToDto)

		shared.ResOK(w, dtos)
	}
}

func handleCreatePlaylistAudio(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*createPlaylistAudioDto](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlistAudio, err := CreatePlaylistAudio(r.Context(), apiCfg.DB, database.CreatePlaylistAudioParams{
			PlaylistID: body.PlaylistID,
			AudioID:    body.AudioID,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := playlistAudioEntityToDto(playlistAudio)

		shared.ResCreated(w, dto)
	}
}

func handleGetAuthUserPlaylists(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		playlists, err := GetUserPlaylists(r.Context(), apiCfg.DB, authPayload.UserId)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dtos := shared.Map(playlists, playlistEntityToDto)

		shared.ResOK(w, dtos)
	}
}
