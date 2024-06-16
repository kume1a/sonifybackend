package spotify

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/usersync"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleSpotifySearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := shared.ValidateRequestQuery[*searchSpotifyQueryDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		if len(query.Keyword) != 1 || len(query.SpotifyAccessToken) != 1 {
			shared.ResBadRequest(w, shared.ErrInvalidQueryParams)
			return
		}

		spotifyRes, err := SpotifySearch(query.SpotifyAccessToken[0], query.Keyword[0])
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := MapSpotifySearchToSearchSpotifyResult(spotifyRes)

		shared.ResOK(w, dto)
	}
}

func handleDownloadPlaylist(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		body, err := shared.ValidateRequestBody[*downloadSpotifyPlaylistDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		// TODO implement
		log.Println("AuthPayload: ", authPayload, "Body: ", body)

		shared.ResInternalServerError(w, shared.ErrNotImplemented)
	}
}

func handleAuthorizeSpotify(w http.ResponseWriter, r *http.Request) {
	body, err := shared.ValidateRequestBody[*authorizeSpotifyDTO](r)
	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	tokenPayload, err := SpotifyGetAuthorizationCodeTokenPayload(body.Code)
	if err != nil {
		shared.ResInternalServerError(w, shared.ErrFailedToGetSpotifyAccessToken)
		return
	}

	dto := spotifyAuthCodeTokenPayloadDTO{
		AccessToken:  tokenPayload.AccessToken,
		RefreshToken: tokenPayload.RefreshToken,
		Scope:        tokenPayload.Scope,
		ExpiresIn:    tokenPayload.ExpiresIn,
		TokenType:    tokenPayload.TokenType,
	}

	shared.ResOK(w, dto)
}

func handleSpotifyRefreshToken(w http.ResponseWriter, r *http.Request) {
	body, err := shared.ValidateRequestBody[*refreshSpotifyTokenDTO](r)
	if err != nil {
		shared.ResBadRequest(w, err.Error())
		return
	}

	tokenPayload, err := SpotifyRefreshToken(body.SpotifyRefreshToken)
	if err != nil {
		shared.ResInternalServerError(w, shared.ErrFailedToGetSpotifyAccessToken)
		return
	}

	dto := spotifyRefreshTokenPayloadDTO{
		AccessToken: tokenPayload.AccessToken,
		Scope:       tokenPayload.Scope,
		ExpiresIn:   tokenPayload.ExpiresIn,
		TokenType:   tokenPayload.TokenType,
	}

	shared.ResOK(w, dto)
}

func handleImportSpotifyUserPlaylists(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		query, err := shared.ValidateRequestQuery[*spotifyAccessTokenDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		if err := downloadSpotifyUserSavedPlaylists(r.Context(), apiCfg, authPayload.UserID, query.SpotifyAccessToken[0]); err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		if _, httpErr := usersync.UpdateUserSyncDatumByUserId(
			r.Context(),
			apiCfg.DB,
			database.UpdateUserSyncDatumByUserIDParams{
				UserID:              authPayload.UserID,
				SpotifyLastSyncedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
			},
		); httpErr != nil {
			shared.ResHttpError(w, httpErr)
			return
		}

		shared.ResNoContent(w)
	}
}
