package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService interface {
	GenerateJWTAccessToken(ctx context.Context, id uuid.UUID) (string, error)
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
	GenerateRefreshToken() (string, error)
}

type tokenService struct {
	jwtSecret []byte // needed for token.SignedString or else it fails
	jwtExpiry time.Duration
}

func NewTokenService(jwtSecret string, jwtExpiry time.Duration) *tokenService {
	return &tokenService{
		jwtSecret: []byte(jwtSecret), // needed for token.SignedString or else it fails
		jwtExpiry: jwtExpiry,
	}
}

// JWT Body. Alter this as you see fit.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (tokenService *tokenService) GenerateJWTAccessToken(ctx context.Context, id uuid.UUID) (string, error) {
	claims := Claims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenService.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(tokenService.jwtSecret)
	if err != nil {
		return "", ErrJWTGenerationFailed{
			Err: err,
		}
	}

	return signedToken, nil
}

func (tokenService *tokenService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			// Make sure the token uses the algorithm we expect.
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrUnknownSigningMethod{
					Err: fmt.Errorf("received signing method: %s", token.Method.Alg()),
				}
			}

			return tokenService.jwtSecret, nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, ErrJWTInvalid{
			Err: errors.New("Invalid JWT"),
		}
	}

	if claims.UserID == uuid.Nil {
		return uuid.Nil, ErrJWTInvalid{
			Err: errors.New("Invalid JWT"),
		}
	}

	return claims.UserID, nil
}

func (tokenService *tokenService) GenerateRefreshToken() (string, error) { // Doesn't provide a JWT. Instead produces a random opaque token.
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
