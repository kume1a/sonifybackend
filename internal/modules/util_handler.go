package modules

import (
	"net/http"
	"time"

	"github.com/kume1a/sonifybackend/internal/shared"
)

type HealthcheckDTO struct {
	Ok bool `json:"ok"`
}

type ServerTimeDTO struct {
	Time time.Time `json:"time"`
}

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	shared.ResOK(w, HealthcheckDTO{Ok: true})
}

func handleGetServerTime(w http.ResponseWriter, r *http.Request) {
	shared.ResOK(w, ServerTimeDTO{Time: time.Now().UTC()})
}
