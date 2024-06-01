package userplaylist

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/userplaylist").Subrouter()

	r.HandleFunc("/myUserPlaylists/full", shared.AuthMW(handleGetUserPlaylistsFullByAuthUserID(apiCfg))).Methods("GET")
	r.HandleFunc("/myUserPlaylists", shared.AuthMW(handleGetUserPlaylistsByAuthUserID(apiCfg))).Methods("GET")
	r.HandleFunc("/myPlaylistIds", shared.AuthMW(handleGetPlaylistIDsByAuthUser(apiCfg))).Methods("GET")
	r.HandleFunc("/myUserPlaylistIds", shared.AuthMW(handleGetUserPlaylistIDsByAuthUser(apiCfg))).Methods("GET")

	return r
}
