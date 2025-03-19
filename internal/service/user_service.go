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
	ErrUserNotFound = errors.New("пользователь не найден")
	ErrUsersFound   = errors.New("не удалось получить список пользователей")
	ErrUserUpdate   = errors.New("не удалось обновить пользователя")
	ErrUserDeletion = errors.New("не удалось удалить пользователя")
)

type UserService struct {
	Repo *repository.UserRepository
}

// NewUserService создаёт новый экземпляр UserService
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// GetAll возвращает список всех пользователей
func (s *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	users, err := s.Repo.GetUsers(ctx)
	if err != nil {
		logger.Error("Ошибка при получении списка пользователей", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrUsersFound, err)
	}
	logger.Info("Список пользователей получен", zap.Int("count", len(users)))
	return users, nil
}

// GetByID ищет пользователя по ID
func (s *UserService) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Пользователь не найден", zap.Uint("userID", userID))
			return nil, ErrUserNotFound
		}
		logger.Error("Ошибка при поиске пользователя", zap.Uint("userID", userID), zap.Error(err))
		return nil, fmt.Errorf("не удалось найти пользователя: %w", err)
	}
	logger.Info("Пользователь найден", zap.Uint("userID", user.ID), zap.String("email", user.Email))
	return user, nil
}

// Update обновляет данные пользователя
func (s *UserService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	updatedUser, err := s.Repo.UpdateUser(ctx, user)
	if err != nil {
		logger.Error("Ошибка при обновлении пользователя", zap.Uint("userID", user.ID), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrUserUpdate, err)
	}
	logger.Info("Данные пользователя обновлены", zap.Uint("userID", updatedUser.ID), zap.String("email", updatedUser.Email))
	return updatedUser, nil
}

// Delete удаляет пользователя по ID
func (s *UserService) Delete(ctx context.Context, userID uint) error {
	err := s.Repo.DeleteUser(ctx, userID)
	if err != nil {
		logger.Error("Ошибка удаления пользователя", zap.Uint("userID", userID), zap.Error(err))
		return fmt.Errorf("%w: %v", ErrUserDeletion, err)
	}
	logger.Info("Пользователь успешно удалён", zap.Uint("userID", userID))
	return nil
}
