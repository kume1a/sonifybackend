package youtube

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetYoutubeSearchSuggestions(keyword string) (*youtubeSearchSuggestions, error) {
	link := "https://invidious.slipfox.xyz/api/v1/search/suggestions?q=" + keyword

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
		return nil, err
	}

	return &res, nil
}
