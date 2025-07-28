package pgstore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Refresh Token types
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type CreateRefreshTokenRequest struct {
	UserID     uuid.UUID `json:"userId" validate:"required"`
	DeviceInfo *string   `json:"deviceInfo,omitempty"`
	IPAddress  *string   `json:"ipAddress,omitempty"`
}

type CreateRefreshTokenParams struct {
	UserID     uuid.UUID
	Token      string
	ExpiresAt  time.Time
	DeviceInfo *string
	IPAddress  *string
}

type GetRefreshTokensByUserIDRow struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Token      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	RevokedAt  *time.Time
	DeviceInfo *string
	IPAddress  *string
}

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token,
    expires_at,
    device_info,
    ip_address
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, user_id, token, expires_at, created_at, updated_at, revoked_at, device_info, ip_address`

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (*RefreshToken, error) {
	row := q.db.QueryRow(ctx, createRefreshToken,
		arg.UserID,
		arg.Token,
		arg.ExpiresAt,
		arg.DeviceInfo,
		arg.IPAddress,
	)

	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RevokedAt,
		&i.DeviceInfo,
		&i.IPAddress,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

const getRefreshToken = `-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at, created_at, updated_at, revoked_at, device_info, ip_address 
FROM refresh_tokens 
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL`

func (q *Queries) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	row := q.db.QueryRow(ctx, getRefreshToken, token)

	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RevokedAt,
		&i.DeviceInfo,
		&i.IPAddress,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const getRefreshTokensByUserID = `-- name: GetRefreshTokensByUserID :many
SELECT id, user_id, token, expires_at, created_at, updated_at, revoked_at, device_info, ip_address 
FROM refresh_tokens 
WHERE user_id = $1 AND expires_at > NOW() AND revoked_at IS NULL
ORDER BY created_at DESC`

func (q *Queries) GetRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) ([]GetRefreshTokensByUserIDRow, error) {
	rows, err := q.db.Query(ctx, getRefreshTokensByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []GetRefreshTokensByUserIDRow
	for rows.Next() {
		var i GetRefreshTokensByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Token,
			&i.ExpiresAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.RevokedAt,
			&i.DeviceInfo,
			&i.IPAddress,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, revokeRefreshToken, token)
	return err
}

const revokeAllUserRefreshTokens = `-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE user_id = $1 AND revoked_at IS NULL`

func (q *Queries) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, revokeAllUserRefreshTokens, userID)
	return err
}

const cleanupExpiredRefreshTokens = `-- name: CleanupExpiredRefreshTokens :exec
DELETE FROM refresh_tokens 
WHERE expires_at < NOW() OR revoked_at < NOW() - INTERVAL '30 days'`

func (q *Queries) CleanupExpiredRefreshTokens(ctx context.Context) error {
	_, err := q.db.Exec(ctx, cleanupExpiredRefreshTokens)
	return err
}

const getRefreshTokenByID = `-- name: GetRefreshTokenByID :one
SELECT id, user_id, token, expires_at, created_at, updated_at, revoked_at, device_info, ip_address 
FROM refresh_tokens 
WHERE id = $1`

func (q *Queries) GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (*RefreshToken, error) {
	row := q.db.QueryRow(ctx, getRefreshTokenByID, id)

	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RevokedAt,
		&i.DeviceInfo,
		&i.IPAddress,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const updateRefreshTokenLastUsed = `-- name: UpdateRefreshTokenLastUsed :exec
UPDATE refresh_tokens 
SET updated_at = NOW()
WHERE token = $1`

func (q *Queries) UpdateRefreshTokenLastUsed(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, updateRefreshTokenLastUsed, token)
	return err
}
