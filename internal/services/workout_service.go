package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

type WorkoutService struct {
	queries *pgstore.Queries
}

func NewWorkoutService(queries *pgstore.Queries) *WorkoutService {
	return &WorkoutService{
		queries: queries,
	}
}

func (s *WorkoutService) GetWorkouts(ctx context.Context, userID uuid.UUID) ([]pgstore.Workout, error) {
	// Get user role to determine access
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	dbWorkouts, err := s.queries.GetWorkouts(ctx, personalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workouts: %w", err)
	}

	workouts := make([]pgstore.Workout, len(dbWorkouts))
	for i, dbWorkout := range dbWorkouts {
		workouts[i] = s.convertWorkoutFromDB(dbWorkout)
	}

	return workouts, nil
}

func (s *WorkoutService) GetWorkoutByID(ctx context.Context, workoutID, userID uuid.UUID) (*pgstore.Workout, error) {
	// Get user role to determine access
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	dbWorkout, err := s.queries.GetWorkoutById(ctx, pgstore.GetWorkoutByIdParams{
		ID:         workoutID,
		PersonalID: personalID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get workout: %w", err)
	}
	if dbWorkout == nil {
		return nil, fmt.Errorf("workout not found")
	}

	workout := s.convertWorkoutByIdFromDB(*dbWorkout)
	return &workout, nil
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, userID uuid.UUID, req *pgstore.CreateWorkoutParams) (*pgstore.Workout, error) {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return nil, fmt.Errorf("only personal trainers can create workouts")
	}

	workoutID := uuid.New()
	now := time.Now()

	// Convert week days
	weekDays := make([]pgstore.Day, len(req.WeekDays))
	for i, day := range req.WeekDays {
		weekDays[i] = pgstore.Day(day)
	}

	var level *pgstore.Level
	if req.Level != nil {
		l := pgstore.Level(*req.Level)
		level = &l
	}

	_, err = s.queries.CreateWorkout(ctx, pgstore.CreateWorkoutParams{
		ID:                       workoutID,
		Name:                     req.Name,
		Description:              req.Description,
		Thumbnail:                req.Thumbnail,
		VideoURL:                 req.VideoURL,
		RestTimeBetweenExercises: req.RestTimeBetweenExercises,
		Level:                    level,
		WeekDays:                 weekDays,
		Exclusive:                req.Exclusive,
		IsTemplate:               req.IsTemplate,
		Modality:                 req.Modality,
		PersonalID:               &userID,
		StudentID:                req.StudentID,
		PlanID:                   req.PlanID,
		CreatedAt:                now,
		UpdatedAt:                now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create workout: %w", err)
	}

	return s.GetWorkoutByID(ctx, workoutID, userID)
}

func (s *WorkoutService) UpdateWorkout(ctx context.Context, workoutID, userID uuid.UUID, req *pgstore.UpdateWorkoutParams) (*pgstore.Workout, error) {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return nil, fmt.Errorf("only personal trainers can update workouts")
	}

	now := time.Now()

	// Convert week days if provided
	var weekDays []pgstore.Day
	if len(req.WeekDays) > 0 {
		weekDays = make([]pgstore.Day, len(req.WeekDays))
		for i, day := range req.WeekDays {
			weekDays[i] = pgstore.Day(day)
		}
	}

	var level *pgstore.Level
	if req.Level != nil {
		l := pgstore.Level(*req.Level)
		level = &l
	}

	err = s.queries.UpdateWorkout(ctx, pgstore.UpdateWorkoutParams{
		ID:                       workoutID,
		Name:                     req.Name,
		Description:              req.Description,
		Thumbnail:                req.Thumbnail,
		VideoURL:                 req.VideoURL,
		RestTimeBetweenExercises: req.RestTimeBetweenExercises,
		Level:                    level,
		WeekDays:                 weekDays,
		Modality:                 req.Modality,
		UpdatedAt:                now,
		PersonalID:               &userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update workout: %w", err)
	}

	return s.GetWorkoutByID(ctx, workoutID, userID)
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, workoutID, userID uuid.UUID) error {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return fmt.Errorf("only personal trainers can delete workouts")
	}

	err = s.queries.DeleteWorkout(ctx, pgstore.DeleteWorkoutParams{
		ID:         workoutID,
		PersonalID: &userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	return nil
}

func (s *WorkoutService) GetExercises(ctx context.Context, userID uuid.UUID) ([]pgstore.ExercisesTemplate, error) {
	// Get user role to determine access
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	dbExercises, err := s.queries.GetExercises(ctx, personalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	exercises := make([]pgstore.ExercisesTemplate, len(dbExercises))
	for i, dbExercise := range dbExercises {
		exercises[i] = s.convertExerciseFromDB(dbExercise)
	}

	return exercises, nil
}

func (s *WorkoutService) GetExerciseByID(ctx context.Context, exerciseID, userID uuid.UUID) (*pgstore.ExercisesTemplate, error) {
	// Get user role to determine access
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	dbExercise, err := s.queries.GetExerciseById(ctx, pgstore.GetExerciseByIdParams{
		ID:         exerciseID,
		PersonalID: personalID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}
	if dbExercise == nil {
		return nil, fmt.Errorf("exercise not found")
	}

	exercise := s.convertExerciseFromDB(*dbExercise)
	return &exercise, nil
}

func (s *WorkoutService) CreateExercise(ctx context.Context, userID uuid.UUID, req *pgstore.CreateExerciseParams) (*pgstore.ExercisesTemplate, error) {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return nil, fmt.Errorf("only personal trainers can create exercises")
	}

	exerciseID := uuid.New()
	now := time.Now()

	_, err = s.queries.CreateExercise(ctx, pgstore.CreateExerciseParams{
		ID:                  exerciseID,
		Name:                req.Name,
		Thumbnail:           req.Thumbnail,
		VideoURL:            req.VideoURL,
		Load:                req.Load,
		Sets:                req.Sets,
		Reps:                req.Reps,
		RestTimeBetweenSets: req.RestTimeBetweenSets,
		PersonalID:          &userID,
		CreatedAt:           now,
		UpdatedAt:           now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create exercise: %w", err)
	}

	return s.GetExerciseByID(ctx, exerciseID, userID)
}

func (s *WorkoutService) UpdateExercise(ctx context.Context, exerciseID, userID uuid.UUID, req *pgstore.UpdateExerciseParams) (*pgstore.ExercisesTemplate, error) {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return nil, fmt.Errorf("only personal trainers can update exercises")
	}

	now := time.Now()

	err = s.queries.UpdateExercise(ctx, pgstore.UpdateExerciseParams{
		ID:                  exerciseID,
		Name:                req.Name,
		Thumbnail:           req.Thumbnail,
		VideoURL:            req.VideoURL,
		Load:                req.Load,
		Sets:                req.Sets,
		Reps:                req.Reps,
		RestTimeBetweenSets: req.RestTimeBetweenSets,
		UpdatedAt:           now,
		PersonalID:          &userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update exercise: %w", err)
	}

	return s.GetExerciseByID(ctx, exerciseID, userID)
}

func (s *WorkoutService) DeleteExercise(ctx context.Context, exerciseID, userID uuid.UUID) error {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return fmt.Errorf("only personal trainers can delete exercises")
	}

	err = s.queries.DeleteExercise(ctx, pgstore.DeleteExerciseParams{
		ID:         exerciseID,
		PersonalID: &userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	return nil
}

// Helper functions to convert between database and domain models
func (s *WorkoutService) convertWorkoutFromDB(dbWorkout pgstore.GetWorkoutsRow) pgstore.Workout {
	weekDays := make([]pgstore.Day, len(dbWorkout.WeekDays))
	for i, day := range dbWorkout.WeekDays {
		weekDays[i] = pgstore.Day(day)
	}

	var level *pgstore.Level
	if dbWorkout.Level != nil {
		l := pgstore.Level(*dbWorkout.Level)
		level = &l
	}

	return pgstore.Workout{
		ID:                       dbWorkout.ID,
		Name:                     dbWorkout.Name,
		Description:              dbWorkout.Description,
		Thumbnail:                dbWorkout.Thumbnail,
		VideoURL:                 dbWorkout.VideoURL,
		RestTimeBetweenExercises: dbWorkout.RestTimeBetweenExercises,
		Level:                    level,
		WeekDays:                 weekDays,
		Exclusive:                dbWorkout.Exclusive,
		IsTemplate:               dbWorkout.IsTemplate,
		Modality:                 dbWorkout.Modality,
		PersonalID:               dbWorkout.PersonalID,
		StudentID:                dbWorkout.StudentID,
		PlanID:                   dbWorkout.PlanID,
		CreatedAt:                dbWorkout.CreatedAt,
		UpdatedAt:                dbWorkout.UpdatedAt,
	}
}

func (s *WorkoutService) convertWorkoutByIdFromDB(dbWorkout pgstore.GetWorkoutByIdRow) pgstore.Workout {
	weekDays := make([]pgstore.Day, len(dbWorkout.WeekDays))
	for i, day := range dbWorkout.WeekDays {
		weekDays[i] = pgstore.Day(day)
	}

	var level *pgstore.Level
	if dbWorkout.Level != nil {
		l := pgstore.Level(*dbWorkout.Level)
		level = &l
	}

	return pgstore.Workout{
		ID:                       dbWorkout.ID,
		Name:                     dbWorkout.Name,
		Description:              dbWorkout.Description,
		Thumbnail:                dbWorkout.Thumbnail,
		VideoURL:                 dbWorkout.VideoURL,
		RestTimeBetweenExercises: dbWorkout.RestTimeBetweenExercises,
		Level:                    level,
		WeekDays:                 weekDays,
		Exclusive:                dbWorkout.Exclusive,
		IsTemplate:               dbWorkout.IsTemplate,
		Modality:                 dbWorkout.Modality,
		PersonalID:               dbWorkout.PersonalID,
		StudentID:                dbWorkout.StudentID,
		PlanID:                   dbWorkout.PlanID,
		CreatedAt:                dbWorkout.CreatedAt,
		UpdatedAt:                dbWorkout.UpdatedAt,
	}
}

func (s *WorkoutService) convertExerciseFromDB(dbExercise pgstore.ExercisesTemplate) pgstore.ExercisesTemplate {
	return pgstore.ExercisesTemplate{
		ID:                  dbExercise.ID,
		Name:                dbExercise.Name,
		Thumbnail:           dbExercise.Thumbnail,
		VideoURL:            dbExercise.VideoURL,
		Load:                dbExercise.Load,
		Sets:                dbExercise.Sets,
		Reps:                dbExercise.Reps,
		RestTimeBetweenSets: dbExercise.RestTimeBetweenSets,
		PersonalID:          dbExercise.PersonalID,
		CreatedAt:           dbExercise.CreatedAt,
		UpdatedAt:           dbExercise.UpdatedAt,
	}
}
