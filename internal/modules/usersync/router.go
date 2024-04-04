package usersync

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/usersync").Subrouter()

	r.HandleFunc("/myUserSyncDatum", shared.AuthMW(handleGetUserSyncDatum(apiCfg))).Methods("GET")

	return r
}
