package hiddenuseraudio

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleHideUserAudio(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*shared.AudioIDDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		newHiddenUserAudio, err := HideUserAudio(
			r.Context(),
			apiCfg.DB,
			HideUnhideAudioParams{
				UserID:  authPayload.UserID,
				AudioID: body.AudioID,
			},
		)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		resDTO := HiddenUserAudioEntityToDTO(newHiddenUserAudio)

		shared.ResOK(w, resDTO)
	}
}

func handleUnhideUserAudio(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*shared.AudioIDDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		err = UnhideAudio(r.Context(), apiCfg.DB, HideUnhideAudioParams{
			UserID:  authPayload.UserID,
			AudioID: body.AudioID,
		})
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		shared.ResOK(w, nil)
	}
}

func handleGetHiddenUserAudiosByAuthUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		// user body for big payload
		body, err := shared.ValidateRequestBody[*shared.OptionalIDsDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		var hiddenUserAudios []database.HiddenUserAudio

		if len(body.IDs) == 0 {
			hiddenUserAudios, err = GetHiddenUserAudiosByUserID(
				r.Context(),
				apiCfg.DB,
				authPayload.UserID,
			)

			if err != nil {
				shared.ResInternalServerErrorDef(w)
				return
			}
		} else {
			hiddenUserAudios, err = GetHiddenUserAudiosByUserIDAndAudioIDs(
				r.Context(),
				apiCfg.DB,
				database.GetHiddenUserAudiosByUserIDAndAudioIDsParams{
					UserID:   authPayload.UserID,
					AudioIds: body.IDs,
				},
			)

			if err != nil {
				shared.ResInternalServerErrorDef(w)
				return
			}
		}

		shared.ResOK(w, HiddenUserAudioEntityListToDTOList(hiddenUserAudios))
	}
}
