package auth

import (
	"database/sql"
	"net/http"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/user"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func handleGoogleSignIn(apiCfg *shared.ApiConfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := shared.ValidateRequestBody[*googleSignInDTO](r)

		if err != nil {
			shared.ResBadRequest(w, err.Error())
			return
		}

		claims, err := ValidateGoogleJWT(body.Token)
		if err != nil {
			shared.ResUnauthorized(w, shared.ErrInvalidGoogleToken)
			return
		}

		authUser, err := user.GetUserByEmail(apiCfg.DB, r.Context(), claims.Email)
		if err != nil {
			newUser, err := user.CreateUser(apiCfg.DB, r.Context(), &database.CreateUserParams{
				Name:  sql.NullString{},
				Email: sql.NullString{String: claims.Email, Valid: true},
			})

			if err != nil {
				shared.ResInternalServerErrorDef(w)
				return
			}

			authUser = newUser
		}

		tokenString, err := shared.GenerateAccessToken(&shared.TokenClaims{
			UserId: authUser.ID,
			Email:  authUser.Email.String,
		})
		if err != nil {
			shared.ResInternalServerErrorDef(w)
			return
		}

		res := tokenPayloadDTO{
			AccessToken: tokenString,
			User:        user.UserEntityToDto(*authUser),
		}

		shared.ResOK(w, res)
	}
}
