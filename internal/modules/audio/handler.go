package audio

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleDownloadYoutubeAudio(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[downloadYoutubeAudioDto](r)
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

		res := struct {
			*UserAudioDto
			Audio *AudioDto `json:"audio"`
		}{
			UserAudioDto: UserAudioEntityToDto(userAudio),
			Audio:        AudioEntityToDto(*audio),
		}

		shared.ResCreated(w, res)
	}
}
