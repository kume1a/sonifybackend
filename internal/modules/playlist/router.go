package playlist

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/playlists").Subrouter()

	r.HandleFunc("", shared.AuthMW(handleCreatePlaylist(apiCfg))).Methods("POST")
	r.HandleFunc("/audio", shared.AuthMW(handleCreatePlaylistAudio(apiCfg))).Methods("POST")

	r.HandleFunc("", shared.AuthMW(handleGetPlaylists(apiCfg))).Methods("GET")
	r.HandleFunc("/myPlaylists", shared.AuthMW(handleGetAuthUserPlaylists(apiCfg))).Methods("GET")
	r.HandleFunc("/{playlistID}", shared.AuthMW(handleGetPlaylistWithAudios(apiCfg))).Methods("GET")

	return r
}
