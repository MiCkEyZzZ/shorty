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
	ErrAuthCreation        = errors.New("не удалось создать пользователя.")
	ErrAuthNotFound        = errors.New("не удалось найти пользователя.")
	ErrAuthWrongCredential = errors.New("Неверный адрес электронной почты или пароль")
)

// AuthService предоставляет методы для работы с авторизацией пользователей.
type AuthService struct {
	Repo *repository.UserRepository
}

// NewAuthService создаёт новый экземпляр AuthService.
func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

// Registration регистрация нового пользователя.
func (s *AuthService) Registration(ctx context.Context, name, email, password string) (string, error) {
	exists, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[AuthService] Ошибка при поиске пользователя: %v", err)
		return "", fmt.Errorf("%w: %v", ErrAuthNotFound, err)
	}
	if exists != nil {
		return "", fmt.Errorf("пользователь с email %s уже зарегистрирован", email)
	}

	hashedPassword, err := models.Hash(password)
	if err != nil {
		log.Printf("[AuthService] Ошибка хеширования пароля: %v", err)
		return "", fmt.Errorf("не удалось хешировать пароль: %w", err)
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	_, err = s.Repo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("[AuthService] Ошибка при создании нового пользователя: %v", err)
		return "", fmt.Errorf("%w: %v", ErrAuthCreation, err)
	}
	return user.Email, nil
}

// Login авторизация существующего пользователя.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	exists, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrAuthWrongCredential
		}
		log.Printf("[AuthService] Ошибка при поиске пользователя: %v", err)
		return "", fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	err = models.VerifyPassword(exists.Password, password)
	if err != nil {
		log.Printf("[AuthService] Ошибка при сравнении паролей: %v", err)
		return "", ErrAuthWrongCredential
	}
	return exists.Email, nil
}
