package user

import (
	"gorm.io/gorm"
)

// User структура произвольного пользователя.
type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"index"`
	Password string
}
