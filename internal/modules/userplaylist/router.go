package userplaylist

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/userplaylist").Subrouter()

	r.HandleFunc("/myPlaylists", shared.AuthMW(handleGetAuthUserPlaylists(apiCfg))).Methods("GET")
	r.HandleFunc("/myPlaylistIds", shared.AuthMW(handleGetAuthUserPlaylistIDs(apiCfg))).Methods("GET")

	return r
}
