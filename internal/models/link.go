package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"
)

// Link represents the entity model for a shortened URL.
type Link struct {
	gorm.Model
	Url       string `json:"url"`
	Hash      string `json:"hash" gorm:"uniqueIndex"`
	Stats     []Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	IsBlocked bool   `json:"is_blocked" gorm:"default:false"`
}

// NewLink creates a new Link instance with a generated short hash.
func NewLink(url string) *Link {
	return &Link{
		Url:  url,
		Hash: generateHash(10),
	}
}

// generateHash generates a random base64-encoded string of the specified length.
// It is used as a short identifier for the link.
func generateHash(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Println("Error generating random data:", err)
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:n]
}
