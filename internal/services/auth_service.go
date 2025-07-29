package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries   *pgstore.Queries
	jwtSecret string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(queries *pgstore.Queries, jwtSecret string) *AuthService {
	return &AuthService{
		queries:   queries,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) AuthenticateWithPassword(ctx context.Context, email, password string) (*TokenPair, *pgstore.UserResponse, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found")
	}

	if !s.verifyPassword(password, user.Password) {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	tokens, err := s.generateTokenPair(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

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

	return tokens, userResponse, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	return nil, fmt.Errorf("refresh token functionality not implemented yet")
}

func (s *AuthService) RevokeToken(ctx context.Context, refreshToken string) error {
	return fmt.Errorf("token revocation not implemented yet")
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

func (s *AuthService) generateTokenPair(userID, role string) (*TokenPair, error) {
	accessClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
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
