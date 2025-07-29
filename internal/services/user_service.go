package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*pgstore.UserResponse, error) {
	user, err := s.queries.GetUserById(ctx, pgstore.GetUserByIdParams{ID: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, nil
	}

	return &pgstore.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) CreateStudentWithUser(ctx context.Context, req pgstore.CreateStudentWithUserRequest) (*pgstore.UserResponse, error) {
	authService := NewAuthService(s.queries, "temp-secret")
	return authService.CreateStudentWithUser(ctx, req)
}

func (s *UserService) CreatePersonalWithUser(ctx context.Context, req pgstore.CreatePersonalWithUserRequest) (*pgstore.UserResponse, error) {
	authService := NewAuthService(s.queries, "temp-secret")
	return authService.CreatePersonalWithUser(ctx, req)
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *pgstore.UpdateProfileRequest) error {
	updateParams := pgstore.UpdateUserParams{
		ID:        userID,
		UpdatedAt: time.Now(),
	}

	if req.Name != nil {
		updateParams.Name = req.Name
	}
	if req.Phone != nil {
		updateParams.Phone = req.Phone
	}

	s.queries.UpdateUser(ctx, updateParams)
	return nil
}

func (s *UserService) UpdateUserAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error {
	s.queries.UpdateUserAvatar(ctx, pgstore.UpdateUserAvatarParams{
		ID:        userID,
		AvatarURL: avatarURL,
		UpdatedAt: time.Now(),
	})
	return nil
}

func (s *UserService) GetAllUsers(page, limit, role, search string) (interface{}, interface{}, error) {
	pageInt := 1
	limitInt := 20

	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageInt = p
		}
	}

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitInt = l
		}
	}

	mockUsers := []map[string]interface{}{
		{
			"id":         uuid.New().String(),
			"name":       "John Doe",
			"email":      "john@example.com",
			"phone":      "+1234567890",
			"role":       "STUDENT",
			"created_at": time.Now().AddDate(0, -2, 0),
			"updated_at": time.Now(),
		},
		{
			"id":         uuid.New().String(),
			"name":       "Jane Smith",
			"email":      "jane@example.com",
			"phone":      "+1234567891",
			"role":       "PERSONAL",
			"created_at": time.Now().AddDate(0, -1, 0),
			"updated_at": time.Now(),
		},
	}

	return mockUsers, len(mockUsers), nil
}

func (s *UserService) UpdateUserStatus(userID, status, reason string) error {
	return nil
}

func (s *UserService) DeleteUser(userID, reason string) error {
	return nil
}

func (s *UserService) GetPersonalTrainers(ctx context.Context) ([]pgstore.PersonalTrainerResponse, error) {
	return []pgstore.PersonalTrainerResponse{
		{
			ID:             uuid.New(),
			Name:           "John Smith",
			Email:          "john.smith@example.com",
			Phone:          "+1234567890",
			AvatarURL:      &[]string{"https://example.com/avatar1.jpg"}[0],
			Description:    &[]string{"Certified personal trainer with 5 years of experience"}[0],
			VideoURL:       &[]string{"https://example.com/video1.mp4"}[0],
			Experience:     &[]string{"5 years"}[0],
			Specialization: &[]string{"Weight Loss, Muscle Building"}[0],
			Qualifications: &[]string{"NASM-CPT, ACSM-CPT"}[0],
			Rating:         &[]float64{4.8}[0],
			ReviewCount:    25,
		},
		{
			ID:             uuid.New(),
			Name:           "Sarah Johnson",
			Email:          "sarah.johnson@example.com",
			Phone:          "+1234567891",
			AvatarURL:      &[]string{"https://example.com/avatar2.jpg"}[0],
			Description:    &[]string{"Specialized in functional fitness and rehabilitation"}[0],
			VideoURL:       &[]string{"https://example.com/video2.mp4"}[0],
			Experience:     &[]string{"8 years"}[0],
			Specialization: &[]string{"Functional Fitness, Rehabilitation"}[0],
			Qualifications: &[]string{"NSCA-CSCS, FMS"}[0],
			Rating:         &[]float64{4.9}[0],
			ReviewCount:    42,
		},
	}, nil
}

