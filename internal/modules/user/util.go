package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func UserEntityToDto(userEntity *database.User) shared.UserDto {
	return shared.UserDto{
		ID:        userEntity.ID,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
		Name:      userEntity.Name.String,
		Email:     userEntity.Email.String,
	}
}

func (dto updateUserDTO) Validate() error {
	_, err := govalidator.ValidateStruct(dto)
	return err
}
