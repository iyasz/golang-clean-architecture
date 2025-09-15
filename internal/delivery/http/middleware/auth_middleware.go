package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/iyasz/golang-clean-architecture/internal/helper"
	"github.com/iyasz/golang-clean-architecture/internal/model"
	"github.com/iyasz/golang-clean-architecture/internal/usecase"
)

type contextKey string
const AuthContextKey contextKey = "authenticated_user"

func extractToken(authHeader string) string {
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	
	return authHeader
}

func NewAuth(userUserCase *usecase.UserUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				userUserCase.Log.Warn("Missing Authorization header")
				helper.ErrorResponse(w, helper.ErrUnauthorized)
				return
			}

			token := extractToken(authHeader)
			if token == "" {
				userUserCase.Log.Warn("Invalid Authorization header format")
				helper.ErrorResponse(w, helper.ErrUnauthorized)
				return
			}

			userUserCase.Log.Debugf("Processing authorization token: %s", token)

			request := &model.VerifyUserRequest{Token: token}
			auth, err := userUserCase.Verify(r.Context(), request)
			if err != nil {
				userUserCase.Log.Warnf("Failed find user by token : %+v", err)
				helper.ErrorResponse(w, helper.ErrUnauthorized)
				return
			}

			userUserCase.Log.Debugf("Authenticated user ID: %+v", auth.ID)

			ctx := context.WithValue(r.Context(), AuthContextKey, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(r *http.Request) *model.Auth {
	return r.Context().Value(AuthContextKey).(*model.Auth)
}