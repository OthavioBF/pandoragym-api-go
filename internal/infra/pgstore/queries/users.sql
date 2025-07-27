-- name: CreateUser :one
INSERT INTO users (id, name, email, phone, password, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: CreateStudent :exec
INSERT INTO student (id, born_date, age, weight, objective, training_frequency, did_bodybuilding, medical_condition, physical_activity_level, observations)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: CreatePersonal :exec
INSERT INTO personal (id, description, video_url, experience, specialization, qualifications)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserById :one
SELECT u.id, u.name, u.email, u.phone, u.avatar_url, u.role, u.created_at, u.updated_at,
       p.rating, p.description, p.video_url, p.experience, p.specialization, p.qualifications,
       st.born_date, st.age, st.weight, st.objective, st.training_frequency, st.did_bodybuilding,
       st.medical_condition, st.physical_activity_level, st.observations
FROM users u
LEFT JOIN personal p ON u.id = p.id
LEFT JOIN student st ON u.id = st.id
WHERE u.id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, phone, avatar_url, password, role, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetAllUsers :many
SELECT id, name, email, phone, avatar_url, role, created_at, updated_at
FROM users
ORDER BY created_at DESC;

-- name: UpdateUserProfile :exec
UPDATE users 
SET name = COALESCE($2, name), 
    phone = COALESCE($3, phone), 
    updated_at = $4
WHERE id = $1;

-- name: UpdateUserAvatar :exec
UPDATE users 
SET avatar_url = $2, updated_at = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
