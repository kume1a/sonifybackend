package usersync

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetUserSyncDatum(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		syncData, err := GetUserSyncDatumByUserId(r.Context(), apiCfg.DB, tokenPayload.UserId)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := userSyncDatumEntityToDTO(syncData)

		shared.ResOK(w, dto)
	}
}
