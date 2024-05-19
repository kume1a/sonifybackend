package audio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/audio").Subrouter()

	r.HandleFunc("/myAudioIds", shared.WrapHandlerFunc(apiCfg, handleGetAuthUserAudioIds, shared.AuthMW)).Methods("GET")
	r.HandleFunc("/myUserAudiosByIds", shared.AuthMW(handleGetAuthUserUserAudiosByIDs(apiCfg))).Methods("GET")
	r.HandleFunc("/myLikes", shared.AuthMW(handleGetAuthUserAudioLikes(apiCfg))).Methods("GET")
	r.HandleFunc("/myLikesByIds", shared.AuthMW(handleGetAuthUserAudioLikesByIDs(apiCfg))).Methods("GET")

	r.HandleFunc("/downloadYoutubeAudio", shared.AuthMW(handleDownloadYoutubeAudio(apiCfg))).Methods("POST")
	r.HandleFunc("/uploadUserLocalMusic", shared.AuthMW(handleUploadUserLocalMusic(apiCfg))).Methods("POST")

	r.HandleFunc("/like", shared.AuthMW(handleLikeAudio(apiCfg))).Methods("POST")
	r.HandleFunc("/unlike", shared.AuthMW(handleUnlikeAudio(apiCfg))).Methods("POST")

	return r
}
