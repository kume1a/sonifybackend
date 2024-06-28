package playlist

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/playlists").Subrouter()

	r.HandleFunc("/{playlistID}", shared.AuthMW(handleGetPlaylistFull(apiCfg))).Methods("GET")

	return r
}
