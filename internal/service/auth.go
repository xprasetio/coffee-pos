package service

import (
	"context"
	"errors"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtExpireHours int
}

// NewAuthService creates a new AuthService instance
func NewAuthService(
	userRepo repository.UserRepository,
	jwtSecret string,
	jwtExpireHours int,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

// RegisterRequestClaims represents JWT claims for authentication
type RegisterRequestClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req entity.RegisterRequest) (*entity.UserResponse, error) {
	// Check if email already exists - do this BEFORE hashing to avoid wasting time
	// bcrypt takes 150-300ms, so we should not hash if email is already registered
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Hash password with cost 12
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user entity
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		IsActive: true,
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

// Login authenticates a user and returns JWT token
func (s *AuthService) Login(ctx context.Context, req entity.LoginRequest) (*entity.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// Don't reveal whether email exists or not
		return nil, errors.New("email atau password salah")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("akun dinonaktifkan")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	// Generate JWT token
	expiredAt := time.Now().Add(time.Duration(s.jwtExpireHours) * time.Hour)
	claims := &RegisterRequestClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &entity.LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiredAt,
		User:      user.ToResponse(),
	}, nil
}
