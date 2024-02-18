package auth

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func AuthMW(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, ok := r.Header["Authorization"]
		if !ok {
			shared.ResUnauthorized(w, shared.ErrMissingToken)
			return
		}

		if _, err := VerifyAccessToken(accessToken[0]); err != nil {
			shared.ResUnauthorized(w, shared.ErrInvalidToken)
			return
		}

		h.ServeHTTP(w, r)
	})
}
