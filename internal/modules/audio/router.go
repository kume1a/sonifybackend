package audio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/audio").Subrouter()

	r.HandleFunc("/downloadYoutubeAudio", shared.AuthMW(handleDownloadYoutubeAudio(apiCfg))).Methods("POST")
	r.HandleFunc("/importUserLocalMusic", shared.AuthMW(handleImportUserLocalMusic(apiCfg))).Methods("POST")

	return r
}
