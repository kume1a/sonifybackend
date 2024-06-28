package spotify

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/spotify").Subrouter()

	r.HandleFunc("/search", shared.AuthMW(handleSpotifySearch(apiCfg))).Methods("GET")

	r.HandleFunc("/authorize", shared.AuthMW(handleAuthorizeSpotify)).Methods("POST")
	r.HandleFunc("/refreshToken", shared.AuthMW(handleSpotifyRefreshToken)).Methods("POST")
	r.HandleFunc("/importPlaylist", shared.AuthMW(handleImportSpotifyPlaylist(apiCfg))).Methods("POST")
	r.HandleFunc("/importUserPlaylists", shared.AuthMW(handleImportSpotifyUserPlaylists(apiCfg))).Methods("POST")

	return r
}
