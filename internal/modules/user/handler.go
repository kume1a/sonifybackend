package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

		user, err := apiCfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
			Name: sql.NullString{String: body.Name, Valid: body.Name != ""},
			ID:   tokenPayload.UserId,
		})

		if err != nil {
			log.Println(err)
			shared.ResInternalServerErrorDef(w)
			return
		}

		shared.ResCreated(w, UserEntityToDto(user))
	}
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/users").Subrouter()

	r.HandleFunc("", shared.AuthMW(handleUpdateUser(apiCfg))).Methods("POST")

	return r
}
