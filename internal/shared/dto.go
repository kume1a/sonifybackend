package shared

import (
	"time"

	"github.com/google/uuid"
)

type HttpErrorDto struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UrlDto struct {
	Url string `json:"url"`
}

type KeywordDto struct {
	Keyword []string `json:"keyword" valid:"required"`
}

type UserDto struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

type LastCreatedAtPageParamsDto struct {
	LastCreatedAt time.Time `json:"lastCreatedAt"`
	Limit         int32     `json:"limit" valid:"required,max(200)"`
}
