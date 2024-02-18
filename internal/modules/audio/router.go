package audio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/modules/auth"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/audio").Subrouter()

	r.HandleFunc("/downloadYoutubeAudio", auth.AuthMW(handleDownloadYoutubeAudio(apiCfg))).Methods("POST")

	return r
}
