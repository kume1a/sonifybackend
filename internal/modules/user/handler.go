package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handlerCreateUser(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Name string `json:"name"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			shared.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}

		createdAt := time.Now().UTC()
		user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			Name:      params.Name,
		})
		if err != nil {
			shared.RespondWithError(w, 400, fmt.Sprintf("Couldn't create user %v", err))
			return
		}

		shared.RespondWithJSON(w, 200, databaseUserToUser(user))
	}
}

func Router(apiCfg *shared.ApiConfg, router *mux.Router) *mux.Router {
	r := router.PathPrefix("/users").Subrouter()

	r.HandleFunc("", handlerCreateUser(apiCfg)).Methods("POST")

	return r
}
