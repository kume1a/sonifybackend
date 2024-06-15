package spotify

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

const spotifyBaseApiUrl = "https://api.spotify.com"

func SpotifySearch(accessToken, query string) (*spotifySearchDTO, error) {
	return shared.HttpGetWithResponse[spotifySearchDTO](
		shared.HttpGetParams{
			URL: spotifyBaseApiUrl + "/v1/search",
			Headers: map[string]string{
				"Authorization": "Bearer " + accessToken,
			},
			Query: url.Values{
				"q":     {query},
				"type":  {"playlist"},
				"limit": {"15"},
			},
		},
	)
}

func SpotifyGetPlaylist(accessToken, playlistID string) (*spotifyPlaylistDTO, error) {
	return shared.HttpGetWithResponse[spotifyPlaylistDTO](
		shared.HttpGetParams{
			URL: spotifyBaseApiUrl + "/v1/playlists/" + playlistID,
			Headers: map[string]string{
				"Authorization": "Bearer " + accessToken,
			},
		},
	)
}

func SpotifyGetUserSavedPlaylists(accessToken string) (*spotifyGetPlaylistsDTO, error) {
	return shared.HttpGetWithResponse[spotifyGetPlaylistsDTO](
		shared.HttpGetParams{
			URL: spotifyBaseApiUrl + "/v1/me/playlists",
			Headers: map[string]string{
				"Authorization": "Bearer " + accessToken,
			},
			Query: url.Values{
				"limit": {"50"},
			},
		},
	)
}

func SpotifyGetPlaylistItems(accessToken, playlistID string) (*spotifyPlaylistItemsDTO, error) {
	return shared.HttpGetWithResponse[spotifyPlaylistItemsDTO](
		shared.HttpGetParams{
			URL: spotifyBaseApiUrl + "/v1/playlists/" + playlistID + "/tracks",
			Headers: map[string]string{
				"Authorization": "Bearer " + accessToken,
			},
			Query: url.Values{
				"limit":  {"50"},
				"fields": {"next,previous,limit,total,items(added_at,added_by(id,type),track(id,preview_url,name,duration_ms,artists(id,name),album(name,images,id)))"},
			},
		},
	)
}

func SpotifyGetPlaylistItemsNext(accessToken, nextUrl string) (*spotifyPlaylistItemsDTO, error) {
	return shared.HttpGetWithResponse[spotifyPlaylistItemsDTO](
		shared.HttpGetParams{
			URL: nextUrl,
			Headers: map[string]string{
				"Authorization": "Bearer " + accessToken,
			},
		},
	)
}

func SpotifyGetGeneralAccessToken() (*spotifyClientCredsTokenDTO, error) {
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

func SpotifyGetAuthorizationCodeTokenPayload(code string) (*spotifyAuthCodeTokenDTO, error) {
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

func SpotifyRefreshToken(refreshToken string) (*spotifyRefreshTokenDTO, error) {
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
