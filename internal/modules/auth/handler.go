package auth

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGoogleAuth(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*googleSignInDTO](r)

		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		tokenPayload, httpErr := AuthWithGoogle(*apiCfg, r.Context(), body.Token)
		if httpErr != nil {
			shared.ResHttpError(w, *httpErr)
			return
		}

		shared.ResOK(w, tokenPayload)
	}
}

func handleEmailAuth(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*emailSignInDTO](r)

		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		tokenPayload, httpErr := AuthWithEmail(*apiCfg, r.Context(), body.Email, body.Password)
		if httpErr != nil {
			shared.ResHttpError(w, *httpErr)
			return
		}

		shared.ResOK(w, tokenPayload)
	}
}
