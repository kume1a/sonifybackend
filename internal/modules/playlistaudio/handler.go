package playlistaudio

import (
	"net/http"

	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/userplaylist"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleCreatePlaylistAudio(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		body, err := shared.ValidateRequestBody[*createPlaylistAudioDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		exists, err := userplaylist.UserPlaylistExistsByUserIDAndPlaylistID(
			r.Context(),
			apiCfg.DB,
			database.UserPlaylistExistsByUserIDAndPlaylistIDParams{
				UserID:     authPayload.UserID,
				PlaylistID: body.PlaylistID,
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		if !exists {
			shared.ResForbidden(w, shared.ErrUserPlaylistNotFound)
			return
		}

		playlistAudio, err := CreatePlaylistAudioTx(
			r.Context(),
			apiCfg.ResourceConfig,
			database.CreatePlaylistAudioParams{
				PlaylistID: body.PlaylistID,
				AudioID:    body.AudioID,
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := playlistAudioEntityToDTO(playlistAudio)

		shared.ResCreated(w, dto)
	}
}

func handleGetPlaylistAudiosByAuthUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		// using body for big payload
		body, err := shared.ValidateRequestBody[*shared.RequiredIDsDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		userPlaylistIDs, err := userplaylist.GetPlaylistIDsByUserID(
			r.Context(),
			apiCfg.DB,
			authPayload.UserID,
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		playlistAudios, err := GetPlaylistAudios(
			r.Context(), apiCfg.DB,
			database.GetPlaylistAudiosParams{
				UserID:      authPayload.UserID,
				PlaylistIds: userPlaylistIDs,
				Ids:         body.IDs,
			},
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dtos := getPlaylistAudioRowListToDTO(playlistAudios)

		shared.ResOK(w, dtos)
	}
}

func handleGetPlaylistAudioIDsByAuthUser(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		playlistAudioIDs, err := GetPlaylistAudioIDsByUserID(
			r.Context(),
			apiCfg.DB,
			authPayload.UserID,
		)
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, playlistAudioIDs)
	}
}
