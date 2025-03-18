package user

import (
	"context"
	"fmt"
	"log"
)

// Ошибки, связанные с работой с пользователями.
var (
	ErrUserNotFound = fmt.Errorf("пользователь не найдена")
	ErrUsersFound   = fmt.Errorf("не удалось получить список пользователей")
	ErrUserUpdate   = fmt.Errorf("не удалось обновить пользователя")
	ErrUserDeletion = fmt.Errorf("не удалось удалить пользователя")
)

// UserService предоставляет методы для работы с пользователями.
type UserService struct {
	Repo *UserRepository
}

// NewUserService создаёт новый экземпляр UserService.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// FindAll получает список всех пользователей.
func (s *UserService) FindAll(ctx context.Context) ([]*User, error) {
	users, err := s.Repo.FindAll(ctx)
	if err != nil {
		log.Printf("[UserService] Ошибка при получении списка пользователей: %v", err)
		return nil, fmt.Errorf("%w: %v", ErrUsersFound, err)
	}
	return users, nil
}

// FindByID ищет пользователя по его ID.
func (s *UserService) FindByID(ctx context.Context, id uint) (*User, error) {
	user, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[UserService] Ошибка при поиске пользователя (ID: %d): %v", id, err)
		if err.Error() == fmt.Sprintf("пользователь с ID %d не найдена", id) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("не удалось найти пользователя: %w", err)
	}
	return user, nil
}

// Update обновляет данные пользователя.
func (s *UserService) Update(ctx context.Context, user *User) (*User, error) {
	updateUser, err := s.Repo.Update(ctx, user)
	if err != nil {
		log.Printf("[UserService] Ошибка при обновлении пользователя (ID: %d): %v", user.ID, err)
		return nil, fmt.Errorf("%w: %v", ErrUserUpdate, err)
	}
	return updateUser, nil
}

// Delete удаляет пользователя по ID.
func (s *UserService) Delete(ctx context.Context, id uint) error {
	err := s.Repo.Delete(ctx, id)
	if err != nil {
		log.Printf("[UserService] Ошибка удаления пользователя (ID: %d): %v", id, err)
		return fmt.Errorf("%w: %v", ErrUserDeletion, err)
	}
	return nil
}
