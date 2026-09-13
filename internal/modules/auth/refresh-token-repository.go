package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/db"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, params db.CreateRefreshTokenParams) (db.RefreshToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (db.RefreshToken, error)
	Update(ctx context.Context, params db.UpdateRefreshTokenParams) (db.RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) (int64, error)
	DeleteByUserId(ctx context.Context, userId uuid.UUID) (int64, error)
	DeleteAll(ctx context.Context, userId uuid.UUID) error
}

type refreshTokenRepository struct {
	DB     *db.Queries // Reference to a DB connection. Will be used to query data form DB.
	expiry time.Duration
}

func NewRefreshTokenRepository(database *db.Queries, expiry time.Duration) *refreshTokenRepository {
	return &refreshTokenRepository{
		DB:     database,
		expiry: expiry,
	}
}

// No need to pass the entire response writer for these refreshTokenRepository functions. Just passing the context from response writer(r.Context()) from where
// this is called from (i.e service) works just fine instead of passing the response writer, then using r.Context() as first param of sqlc function

func (refreshTokenRepository *refreshTokenRepository) Create(ctx context.Context, params db.CreateRefreshTokenParams) (db.RefreshToken, error) {
	return refreshTokenRepository.DB.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID:        params.ID,
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		ExpiresAt: params.ExpiresAt,
		RevokedAt: sql.NullTime{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}

func (refreshTokenRepository *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
	return refreshTokenRepository.DB.GetRefreshTokenByTokenHash(ctx, tokenHash)
}

func (refreshTokenRepository *refreshTokenRepository) Update(ctx context.Context, params db.UpdateRefreshTokenParams) (db.RefreshToken, error) {
	return refreshTokenRepository.DB.UpdateRefreshToken(ctx, db.UpdateRefreshTokenParams{
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		ExpiresAt: params.ExpiresAt,
		RevokedAt: params.RevokedAt,
		UpdatedAt: time.Now().UTC(),
	})
}

func (refreshTokenRepository *refreshTokenRepository) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	return refreshTokenRepository.DB.DeleteRefreshToken(ctx, id)
}

func (refreshTokenRepository *refreshTokenRepository) DeleteByUserId(ctx context.Context, userId uuid.UUID) (int64, error) {
	return refreshTokenRepository.DB.DeleteRefreshTokenByUserId(ctx, userId)
}

func (refreshTokenRepository *refreshTokenRepository) DeleteAll(ctx context.Context, userId uuid.UUID) error {
	return refreshTokenRepository.DB.DeleteAllRefreshTokens(ctx)
}
