package audio

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/audio").Subrouter()

	r.HandleFunc("/myAudioIds", shared.AuthMW(handleGetAuthUserAudioIds(apiCfg))).Methods("GET")
	r.HandleFunc("/myUserAudiosByIds", shared.AuthMW(handleGetAuthUserUserAudiosByIDs(apiCfg))).Methods("GET")

	r.HandleFunc("/uploadUserLocalMusic", shared.AuthMW(handleUploadUserLocalMusic(apiCfg))).Methods("POST")

	// TEMPORARY datamigration, remove after used
	r.HandleFunc("/writeInitialRelCount", handleWriteInitialAudioRelCount(apiCfg)).Methods("POST")

	return r
}
