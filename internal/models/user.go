package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role defines the role of a user in the system.
type Role string

const (
	RoleUser  Role = "user"  // Standard user role
	RoleAdmin Role = "admin" // Administrator role
)

// User represents the user entity.
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"index" json:"email"`
	Password  string         `json:"password,omitempty"`
	Role      Role           `json:"role"`
	IsBlocked bool           `json:"is_blocked" gorm:"default:false"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Hash hashes the given password using bcrypt.
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword compares a bcrypt-hashed password with its possible plaintext equivalent.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
