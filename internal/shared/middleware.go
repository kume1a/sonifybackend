package shared

import (
	"errors"
	"net/http"
	"strings"
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

func GetAccessTokenFromRequest(r *http.Request) (string, error) {
	accessToken, ok := r.Header["Authorization"]
	if !ok {
		return "", errors.New(ErrMissingToken)
	}

	if len(accessToken) == 0 {
		return "", errors.New(ErrInvalidToken)
	}

	accessToken[0] = strings.Replace(accessToken[0], "Bearer ", "", 1)

	return accessToken[0], nil
}
