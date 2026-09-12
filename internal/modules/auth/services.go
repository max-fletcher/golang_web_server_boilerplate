package auth

import (
	"context"
	"errors"
	"log"
	"time"

	constants "github.com/max-fletcher/golang_web_server_boilerplate/helpers/const"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/crypto"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/formatters"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
	common_errors "github.com/max-fletcher/golang_web_server_boilerplate/internal/errors"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
)

type UserService interface {
	Create(ctx context.Context, params users.CreateUserRequest) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
}

type Service interface {
	UserRegistration(ctx context.Context, params UserRegistrationInput) (AuthenticatedUser, string, error)
	UserLogin(ctx context.Context, params UserLoginRequest) (AuthenticatedUser, string, error)
}

type service struct {
	userService UserService
	events      events.Publisher
	JWTSecret   string
	JWTExpiry   time.Duration
}

func NewService(userService UserService, events events.Publisher, JWTSecret string, JWTExpiry time.Duration) *service {
	return &service{
		userService: userService,
		events:      events,
		JWTSecret:   JWTSecret,
		JWTExpiry:   JWTExpiry,
	}
}

func (service *service) UserRegistration(ctx context.Context, params UserRegistrationInput) (AuthenticatedUser, string, error) {
	// Not sure if this is needed anymore since I am checking unique constraint violation below on create
	// and throwing the exact same error
	_, err := service.userService.GetByEmail(ctx, params.Email)
	if err == nil {
		return AuthenticatedUser{}, "", users.ErrUserWithEmailAlreadyExists{
			Email: params.Email,
		}
	}
	// If error exists but doesn't match the errors that GetByEmail() sends back
	var userWithEmailNotFoundErr users.ErrUserWithEmailNotFound
	var userFetchFailedErr users.ErrUserFetchFailed
	if err != nil && !errors.As(err, &userWithEmailNotFoundErr) && !errors.As(err, &userFetchFailedErr) {
		return AuthenticatedUser{}, "", common_errors.ErrInternalServerError{
			Err: err,
		}
	}

	hashedPassword, err := crypto.HashPassword(params.Password)
	if err != nil {
		return AuthenticatedUser{}, "", common_errors.ErrHashingPassword{
			HashErr: err,
		}
	}
	params.Password = hashedPassword
	params.ConfirmPassword = hashedPassword

	// #TODO: Store Avatar Later. Avatar validation is done though.

	// 1st param: context for the request
	// 2nd param: the struct that we want to pass so it saves the underlying data in DB
	user, err := service.userService.Create(ctx, users.CreateUserRequest{
		Name:            params.Name,
		Email:           params.Email,
		Password:        params.Password,
		ConfirmPassword: params.ConfirmPassword,
	})
	if err != nil {
		pgErr := common_errors.GetPostgresError(err)
		if pgErr.Code == constants.PGUniqueViolationCode {
			return AuthenticatedUser{}, "", users.ErrUserWithEmailAlreadyExists{
				Email: params.Email,
			}
		}

		return AuthenticatedUser{}, "", users.ErrUserCreateFailed{
			CreateErr: err,
		}
	}

	err = service.events.Publish(ctx, events.QueueEventAuthRegistration,
		events.AuthRegistration{
			ID:    user.ID,
			Email: user.Email,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user registration event: %v", err)
		return AuthenticatedUser{}, "", err
	}

	authUser := AuthenticatedUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    formatters.StringPointerToNullString(params.Avatar),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	jwt, err := CreateJWT(authUser, service.JWTSecret, service.JWTExpiry)
	if err != nil {
		return AuthenticatedUser{}, "", err
	}

	return authUser, jwt, nil
}

func (service *service) UserLogin(ctx context.Context, params UserLoginRequest) (AuthenticatedUser, string, error) {
	user, err := service.userService.GetByEmail(ctx, params.Email)
	if err != nil {
		return AuthenticatedUser{}, "", users.ErrUserWithEmailNotFound{
			Email: params.Email,
		}
	}
	// If error exists but doesn't match the errors that GetByEmail() sends back
	var userWithEmailNotFoundErr users.ErrUserWithEmailNotFound
	var userFetchFailedErr users.ErrUserFetchFailed
	if err != nil && !errors.As(err, &userWithEmailNotFoundErr) && !errors.As(err, &userFetchFailedErr) {
		return AuthenticatedUser{}, "", common_errors.ErrInternalServerError{
			Err: err,
		}
	}

	err = crypto.CheckPassword(params.Password, user.Password)
	if err != nil {
		return AuthenticatedUser{}, "", ErrPasswordMismatch{
			Err: err,
		}
	}

	err = service.events.Publish(ctx, events.QueueEventAuthLogin,
		events.AuthRegistration{
			ID:    user.ID,
			Email: user.Email,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user login event: %v", err)
		return AuthenticatedUser{}, "", err
	}

	authUser := AuthenticatedUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	jwt, err := CreateJWT(authUser, service.JWTSecret, service.JWTExpiry)
	if err != nil {
		return AuthenticatedUser{}, "", err
	}

	return authUser, jwt, nil
}
