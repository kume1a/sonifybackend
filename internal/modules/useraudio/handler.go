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

		useraudios, err := BulkCreateUserAudios(r.Context(), apiCfg.ResourceConfig, params)
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

		err = DeleteUserAudio(
			r.Context(),
			apiCfg.DB,
			database.DeleteUserAudioParams{
				UserID:  authUser.UserID,
				AudioID: audioID,
			},
		)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		shared.ResNoContent(w)
	}
}
