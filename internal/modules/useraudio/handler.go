package useraudio

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleCreateUserAudiosForAuthUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*shared.AudioIDsDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		params := make([]database.CreateUserAudioParams, 0, len(body.AudioIDs))
		for _, audioID := range body.AudioIDs {
			params = append(params, database.CreateUserAudioParams{
				UserID:  authUser.UserID,
				AudioID: audioID,
			})
		}

		useraudios, err := BulkCreateUserAudiosTx(r.Context(), apiCfg.ResourceConfig, params)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		useraudioDTOs := sharedmodule.UserAudioEntitiesToDTOs(useraudios)

		shared.ResCreated(w, useraudioDTOs)
	}
}

func handleDeleteUserAudioForAuthUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*shared.AudioIDDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		audioID, err := uuid.Parse(body.AudioID)
		if err != nil {
			shared.ResBadRequest(w, shared.ErrInvalidUUID)
			return
		}

		if err := DeleteUserAudioTx(
			r.Context(),
			apiCfg.ResourceConfig,
			database.DeleteUserAudioParams{
				UserID:  authUser.UserID,
				AudioID: audioID,
			},
		); err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		shared.ResNoContent(w)
	}
}

func handleGetAuthUserAudioIDs(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		userAudioIds, err := GetUserAudioIDs(r.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, userAudioIds)
	}
}

func handleGetAuthUserUserAudiosByAudioIDs(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		// user body for big payload
		body, err := shared.GetRequestBody[*shared.AudioIDsDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		audios, err := GetUserAudiosByAudioIDs(
			r.Context(),
			apiCfg.DB,
			database.GetUserAudiosByAudioIdsParams{
				UserID:   authPayload.UserID,
				AudioIds: body.AudioIDs,
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		res := make([]*sharedmodule.UserAudioWithRelDTO, 0, len(audios))
		for _, audio := range audios {
			res = append(res, GetUserAudiosByAudioIdsRowToUserAudioWithRelDTO(audio))
		}

		shared.ResOK(w, res)
	}
}
