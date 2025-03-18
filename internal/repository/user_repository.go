package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"shorty/internal/models"
	"shorty/pkg/db"
)

// UserRepository отвечает за операции с базой данных для сущности User.
type UserRepository struct {
	Database *db.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{Database: db}
}

// CreateUser добавляет нового пользователя в базу данных
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	res := r.Database.DB.WithContext(ctx).Create(user)
	if res.Error != nil {
		log.Printf("[UserRepository] Ошибка создания пользователя: %v", res.Error)
		return nil, fmt.Errorf("ошибка при сохранении пользователя в БД: %w", res.Error)
	}
	return user, nil
}

// GetUsers получает всех пользователей из базы
func (r *UserRepository) GetUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	res := r.Database.DB.WithContext(ctx).Find(&users)
	if res.Error != nil {
		log.Printf("[UserRepository] Ошибка при получении списка пользователей: %v", res.Error)
		return nil, fmt.Errorf("ошибка при получении списка пользователей: %w", res.Error)
	}
	return users, nil
}

// GetUserByID ищет пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	res := r.Database.DB.WithContext(ctx).First(&user, userID)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Printf("[UserRepository] Пользователь с ID %d не найден", userID)
			return nil, fmt.Errorf("пользователь с ID %d не найден", userID)
		}
		log.Printf("[UserRepository] Ошибка при поиске пользователя (ID: %d): %v", userID, res.Error)
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", res.Error)
	}
	return &user, nil
}

// GetUserByEmail ищет пользователя по email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	res := r.Database.DB.
		WithContext(ctx).
		Where("email = ?", email).
		First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Printf("[UserRepository] Пользователь с email %s не найден", email)
			return nil, nil
		}
		log.Printf("[UserRepository] Ошибка при поиске пользователя (email: %s): %v", email, res.Error)
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", res.Error)
	}
	return &user, nil
}

// UpdateUser обновляет данные пользователя
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	res := r.Database.DB.
		WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(user)
	if res.Error != nil {
		log.Printf("[UserRepository] Ошибка обновления пользователя (ID: %d): %v", user.ID, res.Error)
		return nil, fmt.Errorf("ошибка при обновлении пользователя в БД: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("[UserRepository] Пользователь с ID %d не найден для обновления", user.ID)
		return nil, fmt.Errorf("пользователь с ID %d не найден", user.ID)
	}
	return user, nil
}

// DeleteUser удаляет пользователя из базы
func (r *UserRepository) DeleteUser(ctx context.Context, userID uint) error {
	res := r.Database.DB.WithContext(ctx).Delete(&models.User{}, userID)
	if res.Error != nil {
		log.Printf("[UserRepository] Ошибка удаления пользователя (ID: %d): %v", userID, res.Error)
		return fmt.Errorf("ошибка при удалении пользователя из БД: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		log.Printf("[UserRepository] Пользователь с ID %d не найден для удаления", userID)
		return fmt.Errorf("пользователь с ID %d не найден", userID)
	}
	return nil
}
