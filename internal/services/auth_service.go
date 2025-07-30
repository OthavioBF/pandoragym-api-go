package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries        *pgstore.Queries
	sessionManager *scs.SessionManager
}

type SessionData struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

func NewAuthService(queries *pgstore.Queries, sessionManager *scs.SessionManager) *AuthService {
	return &AuthService{
		queries:        queries,
		sessionManager: sessionManager,
	}
}

func (s *AuthService) AuthenticateWithPassword(ctx context.Context, email, password string) (*pgstore.UserResponse, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !s.verifyPassword(password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	s.sessionManager.Put(ctx, "user_id", user.ID.String())
	s.sessionManager.Put(ctx, "role", string(user.Role))
	s.sessionManager.Put(ctx, "email", user.Email)
	s.sessionManager.Put(ctx, "name", user.Name)

	userResponse := &pgstore.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userResponse, nil
}

func (s *AuthService) GetSessionData(ctx context.Context) (*SessionData, error) {
	userID := s.sessionManager.GetString(ctx, "user_id")
	if userID == "" {
		return nil, fmt.Errorf("no active session")
	}

	role := s.sessionManager.GetString(ctx, "role")
	email := s.sessionManager.GetString(ctx, "email")
	name := s.sessionManager.GetString(ctx, "name")

	return &SessionData{
		UserID: userID,
		Role:   role,
		Email:  email,
		Name:   name,
	}, nil
}

func (s *AuthService) IsAuthenticated(ctx context.Context) bool {
	userID := s.sessionManager.GetString(ctx, "user_id")
	return userID != ""
}

func (s *AuthService) GetUserIDFromSession(ctx context.Context) (uuid.UUID, error) {
	userIDStr := s.sessionManager.GetString(ctx, "user_id")
	if userIDStr == "" {
		return uuid.Nil, fmt.Errorf("no user in session")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in session: %w", err)
	}

	return userID, nil
}

func (s *AuthService) GetUserRoleFromSession(ctx context.Context) (string, error) {
	role := s.sessionManager.GetString(ctx, "role")
	if role == "" {
		return "", fmt.Errorf("no role in session")
	}

	return role, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	return s.sessionManager.Destroy(ctx)
}

func (s *AuthService) RefreshSession(ctx context.Context) error {
	if !s.IsAuthenticated(ctx) {
		return fmt.Errorf("no active session to refresh")
	}

	return s.sessionManager.RenewToken(ctx)
}

func (s *AuthService) CreateStudentWithUser(ctx context.Context, req pgstore.CreateStudentWithUserRequest) (*pgstore.UserResponse, error) {
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New()

	_, err = s.queries.CreateUser(ctx, pgstore.CreateUserParams{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
		Role:      pgstore.RoleStudent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.queries.CreateStudent(ctx, pgstore.CreateStudentParams{
		ID:                    userID,
		BornDate:              req.BornDate,
		Age:                   req.Age,
		Weight:                req.Weight,
		Objective:             req.Objective,
		TrainingFrequency:     req.TrainingFrequency,
		DidBodybuilding:       req.DidBodybuilding,
		MedicalCondition:      req.MedicalCondition,
		PhysicalActivityLevel: req.PhysicalActivityLevel,
		Observations:          req.Observations,
	})

	return &pgstore.UserResponse{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Role:      pgstore.RoleStudent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *AuthService) CreatePersonalWithUser(ctx context.Context, req pgstore.CreatePersonalWithUserRequest) (*pgstore.UserResponse, error) {
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New()

	_, err = s.queries.CreateUser(ctx, pgstore.CreateUserParams{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
		Role:      pgstore.RolePersonal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.queries.CreatePersonal(ctx, pgstore.CreatePersonalParams{
		ID:             userID,
		Description:    req.Description,
		VideoURL:       req.VideoURL,
		Experience:     req.Experience,
		Specialization: req.Specialization,
		Qualifications: req.Qualifications,
	})

	return &pgstore.UserResponse{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Role:      pgstore.RolePersonal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *AuthService) InitiatePasswordRecovery(ctx context.Context, email string) error {
	return fmt.Errorf("password recovery not implemented yet")
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	return fmt.Errorf("password reset not implemented yet")
}

func (s *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) generateSecureToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
