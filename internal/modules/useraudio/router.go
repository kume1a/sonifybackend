package useraudio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/useraudio").Subrouter()

	r.HandleFunc("/createForAuthUser", shared.AuthMW(handleCreateUserAudiosForAuthUser(apiCfg))).Methods("POST")

	return r
}
