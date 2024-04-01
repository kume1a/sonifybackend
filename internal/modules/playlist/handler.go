package playlist

import (
	"database/sql"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleCreatePlaylist(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, httpErr := ValidateCreatePlaylistDto(w, r)
		if httpErr != nil {
			shared.ResHttpError(w, *httpErr)
			return
		}

		playlist, err := CreatePlaylist(r.Context(), apiCfg.DB, body.Name, sql.NullString{String: body.ThumbnailPath, Valid: true})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		dto := playlistEntityToDto(playlist)

		shared.ResCreated(w, dto)
	}
}
