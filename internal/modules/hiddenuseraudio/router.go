package hiddenuseraudio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/hiddenuseraudio").Subrouter()

	r.HandleFunc("/myhiddenuseraudios", shared.AuthMW(handleGetHiddenUserAudiosByAuthUser(apiCfg))).Methods("GET")

	r.HandleFunc("/hide", shared.AuthMW(handleHideUserAudio(apiCfg))).Methods("POST")
	r.HandleFunc("/unhide", shared.AuthMW(handleUnhideUserAudio(apiCfg))).Methods("POST")

	return r
}
