package userplaylist

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetUserPlaylistsFullByAuthUserID(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		query, err := shared.ValidateRequestQuery[*PlaylistIDsDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlists, err := GetUserPlaylistsFull(
			r.Context(), apiCfg.DB,
			database.GetFullUserPlaylistsParams{
				UserID:      authPayload.UserID,
				PlaylistIds: query.PlaylistIDs,
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dtos := shared.Map(playlists, MapUserPlaylistFullEntityToDTO)

		shared.ResOK(w, dtos)
	}
}

func handleGetUserPlaylistsByAuthUserID(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		// using body for big payload
		body, err := shared.ValidateRequestBody[*shared.OptionalIDsDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlists, err := GetUserPlaylistsByUserID(
			r.Context(), apiCfg.DB,
			database.GetUserPlaylistsParams{
				UserID: authPayload.UserID,
				Ids:    body.IDs,
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dtos := shared.Map(playlists, MapUserPlaylistEntityToDTO)

		shared.ResOK(w, dtos)
	}
}

func handleGetPlaylistIDsByAuthUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		ids, err := GetPlaylistIDsByUserID(r.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, ids)
	}
}

func handleGetUserPlaylistIDsByAuthUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		ids, err := GetUserPlaylistIDsByUserID(r.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, ids)
	}
}
