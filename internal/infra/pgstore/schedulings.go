package pgstore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Scheduling-related request/response types
type CreateSchedulingParams struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	PersonalID uuid.UUID        `json:"personalId" db:"personal_id" validate:"required"`
	StudentID  uuid.UUID        `json:"studentId" db:"student_id" validate:"required"`
	WorkoutID  *uuid.UUID       `json:"workoutId,omitempty" db:"workout_id"`
	Date       time.Time        `json:"date" db:"date" validate:"required"`
	Type       SchedulingType   `json:"type" db:"type" validate:"required,oneof=ONLINE IN_PERSON"`
	Status     SchedulingStatus `json:"status" db:"status" validate:"required"`
	CreatedAt  time.Time        `json:"createdAt" db:"created_at"`
	UserID     *uuid.UUID       `json:"userId,omitempty" db:"user_id"`
}

type GetSchedulingByIdParams struct {
	ID         uuid.UUID `json:"id" db:"id" validate:"required"`
	PersonalID uuid.UUID `json:"personalId" db:"personal_id" validate:"required"`
}

type UpdateSchedulingStatusParams struct {
	ID     uuid.UUID        `json:"id" db:"id" validate:"required"`
	Status SchedulingStatus `json:"status" db:"status" validate:"required"`
}

type UpdateSchedulingWithStartTimeParams struct {
	ID        uuid.UUID        `json:"id" db:"id" validate:"required"`
	StartedAt *time.Time       `json:"startedAt,omitempty" db:"started_at"`
	Status    SchedulingStatus `json:"status" db:"status" validate:"required"`
}

type UpdateSchedulingWithCompletedTimeParams struct {
	ID          uuid.UUID        `json:"id" db:"id" validate:"required"`
	CompletedAt *time.Time       `json:"completedAt,omitempty" db:"completed_at"`
	Status      SchedulingStatus `json:"status" db:"status" validate:"required"`
}

type UpdateSchedulingWithCanceledTimeParams struct {
	ID     uuid.UUID        `json:"id" db:"id" validate:"required"`
	Status SchedulingStatus `json:"status" db:"status" validate:"required"`
}

type CreateSchedulingHistoryParams struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	ScheduleID uuid.UUID        `json:"scheduleId" db:"schedule_id" validate:"required"`
	UserID     uuid.UUID        `json:"userId" db:"user_id" validate:"required"`
	Status     SchedulingStatus `json:"status" db:"status" validate:"required"`
	ChangedAt  *time.Time       `json:"changedAt,omitempty" db:"changed_at"`
	ChangedBy  string           `json:"changedBy" db:"changed_by" validate:"required"`
	Reason     *string          `json:"reason,omitempty" db:"reason"`
	Notes      *string          `json:"notes,omitempty" db:"notes"`
}

type CancelSchedulingRequest struct {
	Reason string `json:"reason" validate:"required,min=5,max=500"`
}

