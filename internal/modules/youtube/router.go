package youtube

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/youtube").Subrouter()

	r.HandleFunc("/downloadAudioToUserLibrary", shared.AuthMW(handleDownloadYoutubeAudioToUserLibrary(apiCfg))).Methods("POST")
	r.HandleFunc("/downloadAudioToPlaylist", shared.AuthMW(handleDownloadYoutubeAudioPlaylist(apiCfg))).Methods("POST")

	r.HandleFunc("/searchSuggestions", shared.AuthMW(handleGetYoutubeSearchSuggestions)).Methods("GET")

	return r
}
