package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
)

type UserService struct {
	queries *pgstore.Queries
}

func NewUserService(queries *pgstore.Queries) *UserService {
	return &UserService{
		queries: queries,
	}
}

// GetUserByID returns a clean domain model
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*pgstore.UserResponse, error) {
	user, err := s.queries.GetUserById(ctx, pgstore.GetUserByIdParams{ID: id})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	// Convert database model to response model
	userResponse := &pgstore.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      pgstore.Role(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Handle nullable avatar URL
	if user.AvatarURL != nil {
		userResponse.AvatarURL = user.AvatarURL
	}

	// Add personal trainer specific fields if applicable
	if user.Role == pgstore.RolePersonal {
		personal := &pgstore.Personal{
			ID: user.ID,
		}
		if user.Rating != nil {
			personal.Rating = user.Rating
		}
		if user.Description != nil {
			personal.Description = user.Description
		}
		if user.VideoURL != nil {
			personal.VideoURL = user.VideoURL
		}
		if user.Experience != nil {
			personal.Experience = user.Experience
		}
		if user.Specialization != nil {
			personal.Specialization = user.Specialization
		}
		if user.Qualifications != nil {
			personal.Qualifications = user.Qualifications
		}
		userResponse.Personal = personal
	}

	// Add student specific fields if applicable
	if user.Role == pgstore.RoleStudent {
		student := &pgstore.Student{
			ID: user.ID,
		}
		if user.BornDate != nil {
			student.BornDate = *user.BornDate
		}
		if user.Age != nil {
			student.Age = *user.Age
		}
		if user.Weight != nil {
			student.Weight = *user.Weight
		}
		if user.Objective != nil {
			student.Objective = *user.Objective
		}
		if user.TrainingFrequency != nil {
			student.TrainingFrequency = *user.TrainingFrequency
		}
		if user.DidBodybuilding != nil {
			student.DidBodybuilding = *user.DidBodybuilding
		}
		if user.MedicalCondition != nil {
			student.MedicalCondition = user.MedicalCondition
		}
		if user.PhysicalActivityLevel != nil {
			student.PhysicalActivityLevel = user.PhysicalActivityLevel
		}
		if user.Observations != nil {
			student.Observations = user.Observations
		}
		userResponse.Student = student
	}

	return userResponse, nil
}

// GetUserByEmail returns user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*pgstore.UserResponse, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	userResponse := &pgstore.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      pgstore.Role(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.AvatarURL != nil {
		userResponse.AvatarURL = user.AvatarURL
	}

	return userResponse, nil
}

// UpdateUserProfile updates user profile information
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *pgstore.UpdateProfileRequest) error {
	params := pgstore.UpdateUserProfileParams{
		ID:        userID,
		UpdatedAt: time.Now(),
	}

	if req.Name != nil {
		params.Name = req.Name
	}

	if req.Phone != nil {
		params.Phone = req.Phone
	}

	return s.queries.UpdateUserProfile(ctx, params)
}

// CreateStudent creates a new student user
func (s *UserService) CreateStudent(ctx context.Context, req pgstore.CreateStudentWithUserRequest) (*pgstore.UserResponse, error) {
	authService := NewAuthService(s.queries, "temp-secret") // TODO: Get from config
	return authService.CreateStudent(ctx, &req)
}

// CreatePersonal creates a new personal trainer user
func (s *UserService) CreatePersonal(ctx context.Context, req pgstore.CreatePersonalWithUserRequest) (*pgstore.UserResponse, error) {
	authService := NewAuthService(s.queries, "temp-secret") // TODO: Get from config
	return authService.CreatePersonal(ctx, &req)
}
