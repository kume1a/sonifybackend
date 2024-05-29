package playlist

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func (dto *createPlaylistAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func ValidateCreatePlaylistDto(w http.ResponseWriter, r *http.Request) (*createPlaylistDTO, *shared.HttpError) {
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

func ValidateGetPlaylistByIDVars(r *http.Request) (*playlistIDDTO, *shared.HttpError) {
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

func PlaylistEntityToDto(e database.Playlist) playlistDTO {
	return playlistDTO{
		ID:            e.ID,
		CreatedAt:     e.CreatedAt,
		Name:          e.Name,
		ThumbnailPath: e.ThumbnailPath.String,
		ThumbnailUrl:  e.ThumbnailUrl.String,
		SpotifyId:     e.SpotifyID.String,
	}
}

func playlistAudioEntityToDto(e *database.PlaylistAudio) playlistAudioDTO {
	return playlistAudioDTO{
		CreatedAt:  e.CreatedAt,
		PlaylistID: e.PlaylistID,
		AudioID:    e.AudioID,
	}
}
