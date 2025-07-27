-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token,
    expires_at,
    device_info,
    ip_address
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens 
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL;

-- name: GetRefreshTokensByUserID :many
SELECT * FROM refresh_tokens 
WHERE user_id = $1 AND expires_at > NOW() AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: CleanupExpiredRefreshTokens :exec
DELETE FROM refresh_tokens 
WHERE expires_at < NOW() OR revoked_at < NOW() - INTERVAL '30 days';

-- name: GetRefreshTokenByID :one
SELECT * FROM refresh_tokens 
WHERE id = $1;

-- name: UpdateRefreshTokenLastUsed :exec
UPDATE refresh_tokens 
SET updated_at = NOW()
WHERE token = $1;
