package pgstore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Exercise-related request/response types
type CreateExerciseParams struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	Name                string     `json:"name" db:"name" validate:"required,min=2,max=100"`
	Thumbnail           string     `json:"thumbnail" db:"thumbnail" validate:"required"`
	VideoURL            string     `json:"videoUrl" db:"video_url" validate:"required,url"`
	Load                *int32     `json:"load,omitempty" db:"load"`
	Sets                int32      `json:"sets" db:"sets" validate:"required,min=1"`
	Reps                int32      `json:"reps" db:"reps" validate:"required,min=1"`
	RestTimeBetweenSets *int32     `json:"restTimeBetweenSets,omitempty" db:"rest_time_between_sets"`
	PersonalID          *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
	CreatedAt           time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt" db:"updated_at"`
}

type UpdateExerciseParams struct {
	ID                  uuid.UUID  `json:"id" db:"id" validate:"required"`
	Name                *string    `json:"name,omitempty" db:"name" validate:"omitempty,min=2,max=100"`
	Thumbnail           *string    `json:"thumbnail,omitempty" db:"thumbnail"`
	VideoURL            *string    `json:"videoUrl,omitempty" db:"video_url" validate:"omitempty,url"`
	Load                *int32     `json:"load,omitempty" db:"load"`
	Sets                *int32     `json:"sets,omitempty" db:"sets" validate:"omitempty,min=1"`
	Reps                *int32     `json:"reps,omitempty" db:"reps" validate:"omitempty,min=1"`
	RestTimeBetweenSets *int32     `json:"restTimeBetweenSets,omitempty" db:"rest_time_between_sets"`
	PersonalID          *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
	UpdatedAt           time.Time  `json:"updatedAt" db:"updated_at"`
}

type GetExerciseByIdParams struct {
	ID         uuid.UUID  `json:"id" db:"id" validate:"required"`
	PersonalID *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
}

type DeleteExerciseParams struct {
	ID         uuid.UUID  `json:"id" db:"id" validate:"required"`
	PersonalID *uuid.UUID `json:"personalId,omitempty" db:"personal_id"`
}

type AddExerciseToWorkoutParams struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	Name                string    `json:"name" db:"name"`
	Thumbnail           string    `json:"thumbnail" db:"thumbnail"`
	VideoURL            string    `json:"videoUrl" db:"video_url"`
	Sets                int32     `json:"sets" db:"sets"`
	Reps                int32     `json:"reps" db:"reps"`
	RestTimeBetweenSets int32     `json:"restTimeBetweenSets" db:"rest_time_between_sets"`
	Load                int32     `json:"load" db:"load"`
	WorkoutID           uuid.UUID `json:"workoutId" db:"workout_id"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
}

type RemoveExerciseFromWorkoutParams struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"required"`
	WorkoutID uuid.UUID `json:"workoutId" db:"workout_id" validate:"required"`
}

