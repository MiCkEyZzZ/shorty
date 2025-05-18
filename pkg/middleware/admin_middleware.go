package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"shorty/internal/common"
	"shorty/internal/models"
	"shorty/internal/service"
	"shorty/pkg/jwt"
	"shorty/pkg/logger"
	"shorty/pkg/res"
)

// AdminMiddleware проверяет, является ли пользователь администратором.
func AdminMiddleware(jwtService *jwt.JWT, userService service.UserServ) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Отсутствует заголовок Authorization")
				res.ERROR(w, common.ErrUnauthorized, http.StatusUnauthorized)
				return
			}

			// Ожидаем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Неверный формат токена")
				res.ERROR(w, common.ErrInvalidToken, http.StatusUnauthorized)
				return
			}

			// Разбираем JWT
			claims, err := jwtService.ParseToken(parts[1])
			if err != nil {
				logger.Warn("Ошибка парсинга JWT", zap.Error(err))
				res.ERROR(w, common.ErrInvalidToken, http.StatusUnauthorized)
				return
			}

			// Получаем пользователя из БД
			ctx := r.Context()
			user, err := userService.GetByID(ctx, claims.UserID)
			if err != nil {
				logger.Warn("Пользователь не найден", zap.Uint("userID", claims.UserID))
				res.ERROR(w, common.ErrUserNotFound, http.StatusUnauthorized)
				return
			}

			// Проверяем, является ли пользователь администратором
			if user.Role != models.RoleAdmin {
				logger.Warn("Доступ запрещён: недостаточно прав", zap.Uint("userID", user.ID))
				res.ERROR(w, common.ErrForbidden, http.StatusForbidden)
				return
			}

			// Передаём пользователя в контекст
			ctx = context.WithValue(ctx, common.UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
