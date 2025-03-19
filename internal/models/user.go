package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User структура произвольного пользователя.
type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"index"`
	Password string
}

// Hash ф-я для хеширования пароля.
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword ф-я для проверки пароля пользователя.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
