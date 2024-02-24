package shared

import (
	"net/http"
)

func AuthMW(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := GetAccessTokenFromRequest(r)
		if err != nil {
			ResUnauthorized(w, err.Error())
			return
		}

		if _, err := VerifyAccessToken(accessToken); err != nil {
			ResUnauthorized(w, ErrInvalidToken)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func GetAuthPayload(r *http.Request) (*TokenClaims, error) {
	accessToken, err := GetAccessTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	return VerifyAccessToken(accessToken)
}
