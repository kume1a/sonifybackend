package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func GetYoutubeSearchSuggestions(keyword string) (*youtubeSearchSuggestions, error) {
	query := url.Values{}

	query.Add("q", keyword)
	query.Add("client", "youtube")
	query.Add("hl", "en")
	query.Add("gl", "ge")

	link := "https://clients1.google.com/complete/search?" + query.Encode()

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

	suggestions, err := parseSuggestions(string(body))
	if err != nil {
		return nil, err
	}

	return &youtubeSearchSuggestions{
		Suggestions: suggestions,
		Query:       keyword,
	}, nil
}

func parseSuggestions(data string) ([]string, error) {
	re := regexp.MustCompile(`\[\[.*\]\]`)
	innerArray := re.FindString(data)

	var result []interface{}
	err := json.Unmarshal([]byte(innerArray), &result)
	if err != nil {
		return nil, err
	}

	suggestions := []string{}
	for _, suggestionWithJunk := range result {
		suggestion, ok := suggestionWithJunk.([]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse suggestion")
		}

		suggestionText, ok := suggestion[0].(string)
		if !ok {
			return nil, fmt.Errorf("failed to parse suggestion text")
		}

		suggestions = append(suggestions, suggestionText)
	}

	return suggestions, nil
}
