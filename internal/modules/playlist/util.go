package playlist

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func ValidateCreatePlaylistDto(w http.ResponseWriter, r *http.Request) (*createPlaylistDTO, error) {
	thumbnailPath, err := shared.HandleUploadFile(shared.HandleUploadFileArgs{
		ResponseWriter:   w,
		Request:          r,
		FieldName:        "thumbnail",
		Dir:              config.DirPlaylistThumbnails,
		AllowedMimeTypes: shared.ImageMimeTypes,
		IsOptional:       false,
	})
	if err != nil {
		return nil, err
	}

	name := r.FormValue("name")

	if !govalidator.IsByteLength(name, 1, 255) {
		return nil, shared.BadRequest("Name must be between 1 and 255 characters")
	}

	return &createPlaylistDTO{
		Name:          name,
		ThumbnailPath: thumbnailPath,
	}, nil
}

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
