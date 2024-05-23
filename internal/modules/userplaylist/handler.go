package userplaylist

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/playlist"
	"github.com/kume1a/sonifybackend/internal/shared"
)

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

		dtos := shared.Map(playlists, playlist.PlaylistEntityToDto)

		shared.ResOK(w, dtos)
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
