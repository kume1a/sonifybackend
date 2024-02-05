package youtube

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/youtube").Subrouter()

	r.HandleFunc("/musicUrl", handleGetYoutubeMusicUrl).Methods("GET")
	r.HandleFunc("/searchSuggestions", handleGetYoutubeSearchSuggestions).Methods("GET")

	return r
}
