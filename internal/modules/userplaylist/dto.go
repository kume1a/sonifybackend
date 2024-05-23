package userplaylist

import "github.com/google/uuid"

type getMyPlaylistsDTO struct {
	IDs uuid.UUIDs `json:"ids" valid:"-"`
}
