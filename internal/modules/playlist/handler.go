package playlist

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/playlistaudio"
	"github.com/kume1a/sonifybackend/internal/modules/sharedmodule"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGetPlaylistFull(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPaylad, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrUnauthorized)
			return
		}

		playlistIDDTO, err := ValidateGetPlaylistByIDVars(r)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		ctx := r.Context()

		playlist, err := GetPlaylistByID(ctx, apiCfg.DB, playlistIDDTO.PlaylistID)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		playlistAudios, err := playlistaudio.GetPlaylistAudios(
			ctx, apiCfg.DB,
			database.GetPlaylistAudiosParams{
				PlaylistIds: []uuid.UUID{playlistIDDTO.PlaylistID},
				UserID:      authPaylad.UserID,
				Ids:         nil,
			},
		)
		if err != nil {
			shared.ResTryHttpError(w, err)
			return
		}

		dto := PlaylistFullDTO{
			PlaylistDTO:    sharedmodule.PlaylistEntityToDTO(playlist),
			PlaylistAudios: shared.Map(playlistAudios, playlistaudio.GetPlaylistAudioRowToDTO),
		}

		shared.ResOK(w, dto)
	}
}
