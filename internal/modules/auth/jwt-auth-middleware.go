package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
)

type contextKey string

const userIDContextKey contextKey = "userID"

type JWTTokenServiceForMiddleware interface {
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
}

type Middleware struct {
	tokenService JWTTokenServiceForMiddleware
}

func NewMiddleware(tokenService JWTTokenServiceForMiddleware) *Middleware {
	return &Middleware{
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

		ctx := context.WithValue(
			r.Context(),
			userIDContextKey,
			userID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return userID, ok
}
