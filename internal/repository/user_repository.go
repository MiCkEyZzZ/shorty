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

// UserRepository отвечает за операции с базой данных для сущности User.
type UserRepository struct {
	Database *db.DB
}

// NewUserRepository создаёт новый экземпляр UserRepository
func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{Database: db}
}

// CreateUser метод для создания нового пользователя.
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	res := r.Database.DB.WithContext(ctx).Create(user)
	if res.Error != nil {
		logger.Error("Ошибка создания пользователя", zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при сохранении пользователя в БД: %w", res.Error)
	}
	logger.Info("Пользователь успешно создан", zap.String("email", user.Email), zap.Uint("userID", user.ID))
	return user, nil
}

// GetUsers метод для получения списка пользователей.
func (r *UserRepository) GetUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	res := r.Database.DB.WithContext(ctx).Find(&users)
	if res.Error != nil {
		logger.Error("Ошибка при получении списка пользователей", zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при получении списка пользователей: %w", res.Error)
	}
	logger.Info("Получен список пользователей", zap.Int("usersCount", len(users)))
	return users, nil
}

// GetUserByID метод для поиска пользователя по идентификатору.
func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	res := r.Database.DB.WithContext(ctx).First(&user, userID)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Warn("Пользователь не найден по ID", zap.Uint("userID", userID))
			return nil, fmt.Errorf("пользователь с ID %d не найден", userID)
		}
		logger.Error("Ошибка при поиске пользователя по ID", zap.Uint("userID", userID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", res.Error)
	}
	logger.Info("Пользователь найден по ID", zap.Uint("userID", userID), zap.String("email", user.Email))
	return &user, nil
}

// GetUserByEmail метод для поиска пользователя по адресу электронной почты.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	res := r.Database.DB.
		WithContext(ctx).
		Where("email = ?", email).
		First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			logger.Warn("Пользователь не найден по email", zap.String("email", email))
			return nil, nil
		}
		logger.Error("Ошибка при поиске пользователя по email", zap.String("email", email), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", res.Error)
	}
	logger.Info("Пользователь найден по email", zap.String("email", email))
	return &user, nil
}

// UpdateUser метод для обновления данных пользователя.
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	res := r.Database.DB.
		WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(user)
	if res.Error != nil {
		logger.Error("Ошибка обновления пользователя", zap.Uint("userID", user.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при обновлении пользователя в БД: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		logger.Warn("Пользователь не найден для обновления", zap.Uint("userID", user.ID))
		return nil, fmt.Errorf("пользователь с ID %d не найден", user.ID)
	}
	logger.Info("Пользователь успешно обновлен", zap.Uint("userID", user.ID))
	return user, nil
}

// DeleteUser метод для удаления пользователя.
func (r *UserRepository) DeleteUser(ctx context.Context, userID uint) error {
	res := r.Database.DB.WithContext(ctx).Delete(&models.User{}, userID)
	if res.Error != nil {
		logger.Error("Ошибка удаления пользователя", zap.Uint("userID", userID), zap.Error(res.Error))
		return fmt.Errorf("ошибка при удалении пользователя из БД: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		logger.Warn("Пользователь не найден для удаления", zap.Uint("userID", userID))
		return fmt.Errorf("пользователь с ID %d не найден", userID)
	}
	logger.Info("Пользователь успешно удален", zap.Uint("userID", userID))
	return nil
}

// BlockUsers метод для блокировки пользователя.
func (r *UserRepository) BlockUsers(ctx context.Context, user *models.User) (*models.User, error) {
	res := r.Database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"is_blocked": user.IsBlocked,
		})
	if res.Error != nil {
		logger.Error("Ошибка при обновлении пользователя", zap.Uint("id", user.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при обновлении ссылки: %w", res.Error)
	}
	logger.Info("Пользователь обновлён", zap.Uint("id", user.ID), zap.Bool("is_blocked", user.IsBlocked))
	return user, nil
}

// UnBlockUsers метод для разблокировки пользователя.
func (r *UserRepository) UnBlockUsers(ctx context.Context, user *models.User) (*models.User, error) {
	res := r.Database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"is_blocked": user.IsBlocked,
		})
	if res.Error != nil {
		logger.Error("Ошибка снятие блокировки с пользователя", zap.Uint("id", user.ID), zap.Error(res.Error))
		return nil, fmt.Errorf("ошибка при снятии блокировки c пользователя: %w", res.Error)
	}

	logger.Info("Пользователь разблокирована", zap.Uint("id", user.ID), zap.Bool("is_blocked", user.IsBlocked))
	return user, nil
}
