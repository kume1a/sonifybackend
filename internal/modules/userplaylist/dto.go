package userplaylist

import (
	"github.com/google/uuid"
)

type createUserPlaylistDTO struct {
	Name string `json:"name" valid:"required"`
}

type updateUserPlaylistDTO struct {
	Name string `json:"name" valid:"optional"`
}

type PlaylistIDsDTO struct {
	PlaylistIDs []uuid.UUID `json:"playlistIds" valid:"-"`
}
