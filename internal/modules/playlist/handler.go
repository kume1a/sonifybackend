package playlist

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/audio"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleCreatePlaylist(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, httpErr := ValidateCreatePlaylistDto(w, r)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
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

func handleGetPlaylists(apiCfg *shared.ApiConfig) http.HandlerFunc {
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

func handleCreatePlaylistAudio(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*createPlaylistAudioDTO](r)
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

func handleGetAuthUserPlaylists(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		playlists, err := GetUserPlaylists(r.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dtos := shared.Map(playlists, playlistEntityToDto)

		shared.ResOK(w, dtos)
	}
}

func handleGetPlaylistWithAudios(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		playlistID, ok := vars["playlistID"]
		if !ok {
			shared.ResBadRequest(w, "playlistID is required")
			return
		}

		playlistIDUUID, err := uuid.Parse(playlistID)
		if err != nil {
			shared.ResBadRequest(w, "playlistId is not a valid UUID")
			return
		}

		playlist, audios, httpErr := GetPlaylistWithAudios(r.Context(), apiCfg.DB, playlistIDUUID)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		dto := struct {
			playlistDTO
			Audios []*audio.AudioDTO `json:"audios"`
		}{
			playlistDTO: playlistEntityToDto(*playlist),
			Audios:      shared.Map(audios, audio.AudioEntityToDto),
		}

		shared.ResOK(w, dto)
	}
}
