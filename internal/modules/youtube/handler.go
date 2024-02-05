package youtube

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetYoutubeMusicUrl(w http.ResponseWriter, r *http.Request) {
	query, err := shared.ValidateRequestQuery[*getYoutubeMusicUrlDto](r)
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
