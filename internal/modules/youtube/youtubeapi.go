package youtube

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetYoutubeSearchSuggestions(keyword string) (*youtubeSearchSuggestions, error) {
	query := url.Values{}

	query.Add("q", keyword)

	link := "https://invidious.slipfox.xyz/api/v1/search/suggestions?" + query.Encode()

	response, err := http.Get(link)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	var res youtubeSearchSuggestions
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("Error: unmarshal failed:  %v", err)
		return nil, err
	}

	return &res, nil
}
