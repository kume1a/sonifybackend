package audio

import (
	"log"
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

		body, err := shared.ValidateRequestBody[downloadYoutubeAudioDTO](r)
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

func handleImportUserLocalMusic(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		form, httpErr := ValidateImportUserLocalMusicDTO(w, r)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		log.Println("authPayload = ", authPayload)
		log.Println("form.Title = ", form.Title, "form.Author = ", form.Author, "form.LocalId = ", form.LocalId, "form.AudioPath = ", form.AudioPath, "form.ThumbnailPath = ", form.ThumbnailPath)
	}
}
