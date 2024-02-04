package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/modules/user"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGoogleSignIn(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*googleSignInDTO](r)

		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		claims, err := ValidateGoogleJWT(body.Token)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrInvalidGoogleToken)
			return
		}

		user, err := user.GetUserByEmail(apiCfg.DB, r.Context(), claims.Email)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUserNotFound)
			return
		}

		tokenString, err := GenerateAccessToken(&TokenClaims{
			UserId: user.ID.String(),
			Email:  user.Email.String,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, struct {
			Token string `json:"token"`
		}{
			Token: tokenString,
		})
	}
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/auth").Subrouter()

	r.HandleFunc("/googleSignIn", handleGoogleSignIn(apiCfg)).Methods("POST")

	return r
}
