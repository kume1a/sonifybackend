package playlist

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func (dto *createPlaylistAudioDto) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}

func ValidateCreatePlaylistDto(w http.ResponseWriter, r *http.Request) (*createPlaylistDto, *shared.HttpError) {
	thumbnailPath, err := shared.HandleUploadFile(
		w, r,
		"thumbnail",
		shared.DirPlaylistThumbnails,
		[]string{"image/jpeg", "image/png"},
	)
	if err != nil {
		return nil, err
	}

	name := r.FormValue("name")

	if !govalidator.IsByteLength(name, 1, 255) {
		return nil, shared.HttpErrBadRequest("Name must be between 1 and 255 characters")
	}

	return &createPlaylistDto{
		Name:          name,
		ThumbnailPath: thumbnailPath,
	}, nil
}

func playlistEntityToDto(e database.Playlist) playlistDto {
	return playlistDto{
		ID:            e.ID,
		Name:          e.Name,
		ThumbnailPath: e.ThumbnailPath.String,
	}
}

func playlistAudioEntityToDto(e *database.PlaylistAudio) playlistAudioDto {
	return playlistAudioDto{
		PlaylistID: e.PlaylistID,
		AudioID:    e.AudioID,
	}
}