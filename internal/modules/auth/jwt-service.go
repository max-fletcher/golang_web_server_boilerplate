package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateJWTAccessToken(ctx context.Context, id uuid.UUID) (string, error)
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
}

type jwtService struct {
	JWTSecret []byte // needed for token.SignedString or else it fails
	JWTExpiry time.Duration
}

func NewJWTService(JWTSecret string, JWTExpiry time.Duration) *jwtService {
	return &jwtService{
		JWTSecret: []byte(JWTSecret), // needed for token.SignedString or else it fails
		JWTExpiry: JWTExpiry,
	}
}

// JWT Body. Alter this as you see fit.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (jwtService *jwtService) GenerateJWTAccessToken(ctx context.Context, id uuid.UUID) (string, error) {
	claims := Claims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtService.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtService.JWTSecret)
	if err != nil {
		return "", ErrJWTGenerationFailed{
			Err: err,
		}
	}

	return signedToken, nil
}

func (jwtService *jwtService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
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

			return jwtService.JWTSecret, nil
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
