package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	constants "github.com/max-fletcher/golang_web_server_boilerplate/helpers/const"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/cryptography"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
	users_package "github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
)

type JWTTokenServiceForAuth interface {
	GenerateJWTAccessToken(ctx context.Context, id uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
}

type Service interface {
	UserRegistration(ctx context.Context, params UserRegistrationInput) (AuthenticatedUser, string, string, error)
	UserLogin(ctx context.Context, params UserLoginRequest) (AuthenticatedUser, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, AuthenticatedUser, error)
}

type service struct {
	tokenService       JWTTokenServiceForAuth
	events             events.Publisher
	sqlDB              *sql.DB
	DB                 *db.Queries // We are importing(using DI) db.Queries here so we can use database transactions
	refreshTokenExpiry time.Duration
}

func NewService(tokenService JWTTokenServiceForAuth, events events.Publisher, conn *sql.DB, database *db.Queries, refreshTokenExpiry time.Duration) *service {
	return &service{
		tokenService:       tokenService,
		events:             events,
		sqlDB:              conn,
		DB:                 database, // We are importing(using DI) db.Queries here so we can use database transactions
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (service *service) UserRegistration(ctx context.Context, params UserRegistrationInput) (AuthenticatedUser, string, string, error) {
	hashedPassword, err := cryptography.HashPassword(params.Password)
	if err != nil {
		return AuthenticatedUser{}, "", "", common_errors.ErrHashingPassword{
			HashErr: err,
		}
	}
	params.Password = hashedPassword

	tx, err := service.sqlDB.BeginTx(ctx, nil) // begin transaction
	if err != nil {
		return AuthenticatedUser{}, "", "", err
	}
	defer tx.Rollback()
	// IMPORTANT: txDB contains all the queries from sqlc. So initializing txDB like this and using it to make queries will cause
	// all queries to be part of one transaction
	txDB := service.DB.WithTx(tx)
	userRepository := users_package.NewRepository(txDB)
	refreshTokenRepository := NewRefreshTokenRepository(txDB, service.refreshTokenExpiry)

	// #TODO: Store Avatar Later. Avatar validation is done though.

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := userRepository.Create(ctx, db.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		pgErr := common_errors.GetPostgresError(err)
		if pgErr.Code == constants.PGUniqueViolationCode {
			return AuthenticatedUser{}, "", "", users.ErrUserWithEmailAlreadyExists{
				Email: params.Email,
			}
		}

		return AuthenticatedUser{}, "", "", users.ErrUserCreateFailed{
			CreateErr: err,
		}
	}

	authUser := AuthenticatedUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    formatters.StringPointerToNullString(params.Avatar),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	err = service.events.Publish(ctx, events.QueueEventAuthRegistration,
		events.AuthRegistration{
			ID:    user.ID,
			Email: user.Email,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user registration event: %v", err)
		return AuthenticatedUser{}, "", "", err
	}

	jwt, err := service.tokenService.GenerateJWTAccessToken(ctx, user.ID)
	if err != nil {
		return AuthenticatedUser{}, "", "", err
	}

	refreshToken, err := service.tokenService.GenerateRefreshToken()
	if err != nil {
		return AuthenticatedUser{}, "", "", err
	}

	createRefreshTokenParams := db.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: cryptography.HashString(refreshToken),
		ExpiresAt: time.Now().Add(refreshTokenRepository.expiry).UTC(),
		RevokedAt: sql.NullTime{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = refreshTokenRepository.Create(ctx, createRefreshTokenParams)
	if err != nil {
		return AuthenticatedUser{}, "", "", err
	}

	if err := tx.Commit(); err != nil {
		return AuthenticatedUser{}, "", "", err
	}

	return authUser, jwt, refreshToken, nil
}

func (service *service) UserLogin(ctx context.Context, params UserLoginRequest) (AuthenticatedUser, string, string, error) {
	tx, err := service.sqlDB.BeginTx(ctx, nil) // begin transaction
	if err != nil {
		return AuthenticatedUser{}, "", "", err
	}
	defer tx.Rollback()
	// IMPORTANT: txDB contains all the queries from sqlc. So initializing txDB like this and using it to make queries will cause
	// all queries to be part of one transaction
	txDB := service.DB.WithTx(tx)
	userRepository := users_package.NewRepository(txDB)
	refreshTokenRepository := NewRefreshTokenRepository(txDB, service.refreshTokenExpiry)

	user, err := userRepository.GetByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // check if error is of type sql.ErrNoRows
			return AuthenticatedUser{}, "", "", users.ErrUserWithEmailNotFound{
				Email: params.Email,
			}
		}

		return AuthenticatedUser{}, "", "", users.ErrUserFetchFailed{
			FetchErr: err,
		}
	}

	err = cryptography.CheckPassword(params.Password, user.Password)
	if err != nil {
		return AuthenticatedUser{}, "", "", ErrPasswordMismatch{
			Err: err,
		}
	}

	authUser := AuthenticatedUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	err = service.events.Publish(ctx, events.QueueEventAuthLogin,
		events.AuthRegistration{
			ID:    user.ID,
			Email: user.Email,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user login event: %v", err)
		return AuthenticatedUser{}, "", "", err
	}

	jwt, err := service.tokenService.GenerateJWTAccessToken(ctx, user.ID)
	if err != nil {
		return AuthenticatedUser{}, "", "", err
	}

	refreshToken, err := service.tokenService.GenerateRefreshToken()
	if err != nil {
		return AuthenticatedUser{}, "", "", err
	}

	createRefreshTokenParams := db.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: cryptography.HashString(refreshToken),
		ExpiresAt: time.Now().Add(refreshTokenRepository.expiry).UTC(),
		RevokedAt: sql.NullTime{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = refreshTokenRepository.Create(ctx, createRefreshTokenParams)
	if err != nil {
		return AuthenticatedUser{}, "", "", err
	}

	if err := tx.Commit(); err != nil {
		return AuthenticatedUser{}, "", "", err
	}

	return authUser, jwt, refreshToken, nil
}

func (service *service) RefreshToken(ctx context.Context, refreshToken string) (string, AuthenticatedUser, error) {
	tokenHash := cryptography.HashString(refreshToken)

	refreshTokenRepository := NewRefreshTokenRepository(service.DB, service.refreshTokenExpiry)
	storedToken, err := refreshTokenRepository.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return "", AuthenticatedUser{}, ErrInvalidRefreshToken{}
	}

	if storedToken.RevokedAt.Valid {
		return "", AuthenticatedUser{}, ErrInvalidRefreshToken{}
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return "", AuthenticatedUser{}, ErrInvalidRefreshToken{}
	}

	userRepository := users_package.NewRepository(service.DB)
	user, err := userRepository.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return "", AuthenticatedUser{}, ErrInvalidRefreshToken{}
	}

	authUser := AuthenticatedUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	accessToken, err := service.tokenService.GenerateJWTAccessToken(ctx, storedToken.UserID)
	if err != nil {
		return "", AuthenticatedUser{}, err
	}

	return accessToken, authUser, nil
}
