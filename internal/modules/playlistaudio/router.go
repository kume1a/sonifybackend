package playlistaudio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/playlistaudio").Subrouter()

	r.HandleFunc("", shared.AuthMW(handleGetPlaylistAudioIdsByAuthUser(apiCfg))).Methods("GET")
	r.HandleFunc("/{playlistID}", shared.AuthMW(handleGetPlaylistWithAudios(apiCfg))).Methods("GET")

	return r
}
