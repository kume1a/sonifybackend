package audiolike

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/audiolike").Subrouter()

	r.HandleFunc("/myLikes", shared.AuthMW(handleGetAuthUserAudioLikes(apiCfg))).Methods("GET")

	r.HandleFunc("/like", shared.AuthMW(handleLikeAudio(apiCfg))).Methods("POST")
	r.HandleFunc("/unlike", shared.AuthMW(handleUnlikeAudio(apiCfg))).Methods("POST")

	return r
}
