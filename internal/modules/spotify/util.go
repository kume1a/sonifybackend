package spotify

import (
	"encoding/base64"

	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/config"
)

func (dto *downloadSpotifyPlaylistDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *authorizeSpotifyDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *spotifyAccessTokenDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *refreshSpotifyTokenDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func (dto *searchSpotifyQueryDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func getSpotifyBasicAuthHeader() (string, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return "", err
	}

	return "Basic " + base64.StdEncoding.EncodeToString(
		[]byte(env.SpotifyClientID+":"+env.SpotifyClientSecret),
	), nil
}
