package auth

import "github.com/kume1a/sonifybackend/internal/shared"

type googleSignInDTO struct {
	Token string `json:"token" valid:"required"`
}

type tokenPayloadDTO struct {
	AccessToken string         `json:"accessToken"`
	User        shared.UserDto `json:"user"`
}
