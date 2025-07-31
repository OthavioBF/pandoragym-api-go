package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

type WorkoutService struct {
	queries *pgstore.Queries
	pool    *pgxpool.Pool
}

func NewWorkoutService(queries *pgstore.Queries, pool *pgxpool.Pool) *WorkoutService {
	return &WorkoutService{
		queries: queries,
		pool:    pool,
	}
}

// GetWorkouts retrieves workouts based on user role
func (s *WorkoutService) GetWorkouts(ctx context.Context, userID uuid.UUID) ([]pgstore.GetWorkoutsRow, error) {
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	workouts, err := s.queries.GetWorkouts(ctx, personalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workouts: %w", err)
	}

	return workouts, nil
}

// GetWorkoutByID retrieves a specific workout with access control
func (s *WorkoutService) GetWorkoutByID(ctx context.Context, workoutID, userID uuid.UUID) (*pgstore.GetWorkoutByIdRow, error) {
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	workout, err := s.queries.GetWorkoutById(ctx, pgstore.GetWorkoutByIdParams{
		ID:         workoutID,
		PersonalID: personalID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get workout: %w", err)
	}

	// Access control: personal trainers can only access their own workouts
	if role == pgstore.RolePersonal && workout.PersonalID != nil && *workout.PersonalID != userID {
		return nil, fmt.Errorf("access denied: workout belongs to another trainer")
	}

	return workout, nil
}

// CreateWorkout creates a new workout with exercises in a transaction
func (s *WorkoutService) CreateWorkout(ctx context.Context, req pgstore.CreateWorkoutRequest, userID uuid.UUID) (*pgstore.WorkoutResponse, error) {
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	workoutID := uuid.New()
	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	// Start transaction for atomic workout creation
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txQueries := s.queries.WithTx(tx)

	// Set default values
	weekDays := []pgstore.Day{}
	if req.WeekDays != nil {
		weekDays = *req.WeekDays
	}

	exclusive := false
	if req.Exclusive != nil {
		exclusive = *req.Exclusive
	}

	isTemplate := false
	if req.IsTemplate != nil {
		isTemplate = *req.IsTemplate
	}

	// Create the workout
	createdWorkoutID, err := txQueries.CreateWorkout(ctx, pgstore.CreateWorkoutParams{
		ID:                       workoutID,
		Name:                     req.Name,
		Description:              &req.Description,
		Thumbnail:                req.Thumbnail,
		VideoURL:                 req.VideoURL,
		RestTimeBetweenExercises: req.RestTimeBetweenExercises,
		Level:                    req.Level,
		WeekDays:                 weekDays,
		Exclusive:                exclusive,
		IsTemplate:               isTemplate,
		Modality:                 req.Modality,
		PersonalID:               personalID,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create workout: %w", err)
	}

	// Add exercises to the workout
	for i, exercise := range req.Exercises {
		exerciseID := uuid.New()
		_, err = txQueries.AddExerciseToWorkout(ctx, pgstore.AddExerciseToWorkoutParams{
			ID:                  exerciseID,
			Name:                exercise.Name,
			Thumbnail:           exercise.Thumbnail,
			VideoURL:            exercise.VideoURL,
			Sets:                exercise.Sets,
			Reps:                exercise.Reps,
			RestTimeBetweenSets: exercise.RestTimeBetweenSets,
			Load:                exercise.Load,
			WorkoutID:           workoutID,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add exercise %d to workout: %w", i+1, err)
		}
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return workout response
	return &pgstore.WorkoutResponse{
		ID:                       createdWorkoutID,
		Name:                     req.Name,
		Description:              &req.Description,
		Thumbnail:                req.Thumbnail,
		VideoURL:                 req.VideoURL,
		RestTimeBetweenExercises: req.RestTimeBetweenExercises,
		Level:                    req.Level,
		WeekDays:                 weekDays,
		Exclusive:                exclusive,
		IsTemplate:               isTemplate,
		Modality:                 req.Modality,
		PersonalID:               personalID,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}, nil
}

// UpdateWorkout updates an existing workout with access control
func (s *WorkoutService) UpdateWorkout(ctx context.Context, workoutID uuid.UUID, req pgstore.UpdateWorkoutRequest, userID uuid.UUID) (*pgstore.WorkoutResponse, error) {
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	// Check if user has permission to update this workout
	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	workout, err := s.queries.GetWorkoutById(ctx, pgstore.GetWorkoutByIdParams{
		ID:         workoutID,
		PersonalID: personalID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get workout for update: %w", err)
	}

	// Access control
	if role == pgstore.RolePersonal && workout.PersonalID != nil && *workout.PersonalID != userID {
		return nil, fmt.Errorf("access denied: cannot update workout belonging to another trainer")
	}

	// Update the workout
	updateParams := pgstore.UpdateWorkoutParams{
		ID:        workoutID,
		UpdatedAt: time.Now(),
	}

	if req.Name != nil {
		updateParams.Name = req.Name
	}
	if req.Description != nil {
		updateParams.Description = req.Description
	}
	if req.Thumbnail != nil {
		updateParams.Thumbnail = req.Thumbnail
	}
	if req.VideoURL != nil {
		updateParams.VideoURL = req.VideoURL
	}
	if req.RestTimeBetweenExercises != nil {
		updateParams.RestTimeBetweenExercises = req.RestTimeBetweenExercises
	}
	if req.Level != nil {
		updateParams.Level = req.Level
	}
	if req.WeekDays != nil {
		updateParams.WeekDays = *req.WeekDays
	}
	if req.Modality != nil {
		updateParams.Modality = req.Modality
	}

	err = s.queries.UpdateWorkout(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update workout: %w", err)
	}

	// Get updated workout to return
	updatedWorkout, err := s.queries.GetWorkoutById(ctx, pgstore.GetWorkoutByIdParams{
		ID:         workoutID,
		PersonalID: personalID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated workout: %w", err)
	}

	return &pgstore.WorkoutResponse{
		ID:                       updatedWorkout.ID,
		Name:                     updatedWorkout.Name,
		Description:              updatedWorkout.Description,
		Thumbnail:                updatedWorkout.Thumbnail,
		VideoURL:                 updatedWorkout.VideoURL,
		RestTimeBetweenExercises: updatedWorkout.RestTimeBetweenExercises,
		Level:                    updatedWorkout.Level,
		WeekDays:                 updatedWorkout.WeekDays,
		Exclusive:                updatedWorkout.Exclusive,
		IsTemplate:               updatedWorkout.IsTemplate,
		Modality:                 updatedWorkout.Modality,
		PersonalID:               updatedWorkout.PersonalID,
		StudentID:                updatedWorkout.StudentID,
		CreatedAt:                updatedWorkout.CreatedAt,
		UpdatedAt:                updatedWorkout.UpdatedAt,
	}, nil
}

// DeleteWorkout deletes a workout with access control
func (s *WorkoutService) DeleteWorkout(ctx context.Context, workoutID, userID uuid.UUID) error {
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}

	// Check if user has permission to delete this workout
	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	workout, err := s.queries.GetWorkoutById(ctx, pgstore.GetWorkoutByIdParams{
		ID:         workoutID,
		PersonalID: personalID,
	})
	if err != nil {
		return fmt.Errorf("failed to get workout for deletion: %w", err)
	}

	// Access control
	if role == pgstore.RolePersonal && workout.PersonalID != nil && *workout.PersonalID != userID {
		return fmt.Errorf("access denied: cannot delete workout belonging to another trainer")
	}

	err = s.queries.DeleteWorkout(ctx, pgstore.DeleteWorkoutParams{
		ID: workoutID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	return nil
}

// Exercise-related methods

// GetAllExercises retrieves all exercise templates
func (s *WorkoutService) GetAllExercises(ctx context.Context) ([]pgstore.ExerciseResponse, error) {
	exercises, err := s.queries.GetExercises(ctx, nil) // Get all exercises
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	var response []pgstore.ExerciseResponse
	for _, exercise := range exercises {
		response = append(response, pgstore.ExerciseResponse{
			ID:                  exercise.ID,
			Name:                exercise.Name,
			Thumbnail:           exercise.Thumbnail,
			VideoURL:            exercise.VideoURL,
			Load:                exercise.Load,
			Sets:                exercise.Sets,
			Reps:                exercise.Reps,
			RestTimeBetweenSets: exercise.RestTimeBetweenSets,
			PersonalID:          exercise.PersonalID,
			CreatedAt:           exercise.CreatedAt,
			UpdatedAt:           exercise.UpdatedAt,
		})
	}

	return response, nil
}

// GetExerciseByID retrieves a specific exercise template
func (s *WorkoutService) GetExerciseByID(ctx context.Context, exerciseID uuid.UUID) (*pgstore.ExerciseResponse, error) {
	exercise, err := s.queries.GetExerciseById(ctx, pgstore.GetExerciseByIdParams{
		ID: exerciseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	return &pgstore.ExerciseResponse{
		ID:                  exercise.ID,
		Name:                exercise.Name,
		Thumbnail:           exercise.Thumbnail,
		VideoURL:            exercise.VideoURL,
		Load:                exercise.Load,
		Sets:                exercise.Sets,
		Reps:                exercise.Reps,
		RestTimeBetweenSets: exercise.RestTimeBetweenSets,
		PersonalID:          exercise.PersonalID,
		CreatedAt:           exercise.CreatedAt,
		UpdatedAt:           exercise.UpdatedAt,
	}, nil
}

// CreateExercise creates a new exercise template
func (s *WorkoutService) CreateExercise(ctx context.Context, req pgstore.CreateExerciseRequest) (*pgstore.ExerciseResponse, error) {
	exerciseID := uuid.New()

	createdExerciseID, err := s.queries.CreateExercise(ctx, pgstore.CreateExerciseParams{
		ID:                  exerciseID,
		Name:                req.Name,
		Thumbnail:           req.Thumbnail,
		VideoURL:            req.VideoURL,
		Load:                req.Load,
		Sets:                req.Sets,
		Reps:                req.Reps,
		RestTimeBetweenSets: req.RestTimeBetweenSets,
		PersonalID:          req.PersonalID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create exercise: %w", err)
	}

	// Get the created exercise to return
	createdExercise, err := s.queries.GetExerciseById(ctx, pgstore.GetExerciseByIdParams{
		ID: createdExerciseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get created exercise: %w", err)
	}

	return &pgstore.ExerciseResponse{
		ID:                  createdExercise.ID,
		Name:                createdExercise.Name,
		Thumbnail:           createdExercise.Thumbnail,
		VideoURL:            createdExercise.VideoURL,
		Load:                createdExercise.Load,
		Sets:                createdExercise.Sets,
		Reps:                createdExercise.Reps,
		RestTimeBetweenSets: createdExercise.RestTimeBetweenSets,
		PersonalID:          createdExercise.PersonalID,
		CreatedAt:           createdExercise.CreatedAt,
		UpdatedAt:           createdExercise.UpdatedAt,
	}, nil
}

// UpdateExercise updates an existing exercise template
func (s *WorkoutService) UpdateExercise(ctx context.Context, exerciseID uuid.UUID, req pgstore.UpdateExerciseRequest) (*pgstore.ExerciseResponse, error) {
	updateParams := pgstore.UpdateExerciseParams{
		ID:        exerciseID,
		UpdatedAt: time.Now(),
	}

	if req.Name != nil {
		updateParams.Name = req.Name
	}
	if req.Thumbnail != nil {
		updateParams.Thumbnail = req.Thumbnail
	}
	if req.VideoURL != nil {
		updateParams.VideoURL = req.VideoURL
	}
	if req.Load != nil {
		updateParams.Load = req.Load
	}
	if req.Sets != nil {
		updateParams.Sets = req.Sets
	}
	if req.Reps != nil {
		updateParams.Reps = req.Reps
	}
	if req.RestTimeBetweenSets != nil {
		updateParams.RestTimeBetweenSets = req.RestTimeBetweenSets
	}

	err := s.queries.UpdateExercise(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update exercise: %w", err)
	}

	// Get updated exercise to return
	updatedExercise, err := s.queries.GetExerciseById(ctx, pgstore.GetExerciseByIdParams{
		ID: exerciseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated exercise: %w", err)
	}

	return &pgstore.ExerciseResponse{
		ID:                  updatedExercise.ID,
		Name:                updatedExercise.Name,
		Thumbnail:           updatedExercise.Thumbnail,
		VideoURL:            updatedExercise.VideoURL,
		Load:                updatedExercise.Load,
		Sets:                updatedExercise.Sets,
		Reps:                updatedExercise.Reps,
		RestTimeBetweenSets: updatedExercise.RestTimeBetweenSets,
		PersonalID:          updatedExercise.PersonalID,
		CreatedAt:           updatedExercise.CreatedAt,
		UpdatedAt:           updatedExercise.UpdatedAt,
	}, nil
}

// DeleteExercise deletes an exercise template
func (s *WorkoutService) DeleteExercise(ctx context.Context, exerciseID uuid.UUID) error {
	err := s.queries.DeleteExercise(ctx, pgstore.DeleteExerciseParams{
		ID: exerciseID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	return nil
}

// AddExerciseToWorkout adds an exercise to a specific workout
func (s *WorkoutService) AddExerciseToWorkout(ctx context.Context, workoutID, exerciseID, userID uuid.UUID, sets, reps int, restTime *int) error {
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}

	// Check if user has permission to modify this workout
	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	workout, err := s.queries.GetWorkoutById(ctx, pgstore.GetWorkoutByIdParams{
		ID:         workoutID,
		PersonalID: personalID,
	})
	if err != nil {
		return fmt.Errorf("failed to get workout: %w", err)
	}

	// Access control
	if role == pgstore.RolePersonal && workout.PersonalID != nil && *workout.PersonalID != userID {
		return fmt.Errorf("access denied: cannot modify workout belonging to another trainer")
	}

	// Get exercise template
	exercise, err := s.queries.GetExerciseForWorkout(ctx, exerciseID)
	if err != nil {
		return fmt.Errorf("failed to get exercise template: %w", err)
	}

	// Add exercise to workout
	var load int32
	if exercise.Load != nil {
		load = *exercise.Load
	}

	_, err = s.queries.AddExerciseToWorkout(ctx, pgstore.AddExerciseToWorkoutParams{
		ID:                  uuid.New(),
		Name:                exercise.Name,
		Thumbnail:           exercise.Thumbnail,
		VideoURL:            exercise.VideoURL,
		Sets:                int32(sets),
		Reps:                int32(reps),
		RestTimeBetweenSets: int32(*restTime),
		Load:                load,
		WorkoutID:           workoutID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to add exercise to workout: %w", err)
	}

	return nil
}

// RemoveExerciseFromWorkout removes an exercise from a specific workout
func (s *WorkoutService) RemoveExerciseFromWorkout(ctx context.Context, workoutID, exerciseID, userID uuid.UUID) error {
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}

	// Check if user has permission to modify this workout
	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	workout, err := s.queries.GetWorkoutById(ctx, pgstore.GetWorkoutByIdParams{
		ID:         workoutID,
		PersonalID: personalID,
	})
	if err != nil {
		return fmt.Errorf("failed to get workout: %w", err)
	}

	// Access control
	if role == pgstore.RolePersonal && workout.PersonalID != nil && *workout.PersonalID != userID {
		return fmt.Errorf("access denied: cannot modify workout belonging to another trainer")
	}

	err = s.queries.RemoveExerciseFromWorkout(ctx, pgstore.RemoveExerciseFromWorkoutParams{
		ID:        exerciseID,
		WorkoutID: workoutID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove exercise from workout: %w", err)
	}

	return nil
}

// GetWorkoutExercises retrieves all exercises for a specific workout
func (s *WorkoutService) GetWorkoutExercises(ctx context.Context, workoutID, userID uuid.UUID) ([]pgstore.ExercisesSetup, error) {
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	// Check if user has permission to view this workout
	var personalID *uuid.UUID
	if role == pgstore.RolePersonal {
		personalID = &userID
	}

	workout, err := s.queries.GetWorkoutById(ctx, pgstore.GetWorkoutByIdParams{
		ID:         workoutID,
		PersonalID: personalID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get workout: %w", err)
	}

	// Access control
	if role == pgstore.RolePersonal && workout.PersonalID != nil && *workout.PersonalID != userID {
		return nil, fmt.Errorf("access denied: cannot view workout belonging to another trainer")
	}

	exercises, err := s.queries.GetExercisesByWorkoutId(ctx, workoutID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout exercises: %w", err)
	}

	return exercises, nil
}

// Workout execution and history methods

// FinishWorkout records a completed workout session
func (s *WorkoutService) FinishWorkout(ctx context.Context, userID, workoutID string, duration int, exercises []map[string]interface{}, notes string) error {
	// Implementation for recording workout completion
	// This would typically create a workout history record
	return fmt.Errorf("finish workout not implemented yet")
}

// ExecuteWorkout starts a workout session
func (s *WorkoutService) ExecuteWorkout(ctx context.Context, userID, workoutID string) (interface{}, error) {
	// Implementation for starting a workout session
	return nil, fmt.Errorf("execute workout not implemented yet")
}

// RateWorkout allows users to rate completed workouts
func (s *WorkoutService) RateWorkout(ctx context.Context, userID, workoutID string, rating int, comment string) error {
	// Implementation for workout rating
	return fmt.Errorf("rate workout not implemented yet")
}

// GetWorkoutHistory retrieves user's workout history
func (s *WorkoutService) GetWorkoutHistory(ctx context.Context, userID string) ([]pgstore.WorkoutHistoryResponse, error) {
	// Implementation for getting workout history
	return []pgstore.WorkoutHistoryResponse{}, nil
}

// Template and program methods (simplified implementations)

// GetExerciseTemplates retrieves all exercise templates
func (s *WorkoutService) GetExerciseTemplates(ctx context.Context) ([]pgstore.ExerciseTemplateResponse, error) {
	return []pgstore.ExerciseTemplateResponse{}, nil
}

// CreateExerciseTemplate creates a new exercise template
func (s *WorkoutService) CreateExerciseTemplate(ctx context.Context, name, description, videoURL, instructions, category string, muscleGroups, equipment []string, difficulty string) (*pgstore.ExerciseTemplateResponse, error) {
	return nil, fmt.Errorf("create exercise template not implemented yet")
}

// GetAllExerciseTemplates retrieves all exercise templates
func (s *WorkoutService) GetAllExerciseTemplates(ctx context.Context) ([]pgstore.ExerciseTemplateResponse, error) {
	return []pgstore.ExerciseTemplateResponse{}, nil
}

// DeleteExerciseTemplate deletes an exercise template
func (s *WorkoutService) DeleteExerciseTemplate(ctx context.Context, templateID string) error {
	return nil
}

// GetWorkoutTemplates retrieves all workout templates
func (s *WorkoutService) GetWorkoutTemplates(ctx context.Context) ([]pgstore.WorkoutTemplateResponse, error) {
	return []pgstore.WorkoutTemplateResponse{}, nil
}

// CreateWorkoutTemplate creates a new workout template
func (s *WorkoutService) CreateWorkoutTemplate(ctx context.Context, name, description, thumbnail, category, difficulty string, duration int, weekDays []string, exercises []map[string]interface{}, tags []string) (*pgstore.WorkoutTemplateResponse, error) {
	return nil, fmt.Errorf("create workout template not implemented yet")
}

// GetAllWorkoutTemplates retrieves all workout templates
func (s *WorkoutService) GetAllWorkoutTemplates(ctx context.Context) ([]pgstore.WorkoutTemplateResponse, error) {
	return []pgstore.WorkoutTemplateResponse{}, nil
}

// DeleteWorkoutTemplate deletes a workout template
func (s *WorkoutService) DeleteWorkoutTemplate(ctx context.Context, templateID string) error {
	return nil
}

// GetAllPrograms retrieves all training programs for a user
func (s *WorkoutService) GetAllPrograms(ctx context.Context, userID string) ([]pgstore.TrainingProgramResponse, error) {
	return []pgstore.TrainingProgramResponse{}, nil
}

// GetFreePrograms retrieves all free training programs
func (s *WorkoutService) GetFreePrograms(ctx context.Context) ([]pgstore.TrainingProgramResponse, error) {
	return []pgstore.TrainingProgramResponse{}, nil
}

// GetFreeProgramByID retrieves a specific free training program
func (s *WorkoutService) GetFreeProgramByID(ctx context.Context, programID string) (*pgstore.TrainingProgramResponse, error) {
	return nil, fmt.Errorf("get free program by ID not implemented yet")
}
