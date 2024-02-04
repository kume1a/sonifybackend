package youtubemusic

import (
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

	url, err := GetYoutubeMusicUrl(body.VideoID)
	if err != nil {
		shared.ResInternalServerErrorDef(w)
		return
	}

	dto := shared.UrlDto{Url: url}
	shared.ResOK(w, dto)
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/youtubeMusic").Subrouter()

	r.HandleFunc("/musicUrl", handleGetYoutubeMusicUrl).Methods("GET")

	return r
}
