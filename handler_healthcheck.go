package main

import "net/http"

type HealthcheckRes struct {
	Ok bool `json:"ok"`
}

func handlerHealthcheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, HealthcheckRes{Ok: true})
}
