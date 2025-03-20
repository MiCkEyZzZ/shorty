package middleware

import (
	"context"
	"errors"
	"net/http"

	"shorty/internal/models"
)

type contextKey string

const UserContextKey contextKey = "user"

// AdminOnly — middleware для проверки роли администратора.
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromContext(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if user.Role != models.RoleAdmin {
			http.Error(w, "Forbidden: Access denied", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// WithUserContext — middleware для добавления пользователя в контекст запроса.
func WithUserContext(authService AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, err := authService.GetUserFromToken(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next(w, r.WithContext(ctx))
		}
	}
}

// getUserFromContext — функция для извлечения пользователя из контекста.
func getUserFromContext(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok || user == nil {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

// AuthService — интерфейс для работы с авторизацией.
type AuthService interface {
	GetUserFromToken(r *http.Request) (*models.User, error)
}
