package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

type SchedulingService struct {
	queries *pgstore.Queries
}

func NewSchedulingService(queries *pgstore.Queries) *SchedulingService {
	return &SchedulingService{
		queries: queries,
	}
}

func (s *SchedulingService) GetSchedulings(ctx context.Context, userID uuid.UUID) ([]pgstore.Scheduling, error) {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return nil, fmt.Errorf("only personal trainers can view schedulings")
	}

	dbSchedulings, err := s.queries.GetSchedulings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedulings: %w", err)
	}

	schedulings := make([]pgstore.Scheduling, len(dbSchedulings))
	for i, dbScheduling := range dbSchedulings {
		schedulings[i] = s.convertSchedulingFromDB(dbScheduling)
	}

	return schedulings, nil
}

func (s *SchedulingService) GetSchedulingByID(ctx context.Context, schedulingID, userID uuid.UUID) (*pgstore.Scheduling, error) {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return nil, fmt.Errorf("only personal trainers can view schedulings")
	}

	dbScheduling, err := s.queries.GetSchedulingById(ctx, pgstore.GetSchedulingByIdParams{
		ID:         schedulingID,
		PersonalID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduling: %w", err)
	}
	if dbScheduling == nil {
		return nil, fmt.Errorf("scheduling not found")
	}

	scheduling := s.convertSchedulingFromDB(*dbScheduling)
	return &scheduling, nil
}

func (s *SchedulingService) CreateScheduling(ctx context.Context, userID uuid.UUID, req pgstore.CreateSchedulingParams) (*pgstore.Scheduling, error) {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return nil, fmt.Errorf("only personal trainers can create schedulings")
	}

	schedulingID := uuid.New()
	now := time.Now()

	_, err = s.queries.CreateScheduling(ctx, pgstore.CreateSchedulingParams{
		ID:         schedulingID,
		PersonalID: userID,
		StudentID:  req.StudentID,
		WorkoutID:  req.WorkoutID,
		Date:       req.Date,
		Type:       pgstore.SchedulingType(req.Type),
		Status:     pgstore.SchedulingStatus(req.Status),
		CreatedAt:  now,
		UserID:     &userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduling: %w", err)
	}

	return s.GetSchedulingByID(ctx, schedulingID, userID)
}

func (s *SchedulingService) UpdateSchedulingStatus(ctx context.Context, schedulingID, userID uuid.UUID, status pgstore.SchedulingStatus) error {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return fmt.Errorf("only personal trainers can update schedulings")
	}

	err = s.queries.UpdateSchedulingStatus(ctx, pgstore.UpdateSchedulingStatusParams{
		ID:     schedulingID,
		Status: pgstore.SchedulingStatus(status),
	})
	if err != nil {
		return fmt.Errorf("failed to update scheduling status: %w", err)
	}

	return nil
}

func (s *SchedulingService) StartScheduling(ctx context.Context, schedulingID, userID uuid.UUID) error {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return fmt.Errorf("only personal trainers can start schedulings")
	}

	now := time.Now()
	err = s.queries.UpdateSchedulingWithStartTime(ctx, pgstore.UpdateSchedulingWithStartTimeParams{
		ID:        schedulingID,
		StartedAt: &now,
		Status:    pgstore.SchedulingStatusInProgress,
	})
	if err != nil {
		return fmt.Errorf("failed to start scheduling: %w", err)
	}

	return nil
}

func (s *SchedulingService) CompleteScheduling(ctx context.Context, schedulingID, userID uuid.UUID) error {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return fmt.Errorf("only personal trainers can complete schedulings")
	}

	now := time.Now()
	err = s.queries.UpdateSchedulingWithCompletedTime(ctx, pgstore.UpdateSchedulingWithCompletedTimeParams{
		ID:          schedulingID,
		CompletedAt: &now,
		Status:      pgstore.SchedulingStatusCompleted,
	})
	if err != nil {
		return fmt.Errorf("failed to complete scheduling: %w", err)
	}

	return nil
}

func (s *SchedulingService) UpdateScheduling(ctx context.Context, schedulingID, userID uuid.UUID, req pgstore.UpdateSchedulingRequest) error {
	if req.Status != nil {
		return s.UpdateSchedulingStatus(ctx, schedulingID, userID, *req.Status)
	}
	return nil
}

func (s *SchedulingService) CancelScheduling(ctx context.Context, schedulingID, userID uuid.UUID, reason string) error {
	// Verify user is a personal trainer
	role, err := s.queries.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}
	if role != pgstore.RolePersonal {
		return fmt.Errorf("only personal trainers can cancel schedulings")
	}

	// Update scheduling status
	err = s.queries.UpdateSchedulingWithCanceledTime(ctx, pgstore.UpdateSchedulingWithCanceledTimeParams{
		ID:     schedulingID,
		Status: pgstore.SchedulingStatusCanceled,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel scheduling: %w", err)
	}

	// Create history entry
	historyID := uuid.New()
	now := time.Now()
	_, err = s.queries.CreateSchedulingHistory(ctx, pgstore.CreateSchedulingHistoryParams{
		ID:         historyID,
		ScheduleID: schedulingID,
		UserID:     userID,
		Status:     pgstore.SchedulingStatusCanceled,
		ChangedAt:  &now,
		ChangedBy:  "personal_trainer",
		Reason:     &reason,
	})
	if err != nil {
		return fmt.Errorf("failed to create scheduling history: %w", err)
	}

	return nil
}

// Helper function to convert between database and domain models
func (s *SchedulingService) convertSchedulingFromDB(dbScheduling pgstore.Scheduling) pgstore.Scheduling {
	return pgstore.Scheduling{
		ID:          dbScheduling.ID,
		PersonalID:  dbScheduling.PersonalID,
		StudentID:   dbScheduling.StudentID,
		WorkoutID:   dbScheduling.WorkoutID,
		Date:        dbScheduling.Date,
		Type:        pgstore.SchedulingType(dbScheduling.Type),
		Status:      pgstore.SchedulingStatus(dbScheduling.Status),
		StartedAt:   dbScheduling.StartedAt,
		CompletedAt: dbScheduling.CompletedAt,
		CreatedAt:   dbScheduling.CreatedAt,
		UserID:      dbScheduling.UserID,
	}
}
