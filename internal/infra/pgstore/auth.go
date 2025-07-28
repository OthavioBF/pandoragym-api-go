package pgstore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Auth-related request/response types
type AuthenticateRequest struct {
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required"`
	DeviceInfo *string `json:"deviceInfo,omitempty"`
	IPAddress  *string `json:"ipAddress,omitempty"`
}

type AuthenticateResponse struct {
	User         UserResponse `json:"user"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresIn    int64        `json:"expiresIn"`
}

type PasswordRecoverRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type CreatePasswordResetTokenParams struct {
	UserID    uuid.UUID `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

const createPasswordResetToken = `-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)`

func (q *Queries) CreatePasswordResetToken(ctx context.Context, arg CreatePasswordResetTokenParams) error {
	_, err := q.db.Exec(ctx, createPasswordResetToken, arg.UserID, arg.Token, arg.ExpiresAt)
	return err
}

const getPasswordResetToken = `-- name: GetPasswordResetToken :one
SELECT id, user_id, token, expires_at, used_at, created_at
FROM password_reset_tokens
WHERE token = $1`

func (q *Queries) GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	row := q.db.QueryRow(ctx, getPasswordResetToken, token)
	var i PasswordResetToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
		&i.UsedAt,
		&i.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const markPasswordResetTokenAsUsed = `-- name: MarkPasswordResetTokenAsUsed :exec
UPDATE password_reset_tokens SET used_at = NOW() WHERE token = $1`

func (q *Queries) MarkPasswordResetTokenAsUsed(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, markPasswordResetTokenAsUsed, token)
	return err
}
