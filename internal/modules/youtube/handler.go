package youtube

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetYoutubeMusicUrl(w http.ResponseWriter, r *http.Request) {
	query, err := shared.ValidateRequestQuery[*getYoutubeMusicUrlDto](r)
	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	url, err := GetYoutubeAudioUrl(query.VideoID[0])
	if err != nil {
		shared.ResInternalServerErrorDef(w)
		return
	}

	dto := shared.UrlDTO{Url: url}
	shared.ResOK(w, dto)
}

func handleGetYoutubeSearchSuggestions(w http.ResponseWriter, r *http.Request) {
	query, err := shared.ValidateRequestQuery[*shared.KeywordDTO](r)
	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	if len(query.Keyword) != 1 {
		shared.ResBadRequest(w, shared.ErrInvalidQueryParams)
		return
	}

	res, err := GetYoutubeSearchSuggestions(query.Keyword[0])
	if err != nil {
		shared.ResInternalServerErrorDef(w)
		return
	}

	shared.ResOK(w, res)
}

func handleDownloadYoutubeAudio(apiCfg *config.ApiConfig) http.HandlerFunc {
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

		userAudioWithAudio, httpErr := DownloadYoutubeAudioAndSave(DownloadYoutubeAudioParams{
			ApiConfig: apiCfg,
			Context:   r.Context(),
			UserID:    authPayload.UserID,
			VideoID:   body.VideoID,
		})
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		res := sharedmodule.UserAudioWithRelDTO{
			UserAudioDTO: sharedmodule.UserAudioEntityToDTO(userAudioWithAudio.UserAudio),
			Audio:        sharedmodule.AudioEntityToDto(userAudioWithAudio.Audio),
		}

		shared.ResCreated(w, res)
	}
}
