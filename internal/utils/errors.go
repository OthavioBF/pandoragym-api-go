package utils

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
)

// Context keys
type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
	UserKey     contextKey = "user"
)

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrUnauthorized
	}
	return userID, nil
}
