package auth

import (
	"errors"
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

func VerifyAccessToken(tokenString string) (*TokenClaims, error) {
	env, err := shared.ParseEnv()
	if err != nil {
		return nil, err
	}

	return verifyJWT(tokenString, env.AccessTokenSecret)
}

func verifyJWT(tokenString string, secretKey string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return nil, errors.New("invalid userId")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email")
	}

	return &TokenClaims{
		UserId: userId,
		Email:  email,
	}, nil
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
