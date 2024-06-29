package spotify

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/usersync"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleSpotifySearch(apiCfg *config.ApiConfig) http.HandlerFunc {
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

		spotifyResMergedWithDb, err := mergeSpotifySearchWithDBPlaylists(r.Context(), apiCfg.ResourceConfig, spotifyRes)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := MapMergedSpotifySearchToSearchSpotifyResult(spotifyRes, spotifyResMergedWithDb)

		shared.ResOK(w, dto)
	}
}

func handleImportSpotifyPlaylist(apiCfg *config.ApiConfig) http.HandlerFunc {
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

		if err := downloadSpotifyPlaylist(
			r.Context(), apiCfg,
			authPayload.UserID,
			body.SpotifyPlaylistID,
			body.SpotifyAccessToken,
		); err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResNoContent(w)
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

		body, err := shared.ValidateRequestBody[*spotifyAccessTokenDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		if err := downloadSpotifyUserSavedPlaylists(
			r.Context(), apiCfg,
			authPayload.UserID,
			body.SpotifyAccessToken,
		); err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		if _, err := usersync.UpdateUserSyncDatumByUserId(
			r.Context(),
			apiCfg.DB,
			database.UpdateUserSyncDatumByUserIDParams{
				UserID:              authPayload.UserID,
				SpotifyLastSyncedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
			},
		); err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		shared.ResNoContent(w)
	}
}
