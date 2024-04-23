package audio

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleDownloadYoutubeAudio(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*downloadYoutubeAudioDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		userAudio, audio, httpErr := DownloadYoutubeAudio(DownloadYoutubeAudioParams{
			ApiConfig: apiCfg,
			Context:   r.Context(),
			UserID:    authPayload.UserID,
			VideoID:   body.VideoID,
		})
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		res := UserAudioWithRelDTO{
			UserAudioDTO: UserAudioEntityToDto(userAudio),
			Audio:        AudioEntityToDto(*audio),
		}

		shared.ResCreated(w, res)
	}
}

func handleUploadUserLocalMusic(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		form, httpErr := ValidateUploadUserLocalMusicDTO(w, r)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		audioExists, err := DoesAudioExistByLocalId(r.Context(), apiCfg.DB, authPayload.UserID, form.LocalID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		if audioExists {
			shared.DeleteFiles([]string{form.AudioPath, form.ThumbnailPath})

			shared.ResConflict(w, shared.ErrAudioAlreadyExists)
			return
		}

		userAudioWithAudio, httpErr := WriteUserImportedLocalMusic(WriteUserImportedLocalMusicParams{
			ApiConfig:          apiCfg,
			Context:            r.Context(),
			UserID:             authPayload.UserID,
			AudioTitle:         form.Title,
			AudioAuthor:        form.Author,
			AudioPath:          form.AudioPath,
			AudioThumbnailPath: form.ThumbnailPath,
			AudioDurationMs:    form.DurationMs,
			AudioLocalId:       form.LocalID,
		})

		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		res := UserAudioWithRelDTO{
			UserAudioDTO: UserAudioEntityToDto(userAudioWithAudio.UserAudio),
			Audio:        AudioEntityToDto(*userAudioWithAudio.Audio),
		}

		shared.ResCreated(w, res)
	}
}

func handleAuthGetUserAudios(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		userAudios, err := GetUserAudios(r.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		res := make([]*AudioDTO, 0, len(userAudios))
		for _, userAudio := range userAudios {
			res = append(res, AudioEntityToDto(userAudio))
		}

		shared.ResOK(w, res)
	}
}

func handleGetAuthUserAudioIds(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		userAudioIds, err := GetUserAudioIds(r.Context(), apiCfg.DB, authPayload.UserID)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, userAudioIds)
	}
}

func handleGetAuthUserUserAudios(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		// user body for big payload
		body, err := shared.GetRequestBody[getAudiosByIdsDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		if len(body.AudioIDs) == 0 {
			shared.ResOK(w, []UserAudioWithRelDTO{})
			return
		}

		audios, err := GetUserAudiosByAudioIds(r.Context(), apiCfg.DB, database.GetUserAudiosByAudioIdsParams{
			UserID:   authPayload.UserID,
			AudioIds: body.AudioIDs,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		res := make([]*UserAudioWithRelDTO, 0, len(audios))
		for _, audio := range audios {
			res = append(res, GetUserAudiosByAudioIdsRowToUserAudioWithRelDTO(audio))
		}

		shared.ResOK(w, res)
	}
}

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

		_, err = CreateAudioLike(r.Context(), apiCfg.DB, database.CreateAudioLikeParams{
			UserID:  authPayload.UserID,
			AudioID: body.AudioID,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, nil)
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

func handleSyncAudioLikes(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*syncAudioLikesDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		err = SyncAudioLikes(SyncAudioLikesParams{
			Context:  r.Context(),
			ApiCfg:   apiCfg,
			UserID:   authPayload.UserID,
			AudioIDs: body.AudioIDs,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResNoContent(w)
	}
}
