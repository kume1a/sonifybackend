package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		log.Println("Error getting google public key: ", err)
		return "", err
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading google public key: ", err)
		return "", err
	}

	myResp := map[string]string{}
	if err := json.Unmarshal(dat, &myResp); err != nil {
		log.Println("Error unmarshalling google public key: ", err)
		return "", err
	}

	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}

	return key, nil
}

func ValidateGoogleJWT(tokenString string) (*GoogleClaims, error) {
	env, err := shared.ParseEnv()
	if err != nil {
		return nil, err
	}

	claimsStruct := GoogleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				log.Println("Error getting google public key: ", err)
				return nil, err
			}

			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				log.Println("Error parsing google public key: ", err)
				return nil, err
			}

			return key, nil
		},
	)
	if err != nil {
		log.Println("Error parsing google JWT: ", err)
		return nil, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return nil, errors.New("Invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		log.Println("iss is invalid, expected: accounts.google.com, got: ", claims.Issuer)
		return nil, errors.New("iss is invalid")
	}

	if claims.Audience != env.GoogleClientKey {
		log.Println("aud is invalid, expected: ", env.GoogleClientKey, " got: ", claims.Audience)
		return nil, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		log.Println("JWT is expired, exp: ", claims.ExpiresAt, " now: ", time.Now().UTC().Unix())
		return nil, errors.New("JWT is expired")
	}

	return claims, nil
}
