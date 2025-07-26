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

// User service errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsersFetchFailed   = errors.New("failed to fetch users")
	ErrUserUpdateFailed   = errors.New("failed to update user")
	ErrUserDeleteFailed   = errors.New("failed to delete user")
	ErrUserBlockFailed    = errors.New("failed to block user")
	ErrUserUnblockFailed  = errors.New("failed to unblock user")
	ErrUserAlreadyBlocked = errors.New("user is already blocked")
	ErrUserNotBlocked     = errors.New("user is not blocked")
	ErrInvalidUserData    = errors.New("invalid user data")
)

// UserService provides methods for working with users.
type UserService struct {
	Repo repository.UserRepo
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

// GetAll retrieves a list of users with pagination.
func (s *UserService) GetAll(ctx context.Context, limit, offset int) ([]*models.User, error) {
	if limit <= 0 || limit >= 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.Repo.GetUsers(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to get users list",
			zap.Error(err),
			zap.Int("limit", limit),
			zap.Int("offset", offset))
		return nil, fmt.Errorf("%w: %v", ErrUsersFetchFailed, err)
	}
	logger.Info("Users list retrieved successfully", zap.Int("count", len(users)))
	return users, nil
}

// GetByID retrieves a user by their ID.
func (s *UserService) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, ErrInvalidUserData
	}

	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("User not found", zap.Uint("userID", userID))
			return nil, ErrUserNotFound
		}
		logger.Error("Failed to find user",
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	logger.Info("User found successfully",
		zap.Uint("userID", user.ID),
		zap.String("email", user.Email))
	return user, nil
}

// Update updates user information.
func (s *UserService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil || user.ID == 0 {
		return nil, ErrInvalidUserData
	}

	existingUser, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if user.Email != "" {
		existingUser.Email = user.Email
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}

	updatedUser, err := s.Repo.UpdateUser(ctx, user)
	if err != nil {
		logger.Error("Failed to update user",
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrUserUpdateFailed, err)
	}

	logger.Info("User updated successfully",
		zap.Uint("userID", updatedUser.ID),
		zap.String("email", updatedUser.Email))
	return updatedUser, nil
}

// Delete removes a user by their ID.
func (s *UserService) Delete(ctx context.Context, userID uint) error {
	if userID == 0 {
		return ErrInvalidUserData
	}

	_, err := s.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	err = s.Repo.DeleteUser(ctx, userID)
	if err != nil {
		logger.Error("Failed to delete user",
			zap.Uint("userID", userID),
			zap.Error(err))
		return fmt.Errorf("%w: %v", ErrUserDeleteFailed, err)
	}

	logger.Info("User deleted successfully", zap.Uint("userID", userID))
	return nil
}

// Block blocks a user by their ID.
func (s *UserService) Block(ctx context.Context, userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, ErrInvalidUserData
	}

	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		logger.Error("Failed to find user for blocking",
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, err
	}

	user.IsBlocked = true
	updatedUser, err := s.Repo.BlockUsers(ctx, user)
	if err != nil {
		logger.Error("Failed to block user",
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrUserBlockFailed, err)
	}

	logger.Info("User blocked successfully", zap.Uint("userID", updatedUser.ID))
	return updatedUser, nil
}

// UnBlock unblocks a user by their ID.
func (s *UserService) UnBlock(ctx context.Context, userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, ErrInvalidUserData
	}

	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		logger.Error("Failed to find user for unblocking",
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, err
	}

	if !user.IsBlocked {
		logger.Warn("User is not blocked", zap.Uint("userID", userID))
		return nil, ErrUserNotBlocked
	}

	user.IsBlocked = false
	updatedUser, err := s.Repo.UnBlockUsers(ctx, user)
	if err != nil {
		logger.Error("Failed to unblock user",
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrUserUnblockFailed, err)
	}

	logger.Info("User unblocked successfully", zap.Uint("userID", updatedUser.ID))
	return updatedUser, nil
}

// GetBlockedUsersCount returns the count of blocked users.
func (s *UserService) GetBlockedUsersCount(ctx context.Context) (int64, error) {
	count, err := s.Repo.GetBlockedUsersCount(ctx)
	if err != nil {
		logger.Error("Failed to get blocked users count", zap.Error(err))
		return 0, fmt.Errorf("failed to get blocked users count: %w", err)
	}

	logger.Info("Blocked users count retrieved", zap.Int64("count", count))
	return count, nil
}

// IsBlocked checks if a user is blocked.
func (s *UserService) IsBlocked(ctx context.Context, userID uint) (bool, error) {
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.IsBlocked, nil
}

// Count returns the total number of users.
func (s *UserService) Count(ctx context.Context) (int64, error) {
	count, err := s.Repo.CountUsers(ctx)
	if err != nil {
		logger.Error("Failed to count users", zap.Error(err))
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	logger.Info("Users count retrieved", zap.Int64("count", count))
	return count, nil
}
