package pgstore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

// Workout-related request/response types
type CreateWorkoutParams struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	Name                     string     `json:"name" db:"name" validate:"required,min=2,max=100"`
	Description              *string    `json:"description,omitempty" db:"description"`
	Thumbnail                string     `json:"thumbnail" db:"thumbnail" validate:"required"`
	VideoURL                 *string    `json:"video_url,omitempty" db:"video_url"`
	RestTimeBetweenExercises *int32     `json:"rest_time_between_exercises,omitempty" db:"rest_time_between_exercises"`
	Level                    *Level     `json:"level,omitempty" db:"level"`
	WeekDays                 []Day      `json:"week_days" db:"week_days" validate:"required"`
	Exclusive                bool       `json:"exclusive" db:"exclusive"`
	IsTemplate               bool       `json:"is_template" db:"is_template"`
	Modality                 string     `json:"modality" db:"modality" validate:"required"`
	PersonalID               *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
	StudentID                *uuid.UUID `json:"student_id,omitempty" db:"student_id"`
	PlanID                   *string    `json:"plan_id,omitempty" db:"plan_id"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
}

type UpdateWorkoutParams struct {
	ID                       uuid.UUID  `json:"id" db:"id" validate:"required"`
	Name                     *string    `json:"name,omitempty" db:"name" validate:"omitempty,min=2,max=100"`
	Description              *string    `json:"description,omitempty" db:"description"`
	Thumbnail                *string    `json:"thumbnail,omitempty" db:"thumbnail"`
	VideoURL                 *string    `json:"video_url,omitempty" db:"video_url"`
	RestTimeBetweenExercises *int32     `json:"rest_time_between_exercises,omitempty" db:"rest_time_between_exercises"`
	Level                    *Level     `json:"level,omitempty" db:"level"`
	WeekDays                 []Day      `json:"week_days,omitempty" db:"week_days"`
	Modality                 *string    `json:"modality,omitempty" db:"modality"`
	PersonalID               *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
}

type GetWorkoutByIdParams struct {
	ID         uuid.UUID  `json:"id" db:"id" validate:"required"`
	PersonalID *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
}

type DeleteWorkoutParams struct {
	ID         uuid.UUID  `json:"id" db:"id" validate:"required"`
	PersonalID *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
}

type GetWorkoutsRow struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	Name                     string     `json:"name" db:"name"`
	Description              *string    `json:"description,omitempty" db:"description"`
	Thumbnail                string     `json:"thumbnail" db:"thumbnail"`
	VideoURL                 *string    `json:"video_url,omitempty" db:"video_url"`
	RestTimeBetweenExercises *int32     `json:"rest_time_between_exercises,omitempty" db:"rest_time_between_exercises"`
	Level                    *Level     `json:"level,omitempty" db:"level"`
	WeekDays                 []Day      `json:"week_days" db:"week_days"`
	Exclusive                bool       `json:"exclusive" db:"exclusive"`
	IsTemplate               bool       `json:"is_template" db:"is_template"`
	Modality                 string     `json:"modality" db:"modality"`
	PersonalID               *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
	StudentID                *uuid.UUID `json:"student_id,omitempty" db:"student_id"`
	PlanID                   *string    `json:"plan_id,omitempty" db:"plan_id"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
}

type GetWorkoutByIdRow struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	Name                     string     `json:"name" db:"name"`
	Description              *string    `json:"description,omitempty" db:"description"`
	Thumbnail                string     `json:"thumbnail" db:"thumbnail"`
	VideoURL                 *string    `json:"video_url,omitempty" db:"video_url"`
	RestTimeBetweenExercises *int32     `json:"rest_time_between_exercises,omitempty" db:"rest_time_between_exercises"`
	Level                    *Level     `json:"level,omitempty" db:"level"`
	WeekDays                 []Day      `json:"week_days" db:"week_days"`
	Exclusive                bool       `json:"exclusive" db:"exclusive"`
	IsTemplate               bool       `json:"is_template" db:"is_template"`
	Modality                 string     `json:"modality" db:"modality"`
	PersonalID               *uuid.UUID `json:"personal_id,omitempty" db:"personal_id"`
	StudentID                *uuid.UUID `json:"student_id,omitempty" db:"student_id"`
	PlanID                   *string    `json:"plan_id,omitempty" db:"plan_id"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
}

type WorkoutResponse struct {
	ID                       uuid.UUID  `json:"id"`
	Name                     string     `json:"name"`
	Description              *string    `json:"description,omitempty"`
	Thumbnail                string     `json:"thumbnail"`
	VideoURL                 *string    `json:"video_url,omitempty"`
	RestTimeBetweenExercises *int32     `json:"rest_time_between_exercises,omitempty"`
	Level                    *Level     `json:"level,omitempty"`
	WeekDays                 []Day      `json:"week_days"`
	Exclusive                bool       `json:"exclusive"`
	IsTemplate               bool       `json:"is_template"`
	Modality                 string     `json:"modality"`
	PersonalID               *uuid.UUID `json:"personal_id,omitempty"`
	StudentID                *uuid.UUID `json:"student_id,omitempty"`
	PlanID                   *string    `json:"plan_id,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type AddExerciseToWorkoutRequest struct {
	WorkoutID  uuid.UUID `json:"workout_id" validate:"required"`
	ExerciseID uuid.UUID `json:"exercise_id" validate:"required"`
	Sets       int32     `json:"sets" validate:"required,min=1"`
	Reps       int32     `json:"reps" validate:"required,min=1"`
	Load       *int32    `json:"load,omitempty"`
	RestTime   *int32    `json:"rest_time,omitempty"`
}

