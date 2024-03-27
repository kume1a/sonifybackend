package auth

import (
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/auth").Subrouter()

	r.HandleFunc("/googleSignIn", handleGoogleAuth(apiCfg)).Methods("POST")
	r.HandleFunc("/emailSignIn", handleEmailAuth(apiCfg)).Methods("POST")

	return r
}
