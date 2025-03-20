package jwt

import (
	"errors"
	"shorty/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTData содержит информацию, которая будет в токене
type JWTData struct {
	UserID    uint
	Email     string
	Role      models.Role
	IsBlocked bool
}

// JWT - структура для работы с токенами
type JWT struct {
	Secret string
}

// NewJWT создаёт новый экземпляр JWT
func NewJWT(secret string) *JWT {
	return &JWT{Secret: secret}
}

// CreateToken создаёт новый JWT с данными пользователя
func (j *JWT) CreateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"is_blocked": user.IsBlocked,
		"exp":        time.Now().Add(24 * time.Hour).Unix(), // Токен живёт 24 часа
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(j.Secret))
}

// ParseToken парсит JWT и возвращает данные
func (j *JWT) ParseToken(token string) (*JWTData, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64) // JWT использует float64 для чисел
	if !ok {
		return nil, errors.New("invalid user_id")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role")
	}

	isBlocked, ok := claims["is_blocked"].(bool)
	if !ok {
		return nil, errors.New("invalid is_blocked")
	}

	return &JWTData{
		UserID:    uint(userID),
		Email:     email,
		Role:      models.Role(role), // Приводим строку обратно в `models.Role`
		IsBlocked: isBlocked,
	}, nil
}
