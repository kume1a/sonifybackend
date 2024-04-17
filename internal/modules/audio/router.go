package audio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/audio").Subrouter()

	r.HandleFunc("/myAudios", shared.AuthMW(handleGetUserAudios(apiCfg))).Methods("GET")
	r.HandleFunc("/myAudioIds", shared.AuthMW(handleGetUserAudioIds(apiCfg))).Methods("GET")
	r.HandleFunc("/all", shared.AuthMW(handleGetAudiosByIds(apiCfg))).Methods("GET")

	r.HandleFunc("/downloadYoutubeAudio", shared.AuthMW(handleDownloadYoutubeAudio(apiCfg))).Methods("POST")
	r.HandleFunc("/uploadUserLocalMusic", shared.AuthMW(handleUploadUserLocalMusic(apiCfg))).Methods("POST")

	return r
}
