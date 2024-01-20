package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGoogleSignIn(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequest[googleSignInDTO](r)
		if err != nil {
			shared.ResError(w, http.StatusBadRequest, err.Error())
			return
		}

		shared.ResJson(w, http.StatusOK, body)
	}
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/auth").Subrouter()

	r.HandleFunc("/googleSignIn", handleGoogleSignIn(apiCfg)).Methods("POST")

	return r
}
