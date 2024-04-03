package auth

import (
	"context"
	"database/sql"

	"github.com/kume1a/sonifybackend/internal/database"
	"github.com/kume1a/sonifybackend/internal/modules/user"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func AuthWithEmail(apiCfg shared.ApiConfg, ctx context.Context, email string, password string) (*tokenPayloadDTO, *shared.HttpError) {
	userExistsByEmail, err := user.UserExistsByEmail(apiCfg.DB, ctx, email)
	if err != nil {
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	var authUser *database.User
	var httpErr *shared.HttpError
	if userExistsByEmail {
		authUser, httpErr = signInWithEmail(apiCfg, ctx, email, password)
	} else {
		authUser, httpErr = signUpWithEmail(apiCfg, ctx, email, password)
	}

	if httpErr != nil {
		return nil, httpErr
	}

	tokenPayload, err := getTokenPayloadDtoFromUserEntity(authUser)
	if err != nil {
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return tokenPayload, nil
}

func signInWithEmail(apiCfg shared.ApiConfg, ctx context.Context, email string, password string) (*database.User, *shared.HttpError) {
	authUser, err := user.GetUserByEmail(apiCfg.DB, ctx, email)

	if err != nil {
		return nil, shared.HttpErrUnauthorized(shared.ErrInvalidEmailOrPassword)
	}

	if authUser.AuthProvider != database.AuthProviderEMAIL {
		return nil, shared.HttpErrBadRequest(shared.ErrInvalidAuthMethod)
	}

	if !ComparePasswordHash(password, authUser.PasswordHash.String) {
		return nil, shared.HttpErrUnauthorized(shared.ErrInvalidEmailOrPassword)
	}

	return authUser, nil
}

func signUpWithEmail(apiCfg shared.ApiConfg, ctx context.Context, email string, password string) (*database.User, *shared.HttpError) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	newUser, err := user.CreateUser(apiCfg.DB, ctx, &database.CreateUserParams{
		Name:         sql.NullString{},
		Email:        sql.NullString{String: email, Valid: true},
		PasswordHash: sql.NullString{String: passwordHash, Valid: true},
		AuthProvider: database.AuthProviderEMAIL,
	})
	if err != nil {
		return nil, shared.HttpErrInternalServerErrorDef()
	}

	return newUser, nil
}
