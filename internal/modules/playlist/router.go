package playlist

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/playlist").Subrouter()

	r.HandleFunc("", shared.AuthMW(handleCreatePlaylist(apiCfg))).Methods("POST")
	r.HandleFunc("", shared.AuthMW(handleGetPlaylists(apiCfg))).Methods("GET")
	r.HandleFunc("/audio", shared.AuthMW(handleCreatePlaylistAudio(apiCfg))).Methods("POST")

	return r
}
