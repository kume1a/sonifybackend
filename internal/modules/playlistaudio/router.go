package playlistaudio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/playlistaudio").Subrouter()

	r.HandleFunc("", shared.AuthMW(handleCreatePlaylistAudio(apiCfg))).Methods("POST")
	r.HandleFunc("", shared.AuthMW(handleDeletePlaylistAudio(apiCfg))).Methods("DELETE")

	r.HandleFunc("/myPlaylistAudios", shared.AuthMW(handleGetPlaylistAudiosByAuthUser(apiCfg))).Methods("GET")
	r.HandleFunc("/myPlaylistAudioIds", shared.AuthMW(handleGetPlaylistAudioIDsByAuthUser(apiCfg))).Methods("GET")

	return r
}