type SchedulingResponse struct {
	ID          uuid.UUID        `json:"id"`
	PersonalID  uuid.UUID        `json:"personalId"`
	StudentID   uuid.UUID        `json:"studentId"`
	WorkoutID   *uuid.UUID       `json:"workoutId,omitempty"`
	Date        time.Time        `json:"date"`
	Type        SchedulingType   `json:"type"`
	Status      SchedulingStatus `json:"status"`
	StartedAt   *time.Time       `json:"startedAt,omitempty"`
	CompletedAt *time.Time       `json:"completedAt,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UserID      *uuid.UUID       `json:"userId,omitempty"`
}

const createScheduling = `-- name: CreateScheduling :one
INSERT INTO scheduling (
  id, personal_id, student_id, workout_id, date, type, status, created_at, user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, personal_id, student_id, workout_id, date, type, status, started_at, completed_at, stard_at, created_at, user_id`

func (q *Queries) CreateScheduling(ctx context.Context, arg CreateSchedulingParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createScheduling,
		arg.ID,
		arg.PersonalID,
		arg.StudentID,
		arg.WorkoutID,
		arg.Date,
		arg.Type,
		arg.Status,
		arg.CreatedAt,
		arg.UserID,
	)

	var i Scheduling
	err := row.Scan(
		&i.ID,
		&i.PersonalID,
		&i.StudentID,
		&i.WorkoutID,
		&i.Date,
		&i.Type,
		&i.Status,
		&i.StartedAt,
		&i.CompletedAt,
		&i.StardAt,
		&i.CreatedAt,
		&i.UserID,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return i.ID, nil
}

const getSchedulings = `-- name: GetSchedulings :many
SELECT id, personal_id, student_id, workout_id, date, type, status, started_at, completed_at, stard_at, created_at, user_id 
FROM scheduling 
ORDER BY date DESC`

func (q *Queries) GetSchedulings(ctx context.Context, personalID uuid.UUID) ([]Scheduling, error) {
	rows, err := q.db.Query(ctx, getSchedulings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Scheduling
	for rows.Next() {
		var i Scheduling
		if err := rows.Scan(
			&i.ID,
			&i.PersonalID,
			&i.StudentID,
			&i.WorkoutID,
			&i.Date,
			&i.Type,
			&i.Status,
			&i.StartedAt,
			&i.CompletedAt,
			&i.StardAt,
			&i.CreatedAt,
			&i.UserID,
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

const getSchedulingById = `-- name: GetSchedulingById :one
SELECT id, personal_id, student_id, workout_id, date, type, status, started_at, completed_at, stard_at, created_at, user_id 
FROM scheduling 
WHERE id = $1`

func (q *Queries) GetSchedulingById(ctx context.Context, arg GetSchedulingByIdParams) (*Scheduling, error) {
	row := q.db.QueryRow(ctx, getSchedulingById, arg.ID)

	var i Scheduling
	err := row.Scan(
		&i.ID,
		&i.PersonalID,
		&i.StudentID,
		&i.WorkoutID,
		&i.Date,
		&i.Type,
		&i.Status,
		&i.StartedAt,
		&i.CompletedAt,
		&i.StardAt,
		&i.CreatedAt,
		&i.UserID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

const updateSchedulingStatus = `-- name: UpdateSchedulingStatus :one
UPDATE scheduling 
SET status = $2
WHERE id = $1
RETURNING id, personal_id, student_id, workout_id, date, type, status, started_at, completed_at, stard_at, created_at, user_id`

func (q *Queries) UpdateSchedulingStatus(ctx context.Context, arg UpdateSchedulingStatusParams) error {
	_, err := q.db.Exec(ctx, updateSchedulingStatus, arg.ID, arg.Status)
	return err
}

const updateSchedulingWithStartTime = `-- name: UpdateSchedulingWithStartTime :one
UPDATE scheduling 
SET started_at = $2, status = $3
WHERE id = $1
RETURNING id, personal_id, student_id, workout_id, date, type, status, started_at, completed_at, stard_at, created_at, user_id`

func (q *Queries) UpdateSchedulingWithStartTime(ctx context.Context, arg UpdateSchedulingWithStartTimeParams) error {
	_, err := q.db.Exec(ctx, updateSchedulingWithStartTime, arg.ID, arg.StartedAt, arg.Status)
	return err
}

const updateSchedulingWithCompletedTime = `-- name: UpdateSchedulingWithCompletedTime :one
UPDATE scheduling 
SET completed_at = $2, status = $3
WHERE id = $1
RETURNING id, personal_id, student_id, workout_id, date, type, status, started_at, completed_at, stard_at, created_at, user_id`

func (q *Queries) UpdateSchedulingWithCompletedTime(ctx context.Context, arg UpdateSchedulingWithCompletedTimeParams) error {
	_, err := q.db.Exec(ctx, updateSchedulingWithCompletedTime, arg.ID, arg.CompletedAt, arg.Status)
	return err
}

func (q *Queries) UpdateSchedulingWithCanceledTime(ctx context.Context, arg UpdateSchedulingWithCanceledTimeParams) error {
	_, err := q.db.Exec(ctx, updateSchedulingStatus, arg.ID, arg.Status)
	return err
}

const createSchedulingHistory = `-- name: CreateSchedulingHistory :one
INSERT INTO schedulings_history (
  id, schedule_id, user_id, status, changed_at, changed_by, reason, notes
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, schedule_id, user_id, status, changed_at, changed_by, reason, notes`

func (q *Queries) CreateSchedulingHistory(ctx context.Context, arg CreateSchedulingHistoryParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createSchedulingHistory,
		arg.ID,
		arg.ScheduleID,
		arg.UserID,
		arg.Status,
		arg.ChangedAt,
		arg.ChangedBy,
		arg.Reason,
		arg.Notes,
	)

	var i SchedulingsHistory
	err := row.Scan(
		&i.ID,
		&i.ScheduleID,
		&i.UserID,
		&i.Status,
		&i.ChangedAt,
		&i.ChangedBy,
		&i.Reason,
		&i.Notes,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return i.ID, nil
}
