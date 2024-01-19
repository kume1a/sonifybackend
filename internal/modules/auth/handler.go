package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func signInHandler(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle sign in logic here
	}
}

func signUpHandler(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle sign up logic here
	}
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/auth").Subrouter()

	r.HandleFunc("/signIn", signInHandler(apiCfg)).Methods("POST")
	r.HandleFunc("/signUp", signUpHandler(apiCfg)).Methods("POST")

	return r
}
