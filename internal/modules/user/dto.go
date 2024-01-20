package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/kume1a/sonifybackend/internal/database"
)

type UserDto struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
}

func UserEntityToDto(userEntity database.User) UserDto {
	return UserDto{
		ID:        userEntity.ID,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
		Name:      userEntity.Name,
	}
}
