package youtube

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

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

func handleDownloadYoutubeAudioToUserLibrary(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*downloadYoutubeAudioToUserLibraryDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		userAudioWithAudio, err := DownloadYoutubeAudioAndSaveToUserLibrary(
			DownloadYoutubeAudioAndSaveToUserLibraryParams{
				ApiConfig:     apiCfg,
				Context:       r.Context(),
				UserID:        authPayload.UserID,
				YoutueVideoID: body.VideoID,
			},
		)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		res := sharedmodule.UserAudioWithRelDTO{
			UserAudioDTO: sharedmodule.UserAudioEntityToDTO(userAudioWithAudio.UserAudio),
			Audio:        sharedmodule.AudioEntityToDto(userAudioWithAudio.Audio),
		}

		shared.ResCreated(w, res)
	}
}

func handleDownloadYoutubeAudioPlaylist(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*downloadYoutubeAudioToPlaylistDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		playlistAudioWithAudio, err := DownloadYoutubeAudioAndSaveToPlaylist(
			DownloadYoutubeAudioAndSaveToPlaylistParams{
				ApiConfig:      apiCfg,
				Context:        r.Context(),
				UserID:         authPayload.UserID,
				YoutubeVideoID: body.VideoID,
				PlaylistID:     body.PlaylistID,
			},
		)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		res := sharedmodule.PlaylistAudioWithRelDTO{
			PlaylistAudioDTO: sharedmodule.PlaylistAudioEntityToDTO(playlistAudioWithAudio.PlaylistAudio),
			Audio:            sharedmodule.AudioEntityToDto(playlistAudioWithAudio.Audio),
		}

		shared.ResCreated(w, res)
	}
}
