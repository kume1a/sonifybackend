package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleUpdateUser(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		body, err := shared.ValidateRequestBody[*updateUserDTO](r)
		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		user, err := UpdateUser(apiCfg.DB, r.Context(), &database.UpdateUserParams{
			Name: sql.NullString{String: body.Name, Valid: body.Name != ""},
			ID:   tokenPayload.UserId,
		})

		if err != nil {
			log.Println(err)
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResOK(w, UserEntityToDto(user))
	}
}

func handleGetAuthUser(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenPayload, err := shared.GetAuthPayload(r)
		if err != nil {
			shared.ResUnauthorized(w, err.Error())
			return
		}

		user, err := GetUserByID(apiCfg.DB, r.Context(), tokenPayload.UserId)
		if err != nil {
			shared.ResNotFound(w, shared.ErrUserNotFound)
			return
		}

		shared.ResOK(w, UserEntityToDto(user))
	}
}
