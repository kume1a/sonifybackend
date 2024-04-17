package audio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/audio").Subrouter()

	r.HandleFunc("/myAudios", shared.AuthMW(handleAuthGetUserAudios(apiCfg))).Methods("GET")
	r.HandleFunc("/myAudioIds", shared.AuthMW(handleGetAuthUserAudioIds(apiCfg))).Methods("GET")
	r.HandleFunc("/myUserAudios", shared.AuthMW(handleGetAuthUserUserAudios(apiCfg))).Methods("GET")

	r.HandleFunc("/downloadYoutubeAudio", shared.AuthMW(handleDownloadYoutubeAudio(apiCfg))).Methods("POST")
	r.HandleFunc("/uploadUserLocalMusic", shared.AuthMW(handleUploadUserLocalMusic(apiCfg))).Methods("POST")

	return r
}
