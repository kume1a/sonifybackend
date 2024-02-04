package youtubemusic

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetYoutubeMusicUrl(w http.ResponseWriter, r *http.Request) {
	body, err := shared.ValidateRequest[*getYoutubeMusicDto](r)
	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	// Replace with the URL of the video you want to download
	videoURL := "https://www.youtube.com/watch?v=_LWXAPywCV4"

	// Replace with the format code of the bitrate you want to download
	formatCode := "140"

	cmd := exec.Command("yt-dlp", "-f", formatCode, videoURL)

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	shared.ResOK(w, body)
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/youtubeMusic").Subrouter()

	r.HandleFunc("/musicUrl", handleGetYoutubeMusicUrl).Methods("GET")

	return r
}
