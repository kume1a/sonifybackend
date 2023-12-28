package modules

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

type HealthcheckRes struct {
	Ok bool `json:"ok"`
}

func HandlerHealthcheck(w http.ResponseWriter, r *http.Request) {
	shared.RespondWithJSON(w, 200, HealthcheckRes{Ok: true})
}
