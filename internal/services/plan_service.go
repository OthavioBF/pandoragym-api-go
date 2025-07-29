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

func (s *PlanService) GetTrainerPlans(ctx context.Context, trainerID uuid.UUID) (interface{}, error) {
	// Verify trainer exists
	trainer, err := s.queries.GetUserById(ctx, pgstore.GetUserByIdParams{ID: trainerID})
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer: %w", err)
	}

	if trainer == nil || trainer.Role != pgstore.RolePersonal {
		return nil, fmt.Errorf("trainer not found or invalid role")
	}

	// For now, return mock data until we implement plan management in database
	// TODO: Implement actual plan storage and retrieval
	return []map[string]interface{}{
		{
			"id":          "plan-1",
			"trainer_id":  trainerID.String(),
			"name":        "Basic Training Plan",
			"description": "Perfect for beginners looking to start their fitness journey",
			"price":       99.99,
			"duration":    30, // days
			"features": []string{
				"3 workouts per week",
				"Basic nutrition guide",
				"Email support",
			},
			"active":      true,
			"subscribers": 15,
			"created_at":  time.Now().AddDate(0, -2, 0).Format("2006-01-02T15:04:05Z"),
		},
		{
			"id":          "plan-2",
			"trainer_id":  trainerID.String(),
			"name":        "Premium Training Plan",
			"description": "Advanced training with personalized attention",
			"price":       199.99,
			"duration":    30,
			"features": []string{
				"5 workouts per week",
				"Custom meal plans",
				"Video call sessions",
				"24/7 chat support",
			},
			"active":      true,
			"subscribers": 8,
			"created_at":  time.Now().AddDate(0, -1, 0).Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

func (s *PlanService) CreatePlan(ctx context.Context, trainerID uuid.UUID, name, description string, price float64, duration int, features []string) (interface{}, error) {
	// Verify trainer exists
	trainer, err := s.queries.GetUserById(ctx, pgstore.GetUserByIdParams{ID: trainerID})
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer: %w", err)
	}

	if trainer == nil || trainer.Role != pgstore.RolePersonal {
		return nil, fmt.Errorf("trainer not found or invalid role")
	}

	// For now, return mock data until we implement plan storage
	// TODO: Implement actual plan creation in database
	plan := map[string]interface{}{
		"id":          uuid.New().String(),
		"trainer_id":  trainerID.String(),
		"name":        name,
		"description": description,
		"price":       price,
		"duration":    duration,
		"features":    features,
		"active":      true,
		"subscribers": 0,
		"created_at":  time.Now().Format("2006-01-02T15:04:05Z"),
	}

	return plan, nil
}

func (s *PlanService) UpdatePlan(ctx context.Context, trainerID uuid.UUID, planID string, name, description string, price float64, duration int, features []string) (interface{}, error) {
	// Verify trainer exists
	trainer, err := s.queries.GetUserById(ctx, pgstore.GetUserByIdParams{ID: trainerID})
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer: %w", err)
	}

	if trainer == nil || trainer.Role != pgstore.RolePersonal {
		return nil, fmt.Errorf("trainer not found or invalid role")
	}

	// For now, return mock data until we implement plan storage
	// TODO: Implement actual plan update in database
	plan := map[string]interface{}{
		"id":          planID,
		"trainer_id":  trainerID.String(),
		"name":        name,
		"description": description,
		"price":       price,
		"duration":    duration,
		"features":    features,
		"active":      true,
		"subscribers": 5, // Mock existing subscribers
		"updated_at":  time.Now().Format("2006-01-02T15:04:05Z"),
	}

	return plan, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, trainerID uuid.UUID, planID string) error {
	// Verify trainer exists
	trainer, err := s.queries.GetUserById(ctx, pgstore.GetUserByIdParams{ID: trainerID})
	if err != nil {
		return fmt.Errorf("failed to get trainer: %w", err)
	}

	if trainer == nil || trainer.Role != pgstore.RolePersonal {
		return fmt.Errorf("trainer not found or invalid role")
	}

	// For now, just return success until we implement plan storage
	// TODO: Implement actual plan deletion in database
	// In real implementation, check if trainer owns the plan and if there are active subscribers
	return nil
}

func (s *PlanService) GetPlanByID(ctx context.Context, planID string) (interface{}, error) {
	// For now, return mock data until we implement plan storage
	// TODO: Implement actual plan retrieval from database
	return map[string]interface{}{
		"id":          planID,
		"trainer_id":  "trainer-1",
		"name":        "Premium Training Plan",
		"description": "Advanced training with personalized attention",
		"price":       199.99,
		"duration":    30,
		"features": []string{
			"5 workouts per week",
			"Custom meal plans",
			"Video call sessions",
			"24/7 chat support",
		},
		"active":      true,
		"subscribers": 8,
		"trainer": map[string]interface{}{
			"id":           "trainer-1",
			"name":         "John Doe",
			"bio":          "Certified personal trainer with 5 years of experience",
			"specialties":  []string{"Weight Loss", "Muscle Building", "Strength Training"},
		},
	}, nil
}
