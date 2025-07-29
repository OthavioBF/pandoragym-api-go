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

	if role == pgstore.RolePersonal && workout.PersonalID != nil && *workout.PersonalID != userID {
		return nil, fmt.Errorf("access denied: workout belongs to another trainer")
	}

	return workout, nil
}

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

	err = s.queries.CreateWorkout(ctx, pgstore.CreateWorkoutParams{
		ID:          workoutID,
		Name:        req.Name,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Category:    req.Category,
		Difficulty:  req.Difficulty,
		Duration:    req.Duration,
		WeekDays:    req.WeekDays,
		PersonalID:  personalID,
		IsPublic:    req.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create workout: %w", err)
	}

	for _, exercise := range req.Exercises {
		exerciseID, err := uuid.Parse(exercise.ExerciseID)
		if err != nil {
			continue
		}

		err = s.queries.CreateWorkoutExercise(ctx, pgstore.CreateWorkoutExerciseParams{
			ID:         uuid.New(),
			WorkoutID:  workoutID,
			ExerciseID: exerciseID,
			Sets:       exercise.Sets,
			Reps:       exercise.Reps,
			RestTime:   exercise.RestTime,
			Order:      exercise.Order,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add exercise to workout: %w", err)
		}
	}

	return &pgstore.WorkoutResponse{
		ID:          workoutID,
		Name:        req.Name,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Category:    req.Category,
		Difficulty:  req.Difficulty,
		Duration:    req.Duration,
		WeekDays:    req.WeekDays,
		PersonalID:  personalID,
		IsPublic:    req.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s *WorkoutService) UpdateWorkout(ctx context.Context, workoutID uuid.UUID, req pgstore.UpdateWorkoutRequest, userID uuid.UUID) (*pgstore.WorkoutResponse, error) {
	workout, err := s.GetWorkoutByID(ctx, workoutID, userID)
	if err != nil {
		return nil, err
	}

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
	if req.Category != nil {
		updateParams.Category = req.Category
	}
	if req.Difficulty != nil {
		updateParams.Difficulty = req.Difficulty
	}
	if req.Duration != nil {
		updateParams.Duration = req.Duration
	}
	if req.WeekDays != nil {
		updateParams.WeekDays = req.WeekDays
	}
	if req.IsPublic != nil {
		updateParams.IsPublic = req.IsPublic
	}

	err = s.queries.UpdateWorkout(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update workout: %w", err)
	}

	return &pgstore.WorkoutResponse{
		ID:          workoutID,
		Name:        *updateParams.Name,
		Description: updateParams.Description,
		Thumbnail:   updateParams.Thumbnail,
		Category:    *updateParams.Category,
		Difficulty:  *updateParams.Difficulty,
		Duration:    *updateParams.Duration,
		WeekDays:    updateParams.WeekDays,
		PersonalID:  workout.PersonalID,
		IsPublic:    *updateParams.IsPublic,
		UpdatedAt:   time.Now(),
	}, nil
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, workoutID, userID uuid.UUID) error {
	_, err := s.GetWorkoutByID(ctx, workoutID, userID)
	if err != nil {
		return err
	}

	err = s.queries.DeleteWorkout(ctx, pgstore.DeleteWorkoutParams{ID: workoutID})
	if err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	return nil
}

func (s *WorkoutService) GetAllExercises(ctx context.Context) ([]pgstore.ExerciseResponse, error) {
	exercises, err := s.queries.GetExercises(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	var response []pgstore.ExerciseResponse
	for _, exercise := range exercises {
		response = append(response, pgstore.ExerciseResponse{
			ID:           exercise.ID,
			Name:         exercise.Name,
			Description:  exercise.Description,
			VideoURL:     exercise.VideoURL,
			Instructions: exercise.Instructions,
			Category:     exercise.Category,
			MuscleGroups: exercise.MuscleGroups,
			Equipment:    exercise.Equipment,
			Difficulty:   exercise.Difficulty,
			CreatedAt:    exercise.CreatedAt,
			UpdatedAt:    exercise.UpdatedAt,
		})
	}

	return response, nil
}

func (s *WorkoutService) GetExerciseByID(ctx context.Context, exerciseID uuid.UUID) (*pgstore.ExerciseResponse, error) {
	exercise, err := s.queries.GetExerciseById(ctx, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	return &pgstore.ExerciseResponse{
		ID:           exercise.ID,
		Name:         exercise.Name,
		Description:  exercise.Description,
		VideoURL:     exercise.VideoURL,
		Instructions: exercise.Instructions,
		Category:     exercise.Category,
		MuscleGroups: exercise.MuscleGroups,
		Equipment:    exercise.Equipment,
		Difficulty:   exercise.Difficulty,
		CreatedAt:    exercise.CreatedAt,
		UpdatedAt:    exercise.UpdatedAt,
	}, nil
}

func (s *WorkoutService) CreateExercise(ctx context.Context, req pgstore.CreateExerciseRequest) (*pgstore.ExerciseResponse, error) {
	exerciseID := uuid.New()

	err := s.queries.CreateExercise(ctx, pgstore.CreateExerciseParams{
		ID:           exerciseID,
		Name:         req.Name,
		Description:  req.Description,
		VideoURL:     req.VideoURL,
		Instructions: req.Instructions,
		Category:     req.Category,
		MuscleGroups: req.MuscleGroups,
		Equipment:    req.Equipment,
		Difficulty:   req.Difficulty,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create exercise: %w", err)
	}

	return &pgstore.ExerciseResponse{
		ID:           exerciseID,
		Name:         req.Name,
		Description:  req.Description,
		VideoURL:     req.VideoURL,
		Instructions: req.Instructions,
		Category:     req.Category,
		MuscleGroups: req.MuscleGroups,
		Equipment:    req.Equipment,
		Difficulty:   req.Difficulty,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (s *WorkoutService) UpdateExercise(ctx context.Context, exerciseID uuid.UUID, req pgstore.UpdateExerciseRequest) (*pgstore.ExerciseResponse, error) {
	updateParams := pgstore.UpdateExerciseParams{
		ID:        exerciseID,
		UpdatedAt: time.Now(),
	}

	if req.Name != nil {
		updateParams.Name = req.Name
	}
	if req.Description != nil {
		updateParams.Description = req.Description
	}
	if req.VideoURL != nil {
		updateParams.VideoURL = req.VideoURL
	}
	if req.Instructions != nil {
		updateParams.Instructions = req.Instructions
	}
	if req.Category != nil {
		updateParams.Category = req.Category
	}
	if req.MuscleGroups != nil {
		updateParams.MuscleGroups = req.MuscleGroups
	}
	if req.Equipment != nil {
		updateParams.Equipment = req.Equipment
	}
	if req.Difficulty != nil {
		updateParams.Difficulty = req.Difficulty
	}

	err := s.queries.UpdateExercise(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update exercise: %w", err)
	}

	return s.GetExerciseByID(ctx, exerciseID)
}

func (s *WorkoutService) DeleteExercise(ctx context.Context, exerciseID uuid.UUID) error {
	err := s.queries.DeleteExercise(ctx, exerciseID)
	if err != nil {
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	return nil
}

func (s *WorkoutService) AddExerciseToWorkout(ctx context.Context, workoutID, exerciseID, userID uuid.UUID, sets, reps int, restTime *int) error {
	_, err := s.GetWorkoutByID(ctx, workoutID, userID)
	if err != nil {
		return err
	}

	order, err := s.queries.GetNextExerciseOrder(ctx, workoutID)
	if err != nil {
		order = 1
	}

	err = s.queries.CreateWorkoutExercise(ctx, pgstore.CreateWorkoutExerciseParams{
		ID:         uuid.New(),
		WorkoutID:  workoutID,
		ExerciseID: exerciseID,
		Sets:       int32(sets),
		Reps:       int32(reps),
		RestTime:   restTime,
		Order:      int32(order),
	})
	if err != nil {
		return fmt.Errorf("failed to add exercise to workout: %w", err)
	}

	return nil
}

func (s *WorkoutService) RemoveExerciseFromWorkout(ctx context.Context, workoutID, exerciseID, userID uuid.UUID) error {
	_, err := s.GetWorkoutByID(ctx, workoutID, userID)
	if err != nil {
		return err
	}

	err = s.queries.DeleteWorkoutExercise(ctx, pgstore.DeleteWorkoutExerciseParams{
		WorkoutID:  workoutID,
		ExerciseID: exerciseID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove exercise from workout: %w", err)
	}

	return nil
}

func (s *WorkoutService) FinishWorkout(ctx context.Context, userID, workoutID string, duration int, exercises []map[string]interface{}, notes string) error {
	workoutUUID, err := uuid.Parse(workoutID)
	if err != nil {
		return fmt.Errorf("invalid workout ID: %w", err)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = s.queries.CreateWorkoutHistory(ctx, pgstore.CreateWorkoutHistoryParams{
		ID:        uuid.New(),
		UserID:    userUUID,
		WorkoutID: workoutUUID,
		Duration:  int32(duration),
		Notes:     &notes,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to create workout history: %w", err)
	}

	return nil
}

func (s *WorkoutService) ExecuteWorkout(ctx context.Context, userID, workoutID string) (interface{}, error) {
	workoutUUID, err := uuid.Parse(workoutID)
	if err != nil {
		return nil, fmt.Errorf("invalid workout ID: %w", err)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	workout, err := s.GetWorkoutByID(ctx, workoutUUID, userUUID)
	if err != nil {
		return nil, err
	}

	exercises, err := s.queries.GetWorkoutExercises(ctx, workoutUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout exercises: %w", err)
	}

	return map[string]interface{}{
		"workout":   workout,
		"exercises": exercises,
		"started_at": time.Now(),
	}, nil
}

func (s *WorkoutService) RateWorkout(ctx context.Context, userID, workoutID string, rating int, comment string) error {
	workoutUUID, err := uuid.Parse(workoutID)
	if err != nil {
		return fmt.Errorf("invalid workout ID: %w", err)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = s.queries.CreateWorkoutRating(ctx, pgstore.CreateWorkoutRatingParams{
		ID:        uuid.New(),
		UserID:    userUUID,
		WorkoutID: workoutUUID,
		Rating:    int32(rating),
		Comment:   &comment,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to create workout rating: %w", err)
	}

	return nil
}

func (s *WorkoutService) GetWorkoutHistory(ctx context.Context, userID string) ([]pgstore.WorkoutHistoryResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	history, err := s.queries.GetWorkoutHistory(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout history: %w", err)
	}

	var response []pgstore.WorkoutHistoryResponse
	for _, h := range history {
		response = append(response, pgstore.WorkoutHistoryResponse{
			ID:          h.ID,
			WorkoutName: h.WorkoutName,
			Duration:    h.Duration,
			Notes:       h.Notes,
			CreatedAt:   h.CreatedAt,
		})
	}

	return response, nil
}

func (s *WorkoutService) GetExerciseTemplates(ctx context.Context) ([]pgstore.ExerciseTemplateResponse, error) {
	templates, err := s.queries.GetExerciseTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercise templates: %w", err)
	}

	var response []pgstore.ExerciseTemplateResponse
	for _, template := range templates {
		response = append(response, pgstore.ExerciseTemplateResponse{
			ID:           template.ID,
			Name:         template.Name,
			Description:  template.Description,
			VideoURL:     template.VideoURL,
			Instructions: template.Instructions,
			Category:     template.Category,
			MuscleGroups: template.MuscleGroups,
			Equipment:    template.Equipment,
			Difficulty:   template.Difficulty,
		})
	}

	return response, nil
}

func (s *WorkoutService) CreateExerciseTemplate(ctx context.Context, name, description, videoURL, instructions, category string, muscleGroups, equipment []string, difficulty string) (*pgstore.ExerciseTemplateResponse, error) {
	templateID := uuid.New()

	err := s.queries.CreateExerciseTemplate(ctx, pgstore.CreateExerciseTemplateParams{
		ID:           templateID,
		Name:         name,
		Description:  description,
		VideoURL:     videoURL,
		Instructions: instructions,
		Category:     category,
		MuscleGroups: muscleGroups,
		Equipment:    equipment,
		Difficulty:   difficulty,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create exercise template: %w", err)
	}

	return &pgstore.ExerciseTemplateResponse{
		ID:           templateID,
		Name:         name,
		Description:  description,
		VideoURL:     videoURL,
		Instructions: instructions,
		Category:     category,
		MuscleGroups: muscleGroups,
		Equipment:    equipment,
		Difficulty:   difficulty,
	}, nil
}

func (s *WorkoutService) GetAllExerciseTemplates(ctx context.Context) ([]pgstore.ExerciseTemplateResponse, error) {
	return s.GetExerciseTemplates(ctx)
}

func (s *WorkoutService) DeleteExerciseTemplate(ctx context.Context, templateID string) error {
	id, err := uuid.Parse(templateID)
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	err = s.queries.DeleteExerciseTemplate(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete exercise template: %w", err)
	}

	return nil
}

func (s *WorkoutService) GetWorkoutTemplates(ctx context.Context) ([]pgstore.WorkoutTemplateResponse, error) {
	templates, err := s.queries.GetWorkoutTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout templates: %w", err)
	}

	var response []pgstore.WorkoutTemplateResponse
	for _, template := range templates {
		response = append(response, pgstore.WorkoutTemplateResponse{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			Thumbnail:   template.Thumbnail,
			Category:    template.Category,
			Difficulty:  template.Difficulty,
			Duration:    template.Duration,
			WeekDays:    template.WeekDays,
			Tags:        template.Tags,
		})
	}

	return response, nil
}

func (s *WorkoutService) CreateWorkoutTemplate(ctx context.Context, name, description, thumbnail, category, difficulty string, duration int, weekDays []string, exercises []map[string]interface{}, tags []string) (*pgstore.WorkoutTemplateResponse, error) {
	templateID := uuid.New()

	err := s.queries.CreateWorkoutTemplate(ctx, pgstore.CreateWorkoutTemplateParams{
		ID:          templateID,
		Name:        name,
		Description: description,
		Thumbnail:   thumbnail,
		Category:    category,
		Difficulty:  difficulty,
		Duration:    int32(duration),
		WeekDays:    weekDays,
		Tags:        tags,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create workout template: %w", err)
	}

	return &pgstore.WorkoutTemplateResponse{
		ID:          templateID,
		Name:        name,
		Description: description,
		Thumbnail:   thumbnail,
		Category:    category,
		Difficulty:  difficulty,
		Duration:    int32(duration),
		WeekDays:    weekDays,
		Tags:        tags,
	}, nil
}

func (s *WorkoutService) GetAllWorkoutTemplates(ctx context.Context) ([]pgstore.WorkoutTemplateResponse, error) {
	return s.GetWorkoutTemplates(ctx)
}

func (s *WorkoutService) DeleteWorkoutTemplate(ctx context.Context, templateID string) error {
	id, err := uuid.Parse(templateID)
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	err = s.queries.DeleteWorkoutTemplate(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete workout template: %w", err)
	}

	return nil
}

func (s *WorkoutService) GetAllPrograms(ctx context.Context, userID string) ([]pgstore.TrainingProgramResponse, error) {
	programs, err := s.queries.GetTrainingPrograms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get training programs: %w", err)
	}

	var response []pgstore.TrainingProgramResponse
	for _, program := range programs {
		response = append(response, pgstore.TrainingProgramResponse{
			ID:          program.ID,
			Name:        program.Name,
			Description: program.Description,
			Thumbnail:   program.Thumbnail,
			Category:    program.Category,
			Difficulty:  program.Difficulty,
			Duration:    program.Duration,
			IsFree:      program.IsFree,
			Price:       program.Price,
		})
	}

	return response, nil
}

func (s *WorkoutService) GetFreePrograms(ctx context.Context) ([]pgstore.TrainingProgramResponse, error) {
	programs, err := s.queries.GetFreeTrainingPrograms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get free training programs: %w", err)
	}

	var response []pgstore.TrainingProgramResponse
	for _, program := range programs {
		response = append(response, pgstore.TrainingProgramResponse{
			ID:          program.ID,
			Name:        program.Name,
			Description: program.Description,
			Thumbnail:   program.Thumbnail,
			Category:    program.Category,
			Difficulty:  program.Difficulty,
			Duration:    program.Duration,
			IsFree:      program.IsFree,
			Price:       program.Price,
		})
	}

	return response, nil
}

func (s *WorkoutService) GetFreeProgramByID(ctx context.Context, programID string) (*pgstore.TrainingProgramResponse, error) {
	id, err := uuid.Parse(programID)
	if err != nil {
		return nil, fmt.Errorf("invalid program ID: %w", err)
	}

	program, err := s.queries.GetTrainingProgramById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get training program: %w", err)
	}

	if !program.IsFree {
		return nil, fmt.Errorf("program is not free")
	}

	return &pgstore.TrainingProgramResponse{
		ID:          program.ID,
		Name:        program.Name,
		Description: program.Description,
		Thumbnail:   program.Thumbnail,
		Category:    program.Category,
		Difficulty:  program.Difficulty,
		Duration:    program.Duration,
		IsFree:      program.IsFree,
		Price:       program.Price,
	}, nil
}
