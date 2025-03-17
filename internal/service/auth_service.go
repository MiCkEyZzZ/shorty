package service

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"shorty/internal/models"
	"shorty/internal/repository"
)

// Ошибки, связанные с работой авторизации.
var (
	ErrAuthCreation        = fmt.Errorf("не удалось создать пользователя.")
	ErrAuthNotFound        = fmt.Errorf("не удалось найти пользователя.")
	ErrAuthWrongCredential = fmt.Errorf("Неверный адрес электронной почты или пароль")
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
	exists, err := s.Repo.FindByEmail(ctx, email)
	if exists != nil {
		log.Printf("[AuthService] Ошибка при поиске пользователя: %v", err)
		return "", fmt.Errorf("%w: %v", ErrAuthWrongCredential, err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	_, err = s.Repo.Create(ctx, user)
	if err != nil {
		log.Printf("[AuthService] Ошибка при создании нового пользователя: %v", err)
		return "", fmt.Errorf("%w: %v", ErrAuthCreation, err)
	}
	return user.Email, nil
}

// Login авторизация существующего пользователя.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	exists, err := s.Repo.FindByEmail(ctx, email)
	if exists == nil {
		log.Printf("[AuthService] Ошибка при поиске пользователя: %v", err)
		return "", fmt.Errorf("%w: %v", ErrAuthWrongCredential, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(exists.Password), []byte(password))
	if err != nil {
		log.Printf("[AuthService] Ошибка при сравнении паролей: %v", err)
		return "", fmt.Errorf("%w: %v", ErrAuthWrongCredential, err)
	}
	return exists.Email, nil
}
