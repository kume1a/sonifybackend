package playlist

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func ValidateGetPlaylistByIDVars(r *http.Request) (*playlistIDDTO, error) {
	vars := mux.Vars(r)

	playlistID, ok := vars["playlistID"]
	if !ok {
		return nil, shared.BadRequest("playlistID is required")
	}

	playlistIDUUID, err := uuid.Parse(playlistID)
	if err != nil {
		return nil, shared.BadRequest("playlistId is not a valid UUID")
	}

	return &playlistIDDTO{PlaylistID: playlistIDUUID}, nil
}
