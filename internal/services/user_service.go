package services

import (
	"context"
	"fmt"

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
	// For now, return mock response
	// TODO: Implement actual student creation with user
	return &pgstore.UserResponse{
		ID:    uuid.New(),
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Role:  pgstore.RoleStudent,
	}, nil
}

func (s *UserService) CreatePersonalWithUser(ctx context.Context, req pgstore.CreatePersonalWithUserRequest) (*pgstore.UserResponse, error) {
	// For now, return mock response
	// TODO: Implement actual personal trainer creation with user
	return &pgstore.UserResponse{
		ID:    uuid.New(),
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Role:  pgstore.RolePersonal,
	}, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *pgstore.UpdateProfileRequest) error {
	// For now, just return success
	// TODO: Implement actual profile update
	return nil
}

// Placeholder methods for features not yet implemented

func (s *UserService) GetAllUsers(page, limit, role, search string) (interface{}, interface{}, error) {
	// For now, return empty results
	// TODO: Implement actual user listing with pagination
	return []interface{}{}, 0, nil
}

func (s *UserService) UpdateUserStatus(userID, status, reason string) error {
	// For now, just return success
	// TODO: Implement actual user status update
	return nil
}

func (s *UserService) DeleteUser(userID, reason string) error {
	// For now, just return success
	// TODO: Implement actual user deletion
	return nil
}