type ExerciseResponse struct {
	ID                  uuid.UUID  `json:"id"`
	Name                string     `json:"name"`
	Thumbnail           string     `json:"thumbnail"`
	VideoURL            string     `json:"videoUrl"`
	Load                *int32     `json:"load,omitempty"`
	Sets                int32      `json:"sets"`
	Reps                int32      `json:"reps"`
	RestTimeBetweenSets *int32     `json:"restTimeBetweenSets,omitempty"`
	PersonalID          *uuid.UUID `json:"personalId,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

const createExercise = `-- name: CreateExercise :one
INSERT INTO exercises_template (
  id, name, thumbnail, video_url, load, sets, reps, rest_time_between_sets, personal_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id, name, thumbnail, video_url, load, sets, reps, rest_time_between_sets, personal_id, created_at, updated_at`

func (q *Queries) CreateExercise(ctx context.Context, arg CreateExerciseParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createExercise,
		arg.ID,
		arg.Name,
		arg.Thumbnail,
		arg.VideoURL,
		arg.Load,
		arg.Sets,
		arg.Reps,
		arg.RestTimeBetweenSets,
		arg.PersonalID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)

	var i ExercisesTemplate
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Thumbnail,
		&i.VideoURL,
		&i.Load,
		&i.Sets,
		&i.Reps,
		&i.RestTimeBetweenSets,
		&i.PersonalID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return i.ID, nil
}

const getExercises = `-- name: GetExercises :many
SELECT id, name, thumbnail, video_url, load, sets, reps, rest_time_between_sets, personal_id, created_at, updated_at 
FROM exercises_template 
ORDER BY created_at DESC`

func (q *Queries) GetExercises(ctx context.Context, personalID *uuid.UUID) ([]ExercisesTemplate, error) {
	rows, err := q.db.Query(ctx, getExercises)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ExercisesTemplate
	for rows.Next() {
		var i ExercisesTemplate
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Thumbnail,
			&i.VideoURL,
			&i.Load,
			&i.Sets,
			&i.Reps,
			&i.RestTimeBetweenSets,
			&i.PersonalID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getExerciseById = `-- name: GetExerciseById :one
SELECT id, name, thumbnail, video_url, load, sets, reps, rest_time_between_sets, personal_id, created_at, updated_at 
FROM exercises_template 
WHERE id = $1`

func (q *Queries) GetExerciseById(ctx context.Context, arg GetExerciseByIdParams) (*ExercisesTemplate, error) {
	row := q.db.QueryRow(ctx, getExerciseById, arg.ID)

	var i ExercisesTemplate
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Thumbnail,
		&i.VideoURL,
		&i.Load,
		&i.Sets,
		&i.Reps,
		&i.RestTimeBetweenSets,
		&i.PersonalID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const updateExercise = `-- name: UpdateExercise :one
UPDATE exercises_template 
SET 
  name = COALESCE($2, name),
  thumbnail = COALESCE($3, thumbnail),
  video_url = COALESCE($4, video_url),
  load = COALESCE($5, load),
  sets = COALESCE($6, sets),
  reps = COALESCE($7, reps),
  rest_time_between_sets = COALESCE($8, rest_time_between_sets),
  updated_at = $9
WHERE id = $1
RETURNING id, name, thumbnail, video_url, load, sets, reps, rest_time_between_sets, personal_id, created_at, updated_at`

func (q *Queries) GetExerciseForWorkout(ctx context.Context, id uuid.UUID) (*ExercisesTemplate, error) {
	row := q.db.QueryRow(ctx, getExerciseById, id)

	var i ExercisesTemplate
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Thumbnail,
		&i.VideoURL,
		&i.Load,
		&i.Sets,
		&i.Reps,
		&i.RestTimeBetweenSets,
		&i.PersonalID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

func (q *Queries) UpdateExercise(ctx context.Context, arg UpdateExerciseParams) error {
	_, err := q.db.Exec(ctx, updateExercise,
		arg.ID,
		arg.Name,
		arg.Thumbnail,
		arg.VideoURL,
		arg.Load,
		arg.Sets,
		arg.Reps,
		arg.RestTimeBetweenSets,
		arg.UpdatedAt,
	)
	return err
}

const deleteExercise = `-- name: DeleteExercise :exec
DELETE FROM exercises_template 
WHERE id = $1`

func (q *Queries) DeleteExercise(ctx context.Context, arg DeleteExerciseParams) error {
	_, err := q.db.Exec(ctx, deleteExercise, arg.ID)
	return err
}

const addExerciseToWorkout = `-- name: AddExerciseToWorkout :one
INSERT INTO exercises_setup (
  id, name, thumbnail, video_url, sets, reps, rest_time_between_sets, load, workout_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id, name, thumbnail, video_url, sets, reps, rest_time_between_sets, load, workout_id, created_at, updated_at`

func (q *Queries) AddExerciseToWorkout(ctx context.Context, arg AddExerciseToWorkoutParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, addExerciseToWorkout,
		arg.ID,
		arg.Name,
		arg.Thumbnail,
		arg.VideoURL,
		arg.Sets,
		arg.Reps,
		arg.RestTimeBetweenSets,
		arg.Load,
		arg.WorkoutID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)

	var i ExercisesSetup
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Thumbnail,
		&i.VideoURL,
		&i.Sets,
		&i.Reps,
		&i.RestTimeBetweenSets,
		&i.Load,
		&i.WorkoutID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return i.ID, nil
}

const removeExerciseFromWorkout = `-- name: RemoveExerciseFromWorkout :exec
DELETE FROM exercises_setup 
WHERE id = $1 AND workout_id = $2`

func (q *Queries) RemoveExerciseFromWorkout(ctx context.Context, arg RemoveExerciseFromWorkoutParams) error {
	_, err := q.db.Exec(ctx, removeExerciseFromWorkout, arg.ID, arg.WorkoutID)
	return err
}

const getExercisesByWorkoutId = `-- name: GetExercisesByWorkoutId :many
SELECT id, name, thumbnail, video_url, sets, reps, rest_time_between_sets, load, workout_id, created_at, updated_at 
FROM exercises_setup 
WHERE workout_id = $1
ORDER BY created_at ASC`

func (q *Queries) GetExercisesByWorkoutId(ctx context.Context, workoutID uuid.UUID) ([]ExercisesSetup, error) {
	rows, err := q.db.Query(ctx, getExercisesByWorkoutId, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ExercisesSetup
	for rows.Next() {
		var i ExercisesSetup
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Thumbnail,
			&i.VideoURL,
			&i.Sets,
			&i.Reps,
			&i.RestTimeBetweenSets,
			&i.Load,
			&i.WorkoutID,
			&i.CreatedAt,
			&i.UpdatedAt,
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
