package youtubemusic

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetYoutubeMusicUrl(w http.ResponseWriter, r *http.Request) {
	query, err := shared.ValidateRequestQuery[*getYoutubeMusicDto](r)
	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	url, err := GetYoutubeMusicUrl(query.VideoID[0])
	if err != nil {
		shared.ResInternalServerErrorDef(w)
		return
	}

	dto := shared.UrlDto{Url: url}
	shared.ResOK(w, dto)
}

func handleGetYoutubeSearchSuggestions(w http.ResponseWriter, r *http.Request) {
	query, err := shared.ValidateRequestQuery[*shared.KeywordDto](r)

	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	res, err := GetYoutubeSearchSuggestions(query.Keyword[0])
	if err != nil {
		shared.ResInternalServerErrorDef(w)
		return
	}

	shared.ResOK(w, res)
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/youtubeMusic").Subrouter()

	r.HandleFunc("/musicUrl", handleGetYoutubeMusicUrl).Methods("GET")
	r.HandleFunc("/searchSuggestions", handleGetYoutubeSearchSuggestions).Methods("GET")

	return r
}
