package spotify

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetSpotifyAudioDownloadMeta(trackID string) (*downloadSpotifyTrackMetaDTO, error) {
	req, err := http.NewRequest("GET", "https://api.spotifydown.com/download/"+trackID, nil)
	if err != nil {
		return nil, err
	}

	// req.Header.Set(":authority", "api.spotifydown.com")
	// req.Header.Set(":method", "GET")
	// req.Header.Set(":path", "/download/"+trackID)
	// req.Header.Set(":scheme", "https")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("If-None-Match", "W/\"1df-smzelj558DY6eHWjMrN79V/tNhU\"")
	req.Header.Set("Origin", "https://spotifydown.com")
	req.Header.Set("Referer", "https://spotifydown.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

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

	dto := downloadSpotifyTrackMetaDTO{}
	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Println("error unmarshalling response: ", err)
		return nil, err
	}

	if !dto.Success {
		log.Println("not success download meta for URL ", "https://api.spotifydown.com/download/"+trackID)
	}

	return &dto, nil
}
