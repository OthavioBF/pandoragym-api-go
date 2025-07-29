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

func (s *SchedulingService) GetSchedulings(ctx context.Context, userID uuid.UUID) ([]pgstore.SchedulingResponse, error) {
	return []pgstore.SchedulingResponse{
		{
			ID:         uuid.New(),
			PersonalID: uuid.New(),
			StudentID:  userID,
			Date:       time.Now().Add(24 * time.Hour),
			Type:       pgstore.SchedulingTypeInPerson,
			Status:     pgstore.SchedulingStatusScheduled,
			WorkoutID:  &[]uuid.UUID{uuid.New()}[0],
			CreatedAt:  time.Now().AddDate(0, 0, -1),
		},
		{
			ID:         uuid.New(),
			PersonalID: uuid.New(),
			StudentID:  userID,
			Date:       time.Now().Add(72 * time.Hour),
			Type:       pgstore.SchedulingTypeOnline,
			Status:     pgstore.SchedulingStatusPendingConfirmation,
			CreatedAt:  time.Now(),
		},
	}, nil
}

func (s *SchedulingService) GetSchedulingByID(ctx context.Context, schedulingID, userID uuid.UUID) (*pgstore.SchedulingResponse, error) {
	return &pgstore.SchedulingResponse{
		ID:         schedulingID,
		PersonalID: uuid.New(),
		StudentID:  userID,
		Date:       time.Now().Add(24 * time.Hour),
		Type:       pgstore.SchedulingTypeInPerson,
		Status:     pgstore.SchedulingStatusScheduled,
		WorkoutID:  &[]uuid.UUID{uuid.New()}[0],
		CreatedAt:  time.Now().AddDate(0, 0, -1),
	}, nil
}

func (s *SchedulingService) CreateScheduling(ctx context.Context, req pgstore.CreateSchedulingRequest, userID uuid.UUID) (*pgstore.SchedulingResponse, error) {
	if req.Date.Before(time.Now()) {
		return nil, fmt.Errorf("cannot schedule in the past")
	}

	schedulingID := uuid.New()

	return &pgstore.SchedulingResponse{
		ID:         schedulingID,
		PersonalID: req.PersonalID,
		StudentID:  userID,
		Date:       req.Date,
		Type:       req.Type,
		Status:     pgstore.SchedulingStatusPendingConfirmation,
		CreatedAt:  time.Now(),
	}, nil
}

func (s *SchedulingService) UpdateScheduling(ctx context.Context, schedulingID uuid.UUID, req pgstore.UpdateSchedulingRequest, userID uuid.UUID) (*pgstore.SchedulingResponse, error) {
	if req.Date != nil && req.Date.Before(time.Now()) {
		return nil, fmt.Errorf("cannot reschedule to the past")
	}

	return &pgstore.SchedulingResponse{
		ID:         schedulingID,
		PersonalID: uuid.New(),
		StudentID:  userID,
		Date:       *req.Date,
		Type:       *req.Type,
		Status:     *req.Status,
		CreatedAt:  time.Now().AddDate(0, 0, -1),
	}, nil
}

func (s *SchedulingService) CancelScheduling(ctx context.Context, schedulingID uuid.UUID, reason string, userID uuid.UUID) error {
	return nil
}

func (s *SchedulingService) ConfirmScheduling(ctx context.Context, schedulingID, userID uuid.UUID) error {
	return nil
}

func (s *SchedulingService) CompleteScheduling(ctx context.Context, schedulingID, userID uuid.UUID, notes string) error {
	return nil
}

func (s *SchedulingService) GetTrainerAvailability(ctx context.Context, trainerID uuid.UUID, date time.Time) ([]pgstore.AvailabilitySlot, error) {
	return []pgstore.AvailabilitySlot{
		{
			StartTime: time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, date.Location()),
			Available: true,
		},
		{
			StartTime: time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day(), 11, 0, 0, 0, date.Location()),
			Available: false,
		},
		{
			StartTime: time.Date(date.Year(), date.Month(), date.Day(), 14, 0, 0, 0, date.Location()),
			EndTime:   time.Date(date.Year(), date.Month(), date.Day(), 15, 0, 0, 0, date.Location()),
			Available: true,
		},
	}, nil
}

func (s *SchedulingService) SetTrainerAvailability(ctx context.Context, trainerID uuid.UUID, date time.Time, slots []pgstore.AvailabilitySlot) error {
	return nil
}

func (s *SchedulingService) GetSchedulingHistory(ctx context.Context, userID uuid.UUID) ([]pgstore.SchedulingHistoryResponse, error) {
	return []pgstore.SchedulingHistoryResponse{
		{
			ID:            uuid.New(),
			Date:          time.Now().AddDate(0, 0, -7),
			Type:          pgstore.SchedulingTypeInPerson,
			Status:        pgstore.SchedulingStatusCompleted,
			WorkoutName:   &[]string{"Push Day Workout"}[0],
			PartnerName:   "John Smith",
			CompletedAt:   &[]time.Time{time.Now().AddDate(0, 0, -7)}[0],
		},
		{
			ID:            uuid.New(),
			Date:          time.Now().AddDate(0, 0, -14),
			Type:          pgstore.SchedulingTypeOnline,
			Status:        pgstore.SchedulingStatusCompleted,
			WorkoutName:   &[]string{"Leg Day Workout"}[0],
			PartnerName:   "John Smith",
			CompletedAt:   &[]time.Time{time.Now().AddDate(0, 0, -14)}[0],
		},
		{
			ID:            uuid.New(),
			Date:          time.Now().AddDate(0, 0, -21),
			Type:          pgstore.SchedulingTypeInPerson,
			Status:        pgstore.SchedulingStatusCanceled,
			WorkoutName:   nil,
			PartnerName:   "Sarah Johnson",
			CompletedAt:   nil,
		},
	}, nil
}
