package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
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

func NewAuthService(queries *pgstore.Queries, jwtSecret string) *AuthService {
	return &AuthService{
		queries:   queries,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, req pgstore.AuthenticateRequest) (*pgstore.AuthenticateResponse, error) {
	dbUser, err := s.queries.GetUserForAuth(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if dbUser == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if dbUser.Password == "" {
		return nil, fmt.Errorf("user does not have a password, use social login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Get full user details
	userService := NewUserService(s.queries)
	user, err := userService.GetUserByID(ctx, dbUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	// Generate JWT token
	token, err := s.GenerateJWT(dbUser.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.GenerateRefreshToken(ctx, dbUser.ID, req.DeviceInfo, req.IPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &pgstore.AuthenticateResponse{
		User:         *user,
		Token:        token,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    int64((time.Hour * 24 * 7).Seconds()), // 7 days in seconds
	}, nil
}

func (s *AuthService) CreateUser(ctx context.Context, req *pgstore.CreateUserParams) (*pgstore.UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New()
	now := time.Now()

	// Create user
	_, err = s.queries.CreateUser(ctx, pgstore.CreateUserParams{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPassword),
		Role:      pgstore.Role(req.Role),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get the created user
	userService := NewUserService(s.queries)
	return userService.GetUserByID(ctx, userID)
}

func (s *AuthService) CreateStudent(ctx context.Context, req *pgstore.CreateStudentWithUserRequest) (*pgstore.UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New()
	now := time.Now()

	// Create base user first
	_, err = s.queries.CreateUser(ctx, pgstore.CreateUserParams{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPassword),
		Role:      pgstore.RoleStudent,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create student profile
	err = s.queries.CreateStudent(ctx, pgstore.CreateStudentParams{
		ID:                    userID,
		BornDate:              req.BornDate,
		Age:                   int32(req.Age),
		Weight:                req.Weight,
		Objective:             req.Objective,
		TrainingFrequency:     req.TrainingFrequency,
		DidBodybuilding:       req.DidBodybuilding,
		MedicalCondition:      req.MedicalCondition,
		PhysicalActivityLevel: req.PhysicalActivityLevel,
		Observations:          req.Observations,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create student profile: %w", err)
	}

	// Get the complete user with student profile
	userService := NewUserService(s.queries)
	return userService.GetUserByID(ctx, userID)
}

func (s *AuthService) CreatePersonal(ctx context.Context, req *pgstore.CreatePersonalWithUserRequest) (*pgstore.UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New()
	now := time.Now()

	// Create base user first
	_, err = s.queries.CreateUser(ctx, pgstore.CreateUserParams{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPassword),
		Role:      pgstore.RolePersonal,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create personal trainer profile
	err = s.queries.CreatePersonal(ctx, pgstore.CreatePersonalParams{
		ID:             userID,
		Description:    req.Description,
		VideoURL:       req.PresentationVideo,
		Experience:     req.Experience,
		Specialization: req.Specialization,
		Qualifications: req.Qualifications,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create personal trainer profile: %w", err)
	}

	// Get the complete user with personal trainer profile
	userService := NewUserService(s.queries)
	return userService.GetUserByID(ctx, userID)
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Generate reset token
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	// Store reset token
	err = s.queries.CreatePasswordResetToken(ctx, pgstore.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// TODO: Send email with reset token
	// For now, just log it (in production, send via email service)
	fmt.Printf("Password reset token for %s: %s\n", email, token)

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get reset token
	resetToken, err := s.queries.GetPasswordResetToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}
	if resetToken == nil {
		return fmt.Errorf("invalid or expired token")
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return fmt.Errorf("token has expired")
	}

	// Check if token was already used
	if resetToken.UsedAt != nil {
		return fmt.Errorf("token has already been used")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	err = s.queries.UpdateUserPassword(ctx, pgstore.UpdateUserPasswordParams{
		ID:       resetToken.UserID,
		Password: string(hashedPassword),
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	err = s.queries.MarkPasswordResetTokenAsUsed(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

func (s *AuthService) GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
}

// GenerateRefreshToken creates a new refresh token for the user
func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID, deviceInfo, ipAddress *string) (*pgstore.RefreshToken, error) {
	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random token: %w", err)
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Set expiration to 30 days
	expiresAt := time.Now().Add(time.Hour * 24 * 30)

	// Create refresh token in database
	refreshToken, err := s.queries.CreateRefreshToken(ctx, pgstore.CreateRefreshTokenParams{
		UserID:     userID,
		Token:      token,
		ExpiresAt:  expiresAt,
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return refreshToken, nil
}

// RefreshAccessToken validates a refresh token and generates a new access token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (*pgstore.RefreshTokenResponse, error) {
	// Get refresh token from database
	refreshToken, err := s.queries.GetRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	if refreshToken == nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Update last used timestamp
	err = s.queries.UpdateRefreshTokenLastUsed(ctx, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	// Generate new access token
	accessToken, err := s.GenerateJWT(refreshToken.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token (rotate refresh tokens for security)
	newRefreshToken, err := s.GenerateRefreshToken(ctx, refreshToken.UserID, refreshToken.DeviceInfo, refreshToken.IPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// Revoke old refresh token
	err = s.queries.RevokeRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	return &pgstore.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
		ExpiresIn:    int64((time.Hour * 24 * 7).Seconds()), // 7 days in seconds
	}, nil
}

// RevokeRefreshToken revokes a specific refresh token
func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	return s.queries.RevokeRefreshToken(ctx, refreshToken)
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user (useful for logout from all devices)
func (s *AuthService) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return s.queries.RevokeAllUserRefreshTokens(ctx, userID)
}

// GetUserRefreshTokens returns all active refresh tokens for a user
func (s *AuthService) GetUserRefreshTokens(ctx context.Context, userID uuid.UUID) ([]pgstore.GetRefreshTokensByUserIDRow, error) {
	return s.queries.GetRefreshTokensByUserID(ctx, userID)
}

// CleanupExpiredTokens removes expired refresh tokens (should be run periodically)
func (s *AuthService) CleanupExpiredTokens(ctx context.Context) error {
	return s.queries.CleanupExpiredRefreshTokens(ctx)
}

// ExtractIPAddress extracts IP address from request
func ExtractIPAddress(remoteAddr string) *string {
	if remoteAddr == "" {
		return nil
	}
	
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// If splitting fails, assume it's just an IP
		host = remoteAddr
	}
	
	return &host
}
