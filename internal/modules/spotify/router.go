package spotify

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/spotify").Subrouter()

	r.HandleFunc("/downloadPlaylist", shared.AuthMW(handleDownloadPlaylist(apiCfg))).Methods("POST")

	return r
}
