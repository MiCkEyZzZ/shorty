package repository

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"shorty/internal/models"
	"shorty/pkg/db"
	"shorty/pkg/logger"
)

// UserRepository handles database operations for User entity.
type UserRepository struct {
	Database *db.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{Database: db}
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	res := r.Database.DB.WithContext(ctx).Create(user)
	if res.Error != nil {
		logger.Error("Failed to create user",
			zap.Error(res.Error),
			zap.String("email", user.Email))
		return nil, fmt.Errorf("failed to create user in database: %w", res.Error)
	}
	logger.Info("User created successfully",
		zap.String("email", user.Email),
		zap.Uint("userID", user.ID))
	return user, nil
}

// GetUsers retrieves a list of users with pagination.
func (r *UserRepository) GetUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User

	query := r.Database.DB.WithContext(ctx).Model(&models.User{}).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	res := query.Find(&users)
	if res.Error != nil {
		logger.Error("Failed to get users list",
			zap.Error(res.Error),
			zap.Int("limit", limit),
			zap.Int("offset", offset))
		return nil, fmt.Errorf("failed to get users from database: %w", res.Error)
	}

	logger.Info("Users list retrieved successfully", zap.Int("count", len(users)))
	return users, nil
}

// GetUserByID finds a user by their ID.
func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	res := r.Database.DB.WithContext(ctx).First(&user, userID)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Warn("User not found by ID", zap.Uint("userID", userID))
			return nil, gorm.ErrRecordNotFound // Return original GORM error for service layer
		}
		logger.Error("Failed to find user by ID",
			zap.Uint("userID", userID),
			zap.Error(res.Error))
		return nil, fmt.Errorf("failed to find user: %w", res.Error)
	}

	logger.Debug("User found by ID",
		zap.Uint("userID", userID),
		zap.String("email", user.Email))
	return &user, nil
}

// GetUserByEmail finds a user by their email address.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var user models.User
	res := r.Database.DB.WithContext(ctx).
		Where("email = ?", email).
		First(&user)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Warn("User not found by email", zap.String("email", email))
			return nil, gorm.ErrRecordNotFound // Return original GORM error
		}
		logger.Error("Failed to find user by email",
			zap.String("email", email),
			zap.Error(res.Error))
		return nil, fmt.Errorf("failed to find user by email: %w", res.Error)
	}

	logger.Debug("User found by email", zap.String("email", email))
	return &user, nil
}

// UpdateUser updates user data in the database.
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil || user.ID == 0 {
		return nil, errors.New("invalid user data")
	}

	res := r.Database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(user)

	if res.Error != nil {
		logger.Error("Failed to update user",
			zap.Uint("userID", user.ID),
			zap.Error(res.Error))
		return nil, fmt.Errorf("failed to update user in database: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		logger.Warn("No user found to update", zap.Uint("userID", user.ID))
		return nil, gorm.ErrRecordNotFound
	}

	logger.Info("User updated successfully", zap.Uint("userID", user.ID))
	return user, nil
}

// DeleteUser removes a user from the database.
func (r *UserRepository) DeleteUser(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	res := r.Database.DB.WithContext(ctx).Delete(&models.User{}, userID)
	if res.Error != nil {
		logger.Error("Failed to delete user",
			zap.Uint("userID", userID),
			zap.Error(res.Error))
		return fmt.Errorf("failed to delete user from database: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		logger.Warn("No user found to delete", zap.Uint("userID", userID))
		return gorm.ErrRecordNotFound
	}

	logger.Info("User deleted successfully", zap.Uint("userID", userID))
	return nil
}

// BlockUsers blocks a user.
func (r *UserRepository) BlockUsers(ctx context.Context, user *models.User) (*models.User, error) {
	return r.updateUserBlockStatus(ctx, user, true)
}

// UnBlockUsers unblocks a user.
func (r *UserRepository) UnBlockUsers(ctx context.Context, user *models.User) (*models.User, error) {
	return r.updateUserBlockStatus(ctx, user, false)
}

// GetBlockedUsersCount returns the count of blocked users.
func (r *UserRepository) GetBlockedUsersCount(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("is_blocked = ?", true).
		Count(&count)

	if result.Error != nil {
		logger.Error("Failed to get blocked users count", zap.Error(result.Error))
		return 0, fmt.Errorf("failed to count blocked users: %w", result.Error)
	}

	logger.Debug("Blocked users count retrieved", zap.Int64("count", count))
	return count, nil
}

// CountUsers returns the total number of users in the database.
func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	result := r.Database.DB.WithContext(ctx).
		Model(&models.User{}).
		Count(&count)

	if result.Error != nil {
		logger.Error("Failed to count users", zap.Error(result.Error))
		return 0, fmt.Errorf("failed to count users: %w", result.Error)
	}

	logger.Debug("Users count retrieved", zap.Int64("count", count))
	return count, nil
}

// GetUsersByStatus retrieves users by their block status with pagination.
func (r *UserRepository) GetUsersByStatus(ctx context.Context, isBlocked bool, limit, offset int) ([]*models.User, error) {
	var users []*models.User

	query := r.Database.DB.WithContext(ctx).Where("is_blocked = ?", isBlocked).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	res := query.Find(&users)
	if res.Error != nil {
		status := "active"
		if isBlocked {
			status = "blocked"
		}
		logger.Error(fmt.Sprintf("Failed to get %s users", status),
			zap.Error(res.Error),
			zap.Bool("is_blocked", isBlocked))
		return nil, fmt.Errorf("failed to get users by status: %w", res.Error)
	}

	return users, nil
}

// UserExists checks if a user exists by ID.
func (r *UserRepository) UserExists(ctx context.Context, userID uint) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user ID")
	}

	var count int64
	res := r.Database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Count(&count)

	if res.Error != nil {
		logger.Error("Failed to check user existence",
			zap.Uint("userID", userID),
			zap.Error(res.Error))
		return false, fmt.Errorf("failed to check user existence: %w", res.Error)
	}

	return count > 0, nil
}

// EmailExists checks if an email already exists in the database.
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email cannot be empty")
	}

	var count int64
	res := r.Database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count)

	if res.Error != nil {
		logger.Error("Failed to check email existence",
			zap.String("email", email),
			zap.Error(res.Error))
		return false, fmt.Errorf("failed to check email existence: %w", res.Error)
	}

	return count > 0, nil
}

// updateUserBlockStatus is a helper method to update user block status.
func (r *UserRepository) updateUserBlockStatus(ctx context.Context, user *models.User, isBlocked bool) (*models.User, error) {
	if user == nil || user.ID == 0 {
		return nil, errors.New("invalid user data")
	}

	res := r.Database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Update("is_blocked", isBlocked)

	if res.Error != nil {
		action := "block"
		if !isBlocked {
			action = "unblock"
		}
		logger.Error(fmt.Sprintf("Failed to %s user", action),
			zap.Uint("userID", user.ID),
			zap.Error(res.Error))
		return nil, fmt.Errorf("failed to update user block status: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		logger.Warn("No user found to update block status", zap.Uint("userID", user.ID))
		return nil, gorm.ErrRecordNotFound
	}

	user.IsBlocked = isBlocked

	action := "blocked"
	if !isBlocked {
		action = "unblocked"
	}
	logger.Info(fmt.Sprintf("User %s successfully", action),
		zap.Uint("userID", user.ID),
		zap.Bool("is_blocked", isBlocked))

	return user, nil
}
