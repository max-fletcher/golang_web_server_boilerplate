package auth

import (
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
)

// JWT Body. Alter this as you see fit.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateJWT(user AuthenticatedUser, secret string, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", ErrJWTGenerationFailed{
			Err: err,
		}
	}

	return signedToken, nil
}
