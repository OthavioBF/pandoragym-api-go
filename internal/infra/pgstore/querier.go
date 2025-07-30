package pgstore

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	// User operations
	CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
	CreateStudent(ctx context.Context, arg CreateStudentParams) error
	CreatePersonal(ctx context.Context, arg CreatePersonalParams) error
	GetUserById(ctx context.Context, arg GetUserByIdParams) (*GetUserByIdRow, error)
	GetUserByEmail(ctx context.Context, email string) (*GetUserByEmailRow, error)
	GetUserForAuth(ctx context.Context, email string) (*GetUserForAuthRow, error)
	GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error)
	UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) error
	UpdateUserAvatar(ctx context.Context, arg UpdateUserAvatarParams) error
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserRole(ctx context.Context, id uuid.UUID) (Role, error)

	// Auth operations
	CreatePasswordResetToken(ctx context.Context, arg CreatePasswordResetTokenParams) error
	GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error)
	MarkPasswordResetTokenAsUsed(ctx context.Context, token string) error

	// Workout operations
	CreateWorkout(ctx context.Context, arg CreateWorkoutParams) (uuid.UUID, error)
	GetWorkouts(ctx context.Context, personalID *uuid.UUID) ([]GetWorkoutsRow, error)
	GetWorkoutById(ctx context.Context, arg GetWorkoutByIdParams) (*GetWorkoutByIdRow, error)
	UpdateWorkout(ctx context.Context, arg UpdateWorkoutParams) error
	DeleteWorkout(ctx context.Context, arg DeleteWorkoutParams) error

	// Exercise operations
	CreateExercise(ctx context.Context, arg CreateExerciseParams) (uuid.UUID, error)
	GetExercises(ctx context.Context, personalID *uuid.UUID) ([]ExercisesTemplate, error)
	GetExerciseById(ctx context.Context, arg GetExerciseByIdParams) (*ExercisesTemplate, error)
	GetExerciseForWorkout(ctx context.Context, id uuid.UUID) (*ExercisesTemplate, error)
	UpdateExercise(ctx context.Context, arg UpdateExerciseParams) error
	DeleteExercise(ctx context.Context, arg DeleteExerciseParams) error
	AddExerciseToWorkout(ctx context.Context, arg AddExerciseToWorkoutParams) (uuid.UUID, error)
	RemoveExerciseFromWorkout(ctx context.Context, arg RemoveExerciseFromWorkoutParams) error

	// Scheduling operations
	CreateScheduling(ctx context.Context, arg CreateSchedulingParams) (uuid.UUID, error)
	GetSchedulings(ctx context.Context, personalID uuid.UUID) ([]Scheduling, error)
	GetSchedulingById(ctx context.Context, arg GetSchedulingByIdParams) (*Scheduling, error)
	UpdateSchedulingStatus(ctx context.Context, arg UpdateSchedulingStatusParams) error
	UpdateSchedulingWithStartTime(ctx context.Context, arg UpdateSchedulingWithStartTimeParams) error
	UpdateSchedulingWithCompletedTime(ctx context.Context, arg UpdateSchedulingWithCompletedTimeParams) error
	UpdateSchedulingWithCanceledTime(ctx context.Context, arg UpdateSchedulingWithCanceledTimeParams) error
	CreateSchedulingHistory(ctx context.Context, arg CreateSchedulingHistoryParams) (uuid.UUID, error)

	// Transaction support
	WithTx(tx any) *Queries
}

var _ Querier = (*Queries)(nil)
