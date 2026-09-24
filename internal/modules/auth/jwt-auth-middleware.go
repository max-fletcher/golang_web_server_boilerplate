package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	request_context "github.com/max-fletcher/golang_web_server_boilerplate/internal/request-context"
)

type JWTTokenService interface {
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
}

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

type Middleware struct {
	userService  UserService
	tokenService JWTTokenService
}

func NewMiddleware(tokenService JWTTokenService, userService UserService) *Middleware {
	return &Middleware{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (middleware *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			responses.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			responses.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		token := parts[1]

		userID, err := middleware.tokenService.ValidateAccessToken(token)
		if err != nil {
			responses.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := request_context.StoreUserIDInContext(r.Context(), userID)

		_, err = middleware.userService.GetByID(ctx, userID)
		if err != nil {
			responses.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
