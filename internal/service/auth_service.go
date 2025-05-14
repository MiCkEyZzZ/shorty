package service

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"shorty/internal/models"
	"shorty/internal/repository"
	"shorty/pkg/logger"
)

var (
	ErrAuthCreation        = errors.New("failed to create user")
	ErrAuthNotFound        = errors.New("user not found")
	ErrAuthWrongCredential = errors.New("invalid email or password")
)

// AuthService provides methods for user authentication and registration.
type AuthService struct {
	Repo *repository.UserRepository
}

// NewAuthService создаёт новый экземпляр AuthService.
func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

// Registration registers a new user.
func (s *AuthService) Registration(ctx context.Context, name, email, password string, role models.Role, isBlocked bool) (*models.User, error) {
	// Check if user with the given email already exists
	exists, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Error occurred while searching for user", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrAuthNotFound, err)
	}
	if exists != nil {
		logger.Warn("User with this email is already registered", zap.String("email", email))
		return nil, fmt.Errorf("user with email %s is already registered", email)
	}

	// Hash the user's password
	hashedPassword, err := models.Hash(password)
	if err != nil {
		logger.Error("Error hashing password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create a new user model
	newUser := &models.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		IsBlocked: isBlocked,
	}

	// Save the new user in the database
	user, err := s.Repo.CreateUser(ctx, newUser)
	if err != nil {
		logger.Error("Error creating new user", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrAuthCreation, err)
	}
	logger.Info("New user successfully registered", zap.String("email", newUser.Email))
	return user, nil
}

// Login authenticates an existing user.
func (s *AuthService) Login(ctx context.Context, email, password string, role models.Role) (*models.User, error) {
	// Attempt to find user by email
	exists, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Invalid login credentials", zap.String("email", email))
			return nil, ErrAuthWrongCredential
		}
		logger.Error("Error retrieving user", zap.Error(err))
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Verify provided password against the stored hash
	err = models.VerifyPassword(exists.Password, password)
	if err != nil {
		logger.Error("Password verification failed", zap.Error(err))
		return nil, ErrAuthWrongCredential
	}
	logger.Info("User successfully logged in", zap.String("email", email))
	return exists, nil
}
