package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

type PlanService struct {
	queries *pgstore.Queries
}

func NewPlanService(queries *pgstore.Queries) *PlanService {
	return &PlanService{
		queries: queries,
	}
}

func (s *PlanService) GetTrainerPlans(ctx context.Context, trainerID uuid.UUID) ([]pgstore.PlanResponse, error) {
	return []pgstore.PlanResponse{
		{
			ID:              uuid.New(),
			TrainerID:       trainerID,
			Name:            "Basic Training Plan",
			Description:     "A comprehensive training plan for beginners",
			Price:           99.99,
			Duration:        30,
			Features:        []string{"3 workouts per week", "Nutrition guidance", "Progress tracking"},
			IsActive:        true,
			SubscriberCount: 15,
			CreatedAt:       time.Now().AddDate(0, -1, 0),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			TrainerID:       trainerID,
			Name:            "Advanced Training Plan",
			Description:     "Intensive training for experienced athletes",
			Price:           149.99,
			Duration:        30,
			Features:        []string{"5 workouts per week", "Custom meal plans", "1-on-1 sessions", "24/7 support"},
			IsActive:        true,
			SubscriberCount: 8,
			CreatedAt:       time.Now().AddDate(0, -2, 0),
			UpdatedAt:       time.Now(),
		},
	}, nil
}

func (s *PlanService) CreatePlan(ctx context.Context, trainerID uuid.UUID, name, description string, price float64, duration int, features []string) (*pgstore.PlanResponse, error) {
	planID := uuid.New()

	return &pgstore.PlanResponse{
		ID:              planID,
		TrainerID:       trainerID,
		Name:            name,
		Description:     description,
		Price:           price,
		Duration:        int32(duration),
		Features:        features,
		IsActive:        true,
		SubscriberCount: 0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func (s *PlanService) UpdatePlan(ctx context.Context, trainerID uuid.UUID, planID, name, description string, price float64, duration int, features []string) (*pgstore.PlanResponse, error) {
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan ID: %w", err)
	}

	return &pgstore.PlanResponse{
		ID:              planUUID,
		TrainerID:       trainerID,
		Name:            name,
		Description:     description,
		Price:           price,
		Duration:        int32(duration),
		Features:        features,
		IsActive:        true,
		SubscriberCount: 5,
		CreatedAt:       time.Now().AddDate(0, -1, 0),
		UpdatedAt:       time.Now(),
	}, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, trainerID uuid.UUID, planID string) error {
	return nil
}

func (s *PlanService) GetAllPlans(ctx context.Context) ([]pgstore.PlanResponse, error) {
	return []pgstore.PlanResponse{
		{
			ID:              uuid.New(),
			TrainerID:       uuid.New(),
			TrainerName:     "John Smith",
			Name:            "Weight Loss Program",
			Description:     "Effective weight loss program with cardio and strength training",
			Price:           79.99,
			Duration:        30,
			Features:        []string{"4 workouts per week", "Meal planning", "Progress tracking"},
			IsActive:        true,
			SubscriberCount: 25,
			CreatedAt:       time.Now().AddDate(0, -2, 0),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			TrainerID:       uuid.New(),
			TrainerName:     "Sarah Johnson",
			Name:            "Muscle Building Plan",
			Description:     "Comprehensive muscle building program for all levels",
			Price:           119.99,
			Duration:        60,
			Features:        []string{"5 workouts per week", "Supplement guidance", "Form coaching", "Monthly check-ins"},
			IsActive:        true,
			SubscriberCount: 18,
			CreatedAt:       time.Now().AddDate(0, -3, 0),
			UpdatedAt:       time.Now(),
		},
	}, nil
}

func (s *PlanService) GetPlanByID(ctx context.Context, planID string) (*pgstore.PlanResponse, error) {
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan ID: %w", err)
	}

	return &pgstore.PlanResponse{
		ID:              planUUID,
		TrainerID:       uuid.New(),
		TrainerName:     "John Smith",
		Name:            "Weight Loss Program",
		Description:     "Effective weight loss program with cardio and strength training",
		Price:           79.99,
		Duration:        30,
		Features:        []string{"4 workouts per week", "Meal planning", "Progress tracking"},
		IsActive:        true,
		SubscriberCount: 25,
		CreatedAt:       time.Now().AddDate(0, -2, 0),
		UpdatedAt:       time.Now(),
	}, nil
}

