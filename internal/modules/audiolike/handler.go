package audiolike

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleLikeAudio(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*likeAudioDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		newAudioLike, err := CreateAudioLike(r.Context(), apiCfg.DB, database.CreateAudioLikeParams{
			UserID:  authPayload.UserID,
			AudioID: body.AudioID,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		audioLikeDTO := AudioLikeEntityToDTO(newAudioLike)

		shared.ResOK(w, audioLikeDTO)
	}
}

func handleUnlikeAudio(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*unlikeAudioDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		err = DeleteAudioLike(r.Context(), apiCfg.DB, database.DeleteAudioLikeParams{
			UserID:  authPayload.UserID,
			AudioID: body.AudioID,
		})
		if shared.IsDBErrorNotFound(err) {
			shared.ResNotFound(w, shared.ErrAudioLikeNotFound)
			return
		} else if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, nil)
	}
}

func handleGetAuthUserAudioLikes(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		// user body for big payload
		body, err := shared.ValidateRequestBody[*getAudioLikesDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		audioLikes, err := GetAudioLikes(r.Context(), apiCfg.DB, database.GetAudioLikesParams{
			Ids:    body.IDs,
			UserID: authPayload.UserID,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, AudioLikeEntityListToDTOList(audioLikes))
	}
}