const createWorkout = `-- name: CreateWorkout :one
INSERT INTO workout (id, name, description, thumbnail, video_url, rest_time_between_exercises, 
                     level, week_days, exclusive, is_template, modality, personal_id, student_id, 
                     plan_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id`

func (q *Queries) CreateWorkout(ctx context.Context, arg CreateWorkoutParams) (uuid.UUID, error) {
	// Convert []Day to []string for pq.Array
	weekDaysStr := make([]string, len(arg.WeekDays))
	for i, day := range arg.WeekDays {
		weekDaysStr[i] = string(day)
	}

	row := q.db.QueryRow(ctx, createWorkout,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Thumbnail,
		arg.VideoURL,
		arg.RestTimeBetweenExercises,
		arg.Level,
		pq.Array(weekDaysStr),
		arg.Exclusive,
		arg.IsTemplate,
		arg.Modality,
		arg.PersonalID,
		arg.StudentID,
		arg.PlanID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getWorkouts = `-- name: GetWorkouts :many
SELECT id, name, description, thumbnail, video_url, rest_time_between_exercises,
       level, week_days, exclusive, is_template, modality, personal_id, student_id,
       plan_id, created_at, updated_at
FROM workouts
WHERE ($1::uuid IS NULL OR personal_id = $1) AND deleted_at IS NULL
ORDER BY created_at DESC`

func (q *Queries) GetWorkouts(ctx context.Context, personalID *uuid.UUID) ([]GetWorkoutsRow, error) {
	rows, err := q.db.Query(ctx, getWorkouts, personalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkoutsRow
	for rows.Next() {
		var i GetWorkoutsRow
		var weekDaysStr pq.StringArray
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Thumbnail,
			&i.VideoURL,
			&i.RestTimeBetweenExercises,
			&i.Level,
			&weekDaysStr,
			&i.Exclusive,
			&i.IsTemplate,
			&i.Modality,
			&i.PersonalID,
			&i.StudentID,
			&i.PlanID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		// Convert []string to []Day
		i.WeekDays = make([]Day, len(weekDaysStr))
		for j, dayStr := range weekDaysStr {
			i.WeekDays[j] = Day(dayStr)
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkoutById = `-- name: GetWorkoutById :one
SELECT id, name, description, thumbnail, video_url, rest_time_between_exercises,
       level, week_days, exclusive, is_template, modality, personal_id, student_id,
       plan_id, created_at, updated_at
FROM workouts
WHERE id = $1 AND ($2::uuid IS NULL OR personal_id = $2) AND deleted_at IS NULL`

func (q *Queries) GetWorkoutById(ctx context.Context, arg GetWorkoutByIdParams) (*GetWorkoutByIdRow, error) {
	row := q.db.QueryRow(ctx, getWorkoutById, arg.ID, arg.PersonalID)
	var i GetWorkoutByIdRow
	var weekDaysStr pq.StringArray
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Thumbnail,
		&i.VideoURL,
		&i.RestTimeBetweenExercises,
		&i.Level,
		&weekDaysStr,
		&i.Exclusive,
		&i.IsTemplate,
		&i.Modality,
		&i.PersonalID,
		&i.StudentID,
		&i.PlanID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	// Convert []string to []Day
	i.WeekDays = make([]Day, len(weekDaysStr))
	for j, dayStr := range weekDaysStr {
		i.WeekDays[j] = Day(dayStr)
	}
	return &i, nil
}

const updateWorkout = `-- name: UpdateWorkout :exec
UPDATE workouts 
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    thumbnail = COALESCE($4, thumbnail),
    video_url = COALESCE($5, video_url),
    rest_time_between_exercises = COALESCE($6, rest_time_between_exercises),
    level = COALESCE($7, level),
    week_days = COALESCE($8, week_days),
    modality = COALESCE($9, modality),
    updated_at = $10
WHERE id = $1 AND ($11::uuid IS NULL OR personal_id = $11)`

func (q *Queries) UpdateWorkout(ctx context.Context, arg UpdateWorkoutParams) error {
	// Convert []Day to []string for pq.Array
	var weekDaysArray interface{}
	if arg.WeekDays != nil {
		weekDaysStr := make([]string, len(arg.WeekDays))
		for i, day := range arg.WeekDays {
			weekDaysStr[i] = string(day)
		}
		weekDaysArray = pq.Array(weekDaysStr)
	}

	_, err := q.db.Exec(ctx, updateWorkout,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Thumbnail,
		arg.VideoURL,
		arg.RestTimeBetweenExercises,
		arg.Level,
		weekDaysArray,
		arg.Modality,
		arg.UpdatedAt,
		arg.PersonalID,
	)
	return err
}

const deleteWorkout = `-- name: DeleteWorkout :exec
UPDATE workouts 
SET deleted_at = NOW() 
WHERE id = $1 AND ($2::uuid IS NULL OR personal_id = $2)`

func (q *Queries) DeleteWorkout(ctx context.Context, arg DeleteWorkoutParams) error {
	_, err := q.db.Exec(ctx, deleteWorkout, arg.ID, arg.PersonalID)
	return err
}