func (s *PlanService) SubscribeToPlan(ctx context.Context, userID uuid.UUID, planID string) (*pgstore.SubscriptionResponse, error) {
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan ID: %w", err)
	}

	subscriptionID := uuid.New()
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 30)

	return &pgstore.SubscriptionResponse{
		ID:        subscriptionID,
		UserID:    userID,
		PlanID:    planUUID,
		PlanName:  "Weight Loss Program",
		StartDate: startDate,
		EndDate:   endDate,
		Status:    pgstore.SubscriptionStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *PlanService) CancelSubscription(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (s *PlanService) GetUserSubscription(ctx context.Context, userID uuid.UUID) (*pgstore.SubscriptionResponse, error) {
	return &pgstore.SubscriptionResponse{
		ID:        uuid.New(),
		UserID:    userID,
		PlanID:    uuid.New(),
		PlanName:  "Weight Loss Program",
		StartDate: time.Now().AddDate(0, 0, -15),
		EndDate:   time.Now().AddDate(0, 0, 15),
		Status:    pgstore.SubscriptionStatusActive,
		CreatedAt: time.Now().AddDate(0, 0, -15),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *PlanService) GetSubscriptionHistory(ctx context.Context, userID uuid.UUID) ([]pgstore.SubscriptionHistoryResponse, error) {
	return []pgstore.SubscriptionHistoryResponse{
		{
			ID:          uuid.New(),
			PlanName:    "Basic Training Plan",
			TrainerName: "John Smith",
			StartDate:   time.Now().AddDate(0, -2, 0),
			EndDate:     time.Now().AddDate(0, -1, 0),
			Status:      pgstore.SubscriptionStatusExpired,
			Amount:      99.99,
			CreatedAt:   time.Now().AddDate(0, -2, 0),
		},
		{
			ID:          uuid.New(),
			PlanName:    "Weight Loss Program",
			TrainerName: "Sarah Johnson",
			StartDate:   time.Now().AddDate(0, 0, -15),
			EndDate:     time.Now().AddDate(0, 0, 15),
			Status:      pgstore.SubscriptionStatusActive,
			Amount:      79.99,
			CreatedAt:   time.Now().AddDate(0, 0, -15),
		},
	}, nil
}

func (s *PlanService) ProcessPayment(ctx context.Context, subscriptionID uuid.UUID, paymentMethod, paymentToken string) error {
	return nil
}

func (s *PlanService) GetPlanSubscribers(ctx context.Context, trainerID uuid.UUID, planID string) ([]pgstore.SubscriberResponse, error) {
	return []pgstore.SubscriberResponse{
		{
			UserID:       uuid.New(),
			UserName:     "Alice Johnson",
			UserEmail:    "alice@example.com",
			StartDate:    time.Now().AddDate(0, 0, -10),
			EndDate:      time.Now().AddDate(0, 0, 20),
			Status:       pgstore.SubscriptionStatusActive,
			SubscribedAt: time.Now().AddDate(0, 0, -10),
		},
		{
			UserID:       uuid.New(),
			UserName:     "Bob Wilson",
			UserEmail:    "bob@example.com",
			StartDate:    time.Now().AddDate(0, 0, -5),
			EndDate:      time.Now().AddDate(0, 0, 25),
			Status:       pgstore.SubscriptionStatusActive,
			SubscribedAt: time.Now().AddDate(0, 0, -5),
		},
	}, nil
}

func (s *PlanService) GetPlanRevenue(ctx context.Context, trainerID uuid.UUID, planID string) (*pgstore.PlanRevenueResponse, error) {
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan ID: %w", err)
	}

	return &pgstore.PlanRevenueResponse{
		PlanID:            planUUID,
		PlanName:          "Weight Loss Program",
		TotalRevenue:      1599.80,
		MonthlyRevenue:    479.94,
		ActiveSubscribers: 6,
		TotalSubscribers:  20,
		AverageRevenue:    79.99,
	}, nil
}
