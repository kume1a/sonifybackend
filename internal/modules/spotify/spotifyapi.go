package spotify

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func GetSpotifyPlaylist(accessToken, playlistID string) (*spotifyPlaylistDTO, error) {
	return getSpotifyEndpoint[spotifyPlaylistDTO](
		"/v1/playlists/"+playlistID,
		accessToken,
		nil,
	)
}

func GetSavedSpotifyPlaylists(accessToken string) (*getSpotifyPlaylistsDTO, error) {
	return getSpotifyEndpoint[getSpotifyPlaylistsDTO](
		"/v1/me/playlists",
		accessToken,
		url.Values{
			"limit": {"50"},
		},
	)
}

func GetSpotifyPlaylistItems(accessToken, playlistID string) (*spotifyPlaylistItemsDTO, error) {
	return getSpotifyEndpoint[spotifyPlaylistItemsDTO](
		"/v1/playlists/"+playlistID+"/tracks",
		accessToken,
		url.Values{
			"limit":  {"50"},
			"fields": {"next,previous,limit,total,items(added_at,added_by(id,type),track(id,preview_url,name,duration_ms,artists(id,name),album(name,images,id)))"},
		},
	)
}

func GetGeneralSpotifyAccessToken() (*spotifyClientCredsTokenDTO, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return nil, err
	}

	resp, respBody, err := shared.XWWWFormUrlencoded(
		shared.XWWWFormUrlencodedParams{
			URL: "https://accounts.spotify.com/api/token",
			Form: url.Values{
				"grant_type":    {"client_credentials"},
				"client_id":     {env.SpotifyClientID},
				"client_secret": {env.SpotifyClientSecret},
			},
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			},
		},
	)

	if resp.StatusCode != http.StatusOK {
		log.Println("error getting general spotify access token: ", resp.StatusCode, ", body = ", respBody)
		return nil, errors.New("error getting general spotify access token")
	}

	dto := spotifyClientCredsTokenDTO{}
	if err := json.Unmarshal([]byte(respBody), &dto); err != nil {
		return nil, err
	}

	return &dto, err
}

func GetAuthorizationCodeSpotifyTokenPayload(code string) (*spotifyAuthCodeTokenDTO, error) {
	env, err := config.ParseEnv()
	if err != nil {
		return nil, err
	}

	basicAuth, err := getSpotifyBasicAuthHeader()
	if err != nil {
		return nil, err
	}

	resp, respBody, err := shared.XWWWFormUrlencoded(
		shared.XWWWFormUrlencodedParams{
			URL: "https://accounts.spotify.com/api/token",
			Form: url.Values{
				"grant_type":   {"authorization_code"},
				"code":         {code},
				"redirect_uri": {env.SpotifyRedirectURI},
			},
			Headers: map[string]string{
				"Content-Type":  "application/x-www-form-urlencoded",
				"Authorization": basicAuth,
			},
		},
	)

	if resp.StatusCode != http.StatusOK {
		log.Println("error getting auth code: ", resp.StatusCode, ", body = ", respBody)
		return nil, errors.New("error getting auth code")
	}

	dto := spotifyAuthCodeTokenDTO{}
	if err := json.Unmarshal([]byte(respBody), &dto); err != nil {
		return nil, err
	}

	return &dto, err
}

func RefreshSpotifyToken(refreshToken string) (*spotifyRefreshTokenDTO, error) {
	basicAuthHeader, err := getSpotifyBasicAuthHeader()
	if err != nil {
		return nil, err
	}

	resp, respBody, err := shared.XWWWFormUrlencoded(
		shared.XWWWFormUrlencodedParams{
			URL: "https://accounts.spotify.com/api/token",
			Form: url.Values{
				"grant_type":    {"refresh_token"},
				"refresh_token": {refreshToken},
			},
			Headers: map[string]string{
				"Authorization": basicAuthHeader,
				"Content-Type":  "application/x-www-form-urlencoded",
			},
		},
	)

	if resp.StatusCode != http.StatusOK {
		log.Println("error refreshing spotify token: ", resp.StatusCode, ", body = ", respBody)
		return nil, errors.New("error refreshing spotify token")
	}

	dto := spotifyRefreshTokenDTO{}
	if err := json.Unmarshal([]byte(respBody), &dto); err != nil {
		return nil, err
	}

	return &dto, err
}

func getSpotifyEndpoint[DTO interface{}](endpoint, accessToken string, queryParams url.Values) (*DTO, error) {
	url := "https://api.spotify.com" + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	req.URL.RawQuery = queryParams.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dto DTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}
