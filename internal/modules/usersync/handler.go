package usersync

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetUserSyncDatumByUserId(apiCfg *shared.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		syncData, httpErr := GetOrCreateUserSyncDatumByUserId(r.Context(), apiCfg.DB, tokenPayload.UserId)
		if httpErr != nil {
			shared.ResHttpError(w, *httpErr)
			return
		}

		dto := userSyncDatumEntityToDTO(syncData)

		shared.ResOK(w, dto)
	}
}