func (s *UserService) GetPersonalTrainerByID(ctx context.Context, trainerID uuid.UUID) (*pgstore.PersonalTrainerResponse, error) {
	return &pgstore.PersonalTrainerResponse{
		ID:             trainerID,
		Name:           "John Smith",
		Email:          "john.smith@example.com",
		Phone:          "+1234567890",
		AvatarURL:      &[]string{"https://example.com/avatar1.jpg"}[0],
		Description:    &[]string{"Certified personal trainer with 5 years of experience"}[0],
		VideoURL:       &[]string{"https://example.com/video1.mp4"}[0],
		Experience:     &[]string{"5 years"}[0],
		Specialization: &[]string{"Weight Loss, Muscle Building"}[0],
		Qualifications: &[]string{"NASM-CPT, ACSM-CPT"}[0],
		Rating:         &[]float64{4.8}[0],
		ReviewCount:    25,
	}, nil
}

func (s *UserService) AddTrainerComment(ctx context.Context, trainerID, studentID uuid.UUID, comment string, rating int) error {
	return nil
}

func (s *UserService) GetTrainerComments(ctx context.Context, trainerID uuid.UUID) ([]pgstore.TrainerCommentResponse, error) {
	return []pgstore.TrainerCommentResponse{
		{
			ID:          uuid.New(),
			StudentName: "Alice Johnson",
			Comment:     "Great trainer! Very knowledgeable and motivating.",
			Rating:      5,
			CreatedAt:   time.Now().AddDate(0, 0, -7),
		},
		{
			ID:          uuid.New(),
			StudentName: "Bob Wilson",
			Comment:     "Helped me achieve my fitness goals. Highly recommended!",
			Rating:      4,
			CreatedAt:   time.Now().AddDate(0, 0, -14),
		},
	}, nil
}

func (s *UserService) GetTrainerStudents(ctx context.Context, trainerID uuid.UUID) ([]pgstore.StudentResponse, error) {
	return []pgstore.StudentResponse{
		{
			ID:                    uuid.New(),
			Name:                  "Alice Johnson",
			Email:                 "alice@example.com",
			Phone:                 "+1234567892",
			AvatarURL:             &[]string{"https://example.com/avatar3.jpg"}[0],
			BornDate:              time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
			Age:                   33,
			Weight:                65.5,
			Objective:             "Weight Loss",
			TrainingFrequency:     "3 times per week",
			DidBodybuilding:       false,
			MedicalCondition:      nil,
			PhysicalActivityLevel: &[]string{"Beginner"}[0],
			Observations:          &[]string{"Prefers morning workouts"}[0],
		},
		{
			ID:                    uuid.New(),
			Name:                  "Bob Wilson",
			Email:                 "bob@example.com",
			Phone:                 "+1234567893",
			AvatarURL:             nil,
			BornDate:              time.Date(1985, 8, 22, 0, 0, 0, 0, time.UTC),
			Age:                   38,
			Weight:                80.0,
			Objective:             "Muscle Building",
			TrainingFrequency:     "4 times per week",
			DidBodybuilding:       true,
			MedicalCondition:      &[]string{"Previous knee injury"}[0],
			PhysicalActivityLevel: &[]string{"Intermediate"}[0],
			Observations:          nil,
		},
	}, nil
}

func (s *UserService) AssignStudentToTrainer(ctx context.Context, studentID, trainerID uuid.UUID) error {
	return nil
}

func (s *UserService) RemoveStudentFromTrainer(ctx context.Context, studentID, trainerID uuid.UUID) error {
	return nil
}

func (s *UserService) TrainerHasAccessToStudent(ctx context.Context, trainerID, studentID uuid.UUID) (bool, error) {
	return true, nil
}
