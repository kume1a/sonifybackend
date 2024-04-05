package user

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfig, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/users").Subrouter()

	r.HandleFunc("", shared.AuthMW(handleUpdateUser(apiCfg))).Methods("PATCH")
	r.HandleFunc("/authUser", shared.AuthMW(handleGetAuthUser(apiCfg))).Methods("GET")

	return r
}
