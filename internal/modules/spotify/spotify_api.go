package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

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

func GetGeneralSpotifyAccessToken() (*getSpotifyGeneralTokenDTO, error) {
	env, err := shared.ParseEnv()
	if err != nil {
		return nil, err
	}

	resp, err := http.PostForm("https://api.spotify.com/v1/token", url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {env.SpotifyClientID},
		"client_secret": {env.SpotifyClientSecret},
	})

	if err != nil {
		log.Println("error sending request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dto := getSpotifyGeneralTokenDTO{}
	if err := json.Unmarshal(body, &dto); err != nil {
		return nil, err
	}

	return &dto, nil
}

func GetAuthorizationCodeSpotifyTokenPayload(code string) (*getAuthorizationCodeSpotifyTokenPayloadDTO, error) {
	env, err := shared.ParseEnv()
	if err != nil {
		return nil, err
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString(
		[]byte(env.SpotifyClientID+":"+env.SpotifyClientSecret),
	)

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {env.SpotifyRedirectURI},
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", basicAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("status code: ", resp.StatusCode, " body: ", string(body))
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	dto := getAuthorizationCodeSpotifyTokenPayloadDTO{}
	if err := json.Unmarshal(body, &dto); err != nil {
		return nil, err
	}

	return &dto, nil
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
