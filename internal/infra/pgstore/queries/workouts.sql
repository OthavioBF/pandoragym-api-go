-- name: CreateWorkout :one
INSERT INTO workouts (id, name, description, thumbnail, video_url, rest_time_between_exercises, 
                     level, week_days, exclusive, is_template, modality, personal_id, student_id, 
                     plan_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id;

-- name: GetWorkouts :many
SELECT id, name, description, thumbnail, video_url, rest_time_between_exercises, 
       level, week_days, exclusive, is_template, modality, personal_id, student_id, 
       plan_id, created_at, updated_at
FROM workouts 
WHERE (personal_id = $1 OR student_id = $1 OR exclusive = false) 
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetWorkoutById :one
SELECT id, name, description, thumbnail, video_url, rest_time_between_exercises, 
       level, week_days, exclusive, is_template, modality, personal_id, student_id, 
       plan_id, created_at, updated_at
FROM workouts 
WHERE id = $1 AND (personal_id = $2 OR student_id = $2 OR exclusive = false) 
  AND deleted_at IS NULL;

-- name: UpdateWorkout :exec
UPDATE workouts 
SET name = $3, description = $4, thumbnail = $5, video_url = $6, 
    rest_time_between_exercises = $7, level = $8, week_days = $9, 
    exclusive = $10, is_template = $11, modality = $12, updated_at = $13
WHERE id = $1 AND (personal_id = $2 OR student_id = $2);

-- name: DeleteWorkout :exec
UPDATE workouts 
SET deleted_at = $3 
WHERE id = $1 AND (personal_id = $2 OR student_id = $2);

-- name: GetUserRole :one
SELECT role FROM users WHERE id = $1;
