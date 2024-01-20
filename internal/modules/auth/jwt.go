package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type TokenClaims struct {
	UserId string
	Email  string
}

func GenerateAccessToken(tokenClaims *TokenClaims) (string, error) {
	env, err := shared.ParseEnv()
	if err != nil {
		return "", err
	}

	return generateJWT(tokenClaims, env.AccessTokenSecret)
}

func generateJWT(tokenClaims *TokenClaims, secretKey string) (string, error) {
	env, err := shared.ParseEnv()
	if err != nil {
		return "", err
	}

	expDuration := time.Millisecond * time.Duration(env.AccessTokenExp)

	claims := jwt.MapClaims{
		"userId": tokenClaims.UserId,
		"email":  tokenClaims.Email,
		"exp":    time.Now().Add(expDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
