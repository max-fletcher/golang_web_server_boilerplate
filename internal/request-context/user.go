package request_context

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

func StoreUserIDInContext(reqCtx context.Context, userID uuid.UUID) context.Context {
	ctx := context.WithValue(reqCtx, UserIDContextKey, userID)
	return ctx
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(uuid.UUID)
	return userID, ok
}
