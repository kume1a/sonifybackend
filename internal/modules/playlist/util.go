package playlist

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func (dto *createPlaylistAudioDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func ValidateCreatePlaylistDto(w http.ResponseWriter, r *http.Request) (*createPlaylistDTO, *shared.HttpError) {
	thumbnailPath, err := shared.HandleUploadFile(
		w, r,
		"thumbnail",
		shared.DirPlaylistThumbnails,
		shared.ImageMimeTypes,
	)
	if err != nil {
		return nil, err
	}

	name := r.FormValue("name")

	if !govalidator.IsByteLength(name, 1, 255) {
		return nil, shared.HttpErrBadRequest("Name must be between 1 and 255 characters")
	}

	return &createPlaylistDTO{
		Name:          name,
		ThumbnailPath: thumbnailPath,
	}, nil
}

func playlistEntityToDto(e database.Playlist) playlistDTO {
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
