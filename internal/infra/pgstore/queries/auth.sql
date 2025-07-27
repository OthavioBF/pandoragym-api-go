-- name: GetUserForAuth :one
SELECT id, name, email, phone, avatar_url, password, role, created_at, updated_at
FROM users
WHERE email = $1;

-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (id, user_id, token, expires_at, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetPasswordResetToken :one
SELECT id, user_id, token, expires_at, used_at, created_at
FROM password_reset_tokens
WHERE token = $1 AND expires_at > NOW() AND used_at IS NULL;

-- name: MarkPasswordResetTokenAsUsed :exec
UPDATE password_reset_tokens
SET used_at = NOW()
WHERE token = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2, updated_at = $3
WHERE id = $1;
