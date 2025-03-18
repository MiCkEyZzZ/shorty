package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"shorty/pkg/db"
)

// UserRepository предоставляет методы для работы с пользователями в БД
type UserRepository struct {
	Database *db.DB
}

// NewUserRepository создаёт новый экземпляр UserRepository
func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{Database: db}
}

// Create добавляет нового пользователя в базу данных
func (r *UserRepository) Create(ctx context.Context, user *User) (*User, error) {
	res := r.Database.WithContext(ctx).Create(user)
	if res.Error != nil {
		log.Printf("[UserRepository] Ошибка создания пользователя: %v", res.Error)
		return nil, fmt.Errorf("ошибка при сохранении пользователя в БД: %w", res.Error)
	}
	return user, nil
}

// GetAll возвращает список всех пользователей
func (r *UserRepository) FindAll(ctx context.Context) ([]*User, error) {
	var users []*User
	res := r.Database.WithContext(ctx).Find(&users)
	if res.Error != nil {
		log.Printf("[UserRepository] Ошибка при получении списка пользователей: %v", res.Error)
		return nil, fmt.Errorf("ошибка при получении списка пользователей: %w", res.Error)
	}
	return users, nil
}

// FindByID ищет пользователя по его ID
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	var user User
	res := r.Database.WithContext(ctx).Find(&user, id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Printf("[UserRepository] Пользователь с ID %d не найден", id)
			return nil, fmt.Errorf("пользователь с ID %d не найден", id)
		}
		log.Printf("[UserRepository] Ошибка при поиске пользователя (ID: %d): %v", id, res.Error)
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", res.Error)
	}
	return &user, nil
}

// FindByEmail ищет пользователя по email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	res := r.Database.WithContext(ctx).Where("email = ?", email).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Printf("[UserRepository] Пользователь с email %s не найден", email)
			return nil, fmt.Errorf("пользователь с email %s не найден", email)
		}
		log.Printf("[UserRepository] Ошибка при поиске пользователя (email: %s): %v", email, res.Error)
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", res.Error)
	}
	return &user, nil
}

// Update обновляет данные пользователя в базе
func (r *UserRepository) Update(ctx context.Context, user *User) (*User, error) {
	res := r.Database.DB.WithContext(ctx).Clauses(clause.Returning{}).Updates(user)
	if res.Error != nil {
		log.Printf("[UserRepository] Ошибка обновления пользователя (ID: %d): %v", user.ID, res.Error)
		return nil, fmt.Errorf("ошибка при обновлении пользователя в БД: %w", res.Error)
	}
	return user, nil
}

// Delete удаляет пользователя по ID
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	res := r.Database.WithContext(ctx).Delete(&User{}, id)
	if res.Error != nil {
		log.Printf("[UserRepository] Ошибка удаления пользователя (ID: %d): %v", id, res.Error)
		return fmt.Errorf("ошибка при удалении пользователя из БД: %w", res.Error)
	}
	return nil
}
