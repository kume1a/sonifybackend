package playlist

import (
	"database/sql"
	"net/http"

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

		playlists, err := GetPlaylists(r.Context(), apiCfg.DB, database.GetPlaylistsParams{
			CreatedAt: query.LastCreatedAt,
			Limit:     query.Limit,
		})
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

		query, err := shared.ValidateRequestQuery[*getMyPlaylistsDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlists, err := GetUserPlaylists(r.Context(), apiCfg.DB, database.GetUserPlaylistsParams{
			UserID: authPayload.UserID,
			Ids:    query.IDs,
		})
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
		authPaylad, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		playlistIDDTO, httpErr := ValidateGetPlaylistByIDVars(r)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		playlist, audios, httpErr := GetPlaylistWithAudios(r.Context(), apiCfg.DB, playlistIDDTO.PlaylistID, authPaylad.UserID)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		dto := struct {
			playlistDTO
			Audios []*audio.AudioDTO `json:"audios"`
		}{
			playlistDTO: playlistEntityToDto(*playlist),
			Audios:      shared.Map(audios, audio.AudioWithAudioLikeToAudioDTO),
		}

		shared.ResOK(w, dto)
	}
}

func handleGetAuthUserPlaylistIDs(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		ids, err := GetUserPlaylistIDs(r.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, ids)
	}
}
