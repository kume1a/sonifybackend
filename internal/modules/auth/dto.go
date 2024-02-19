package auth

type googleSignInDTO struct {
	Token string `json:"token" valid:"required"`
}

type tokenPayloadDTO struct {
	AccessToken string `json:"accessToken"`
}
