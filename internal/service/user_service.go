package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"shorty/internal/models"
	"shorty/internal/repository"
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
		log.Printf("[UserService] Ошибка при получении списка пользователей: %v", err)
		return nil, fmt.Errorf("%w: %v", ErrUsersFound, err)
	}
	return users, nil
}

// GetByID ищет пользователя по ID
func (s *UserService) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("[UserService] Ошибка при поиске пользователя (ID: %d): %v", userID, err)
		return nil, fmt.Errorf("не удалось найти пользователя: %w", err)
	}
	return user, nil
}

// Update обновляет данные пользователя
func (s *UserService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	updatedUser, err := s.Repo.UpdateUser(ctx, user)
	if err != nil {
		log.Printf("[UserService] Ошибка при обновлении пользователя (ID: %d): %v", user.ID, err)
		return nil, fmt.Errorf("%w: %v", ErrUserUpdate, err)
	}
	return updatedUser, nil
}

// Delete удаляет пользователя по ID
func (s *UserService) Delete(ctx context.Context, userID uint) error {
	err := s.Repo.DeleteUser(ctx, userID)
	if err != nil {
		log.Printf("[UserService] Ошибка удаления пользователя (ID: %d): %v", userID, err)
		return fmt.Errorf("%w: %v", ErrUserDeletion, err)
	}
	log.Printf("[UserService] Пользователь (ID: %d) успешно удалён", userID)
	return nil
}
