package userplaylist

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleCreateUserPlaylist(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		body, err := shared.ValidateRequestBody[*createUserPlaylistDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlist, err := CreatePlaylistAndUserPlaylist(
			r.Context(),
			apiCfg.ResourceConfig,
			CreatePlaylistAndUserPlaylistParams{
				Name:   body.Name,
				UserID: authPayload.UserID,
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, sharedmodule.UserPlaylistWithRelToDTO(playlist))
	}
}

func handleUpdateUserPlaylist(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		userPlaylistID, err := shared.GetURLParamUUID(r, "userPlaylistID")
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		dto, err := shared.ValidateRequestBody[*updateUserPlaylistDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlist, err := UpdateUserPlaylist(
			r.Context(), apiCfg.DB,
			UpdateUserPlaylistParams{
				UserID:         authPayload.UserID,
				UserPlaylistID: userPlaylistID,
				Name:           dto.Name,
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, sharedmodule.UserPlaylistWithRelToDTO(playlist))
	}
}

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

		dtos := shared.Map(
			playlists,
			func(e database.GetFullUserPlaylistsRow) *sharedmodule.UserPlaylistDTO {
				return sharedmodule.MapUserPlaylistFullEntityToDTO(&e)
			},
		)

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

		dtos := shared.Map(
			playlists,
			func(e database.UserPlaylist) *sharedmodule.UserPlaylistDTO {
				return sharedmodule.UserPlaylistEntityToDTO(&e)
			},
		)

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

		ids, err := sharedmodule.GetPlaylistIDsByUserID(
			r.Context(),
			apiCfg.DB,
			authPayload.UserID,
		)
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

		ids, err := GetUserPlaylistIDsByUserID(
			r.Context(),
			apiCfg.DB,
			authPayload.UserID,
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, ids)
	}
}
