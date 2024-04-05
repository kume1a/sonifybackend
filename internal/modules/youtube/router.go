package youtube

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/youtube").Subrouter()

	r.HandleFunc("/musicUrl", shared.AuthMW(handleGetYoutubeMusicUrl)).Methods("GET")
	r.HandleFunc("/searchSuggestions", shared.AuthMW(handleGetYoutubeSearchSuggestions)).Methods("GET")

	return r
}
