-- name: CreateScheduling :one
INSERT INTO schedulings (id, personal_id, student_id, workout_id, date, type, status, created_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: GetSchedulings :many
SELECT id, personal_id, student_id, workout_id, date, type, status, 
       started_at, completed_at, stard_at, created_at, user_id
FROM schedulings 
WHERE personal_id = $1 OR student_id = $1
ORDER BY date DESC;

-- name: GetSchedulingById :one
SELECT id, personal_id, student_id, workout_id, date, type, status, 
       started_at, completed_at, stard_at, created_at, user_id
FROM schedulings 
WHERE id = $1 AND (personal_id = $2 OR student_id = $2);

-- name: UpdateSchedulingStatus :exec
UPDATE schedulings 
SET status = $3 
WHERE id = $1 AND (personal_id = $2 OR student_id = $2);

-- name: UpdateSchedulingWithStartTime :exec
UPDATE schedulings 
SET status = $3, started_at = $4 
WHERE id = $1 AND (personal_id = $2 OR student_id = $2);

-- name: UpdateSchedulingWithCompletedTime :exec
UPDATE schedulings 
SET status = $3, completed_at = $4 
WHERE id = $1 AND (personal_id = $2 OR student_id = $2);

-- name: UpdateSchedulingWithCanceledTime :exec
UPDATE schedulings 
SET status = $3, stard_at = $4 
WHERE id = $1 AND (personal_id = $2 OR student_id = $2);

-- name: CreateSchedulingHistory :one
INSERT INTO schedulings_history (id, schedule_id, user_id, status, changed_at, changed_by, reason)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
