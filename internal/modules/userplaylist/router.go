package userplaylist

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *config.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/userplaylist").Subrouter()

	// POST
	r.HandleFunc(
		"",
		shared.AuthMW(handleCreateUserPlaylist(apiCfg)),
	).Methods("POST")

	// PATCH
	r.HandleFunc(
		"/{userPlaylistID}",
		shared.AuthMW(handleUpdateUserPlaylist(apiCfg)),
	).Methods("PATCH")

	// DELETE
	r.HandleFunc(
		"/{userPlaylistID}",
		shared.AuthMW(handleDeleteUserPlaylist(apiCfg)),
	).Methods("DELETE")

	// GET
	r.HandleFunc(
		"/myUserPlaylists/full",
		shared.AuthMW(handleGetUserPlaylistsFullByAuthUserID(apiCfg)),
	).Methods("GET")

	r.HandleFunc(
		"/myUserPlaylists",
		shared.AuthMW(handleGetUserPlaylistsByAuthUserID(apiCfg)),
	).Methods("GET")

	r.HandleFunc(
		"/myPlaylistIds",
		shared.AuthMW(handleGetPlaylistIDsByAuthUser(apiCfg)),
	).Methods("GET")

	r.HandleFunc(
		"/myUserPlaylistIds",
		shared.AuthMW(handleGetUserPlaylistIDsByAuthUser(apiCfg)),
	).Methods("GET")

	return r
}
