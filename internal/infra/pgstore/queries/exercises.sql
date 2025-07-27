-- name: CreateExercise :one
INSERT INTO exercises_templates (id, name, thumbnail, video_url, load, sets, reps, 
                                rest_time_between_sets, personal_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;

-- name: GetExercises :many
SELECT id, name, thumbnail, video_url, load, sets, reps, rest_time_between_sets, 
       personal_id, created_at, updated_at
FROM exercises_templates 
WHERE personal_id = $1 OR personal_id IS NULL
ORDER BY created_at DESC;

-- name: GetExerciseById :one
SELECT id, name, thumbnail, video_url, load, sets, reps, rest_time_between_sets, 
       personal_id, created_at, updated_at
FROM exercises_templates 
WHERE id = $1 AND (personal_id = $2 OR personal_id IS NULL);

-- name: UpdateExercise :exec
UPDATE exercises_templates 
SET name = $3, thumbnail = $4, video_url = $5, load = $6, sets = $7, 
    reps = $8, rest_time_between_sets = $9, updated_at = $10
WHERE id = $1 AND personal_id = $2;

-- name: DeleteExercise :exec
DELETE FROM exercises_templates WHERE id = $1 AND personal_id = $2;

-- name: AddExerciseToWorkout :one
INSERT INTO exercises_setup (id, name, thumbnail, video_url, sets, reps, 
                           rest_time_between_sets, load, workout_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;

-- name: RemoveExerciseFromWorkout :exec
DELETE FROM exercises_setup 
WHERE exercises_setup.id = $1 AND workout_id = $2 
  AND workout_id IN (SELECT w.id FROM workouts w WHERE w.personal_id = $3 OR w.student_id = $3);

-- name: GetExerciseForWorkout :one
SELECT id, name, thumbnail, video_url, load, sets, reps, rest_time_between_sets, 
       personal_id, created_at, updated_at
FROM exercises_templates 
WHERE id = $1;
