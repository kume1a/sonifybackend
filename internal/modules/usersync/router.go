package usersync

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/usersync").Subrouter()

	r.HandleFunc("/myUserSyncDatum", shared.AuthMW(handleGetUserSyncDatumByUserId(apiCfg))).Methods("GET")
	r.HandleFunc("/markUserAudioLastUpdatedAtAsNow", shared.AuthMW(handleMarkUserAudioLastUpdatedAtAsNow(apiCfg))).Methods("POST")

	return r
}
