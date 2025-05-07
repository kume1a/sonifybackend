package usersync

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetUserSyncDatumByUserId(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		syncData, httpErr := GetOrCreateUserSyncDatumByUserId(r.Context(), apiCfg.DB, tokenPayload.UserID)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		dto := userSyncDatumEntityToDTO(syncData)

		shared.ResOK(w, dto)
	}
}

func handleMarkUserAudioLastUpdatedAtAsNow(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		_, httpErr := GetOrCreateUserSyncDatumByUserId(r.Context(), apiCfg.DB, tokenPayload.UserID)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		_, httpErr = UpsertUserSyncDatumByUserId(
			r.Context(),
			apiCfg.DB,
			&database.UserSyncDatum{
				UserID:                tokenPayload.UserID,
				UserAudioLastSyncedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
			},
		)
		if httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		shared.ResNoContent(w)
	}
}
