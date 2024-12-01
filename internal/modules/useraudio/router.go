package useraudio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/useraudio").Subrouter()

	r.HandleFunc("/createForAuthUser", shared.AuthMW(handleCreateUserAudiosForAuthUser(apiCfg))).Methods("POST")
	r.HandleFunc("/deleteForAuthUser", shared.AuthMW(handleDeleteUserAudioForAuthUser(apiCfg))).Methods("DELETE")

	r.HandleFunc("/myAudioIds", shared.AuthMW(handleGetAuthUserAudioIDs(apiCfg))).Methods("GET")
	r.HandleFunc("/myUserAudiosByIds", shared.AuthMW(handleGetAuthUserUserAudiosByAudioIDs(apiCfg))).Methods("GET")

	return r
}
