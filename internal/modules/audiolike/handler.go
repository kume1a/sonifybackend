package audiolike

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleLikeAudio(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*likeUnlikeAudioDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		newAudioLike, err := LikeAudio(r.Context(), apiCfg.DB, LikeUnlikeAudioParams{
			UserID:  authPayload.UserID,
			AudioID: body.AudioID,
		})
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		audioLikeDTO := sharedmodule.AudioLikeEntityToDTO(newAudioLike)

		shared.ResOK(w, audioLikeDTO)
	}
}

func handleUnlikeAudio(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*likeUnlikeAudioDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		err = UnlikeAudio(r.Context(), apiCfg.DB, LikeUnlikeAudioParams{
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

func handleGetAuthUserAudioLikes(apiCfg *config.ApiConfig) http.HandlerFunc {
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

		var audioLikes []database.AudioLike

		if len(body.IDs) == 0 {
			audioLikes, err = GetAudioLikesByUserID(r.Context(), apiCfg.DB, authPayload.UserID)

			if err != nil {
				shared.ResInternalServerErrorDef(w)
				return
			}
		} else {
			audioLikes, err = GetAudioLikesByUserIDAndAudioIDs(
				r.Context(),
				apiCfg.DB,
				database.GetAudioLikesByUserIDAndAudioIDsParams{
					UserID:   authPayload.UserID,
					AudioIds: body.IDs,
				},
			)

			if err != nil {
				shared.ResInternalServerErrorDef(w)
				return
			}
		}

		shared.ResOK(w, sharedmodule.AudioLikeEntityListToDTOList(audioLikes))
	}
}
