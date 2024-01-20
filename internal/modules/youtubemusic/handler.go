package youtubemusic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetYoutubeMusicUrl(w http.ResponseWriter, r *http.Request) {
	body, err := shared.ValidateRequest[*getYoutubeMusicDto](r)
	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	headers := map[string]string{
		"Accept":             "*/*",
		"Accept-Encoding":    "application/json",
		"Accept-Language":    "en-GB,en-US;q=0.9,en;q=0.8",
		"Referer":            "https://ytjar.downloader-ytjar.online/1.php",
		"Sec-Ch-Ua":          "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"",
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": "\"macOS\"",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"X-Requested-With":   "XMLHttpRequest",
	}

	url := fmt.Sprintf("https://ytjar.downloader-ytjar.online/download1.php?id=%s", body.VideoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		shared.ResInternalServerErrorDef(w)
		return
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		shared.ResInternalServerErrorDef(w)
		return
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		shared.ResInternalServerErrorDef(w)
		return
	}

	shared.ResOK(w, data)
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/youtubeMusic").Subrouter()

	r.HandleFunc("/musicUrl", handleGetYoutubeMusicUrl).Methods("GET")

	return r
}
