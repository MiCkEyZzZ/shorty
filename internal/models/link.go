package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"
)

// Link стурктура представляет сущность ссылки.
type Link struct {
	gorm.Model
	Url       string `json:"url"`
	Hash      string `json:"hash" gorm:"uniqueIndex"`
	Stats     []Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	IsBlocked bool   `json:"is_blocked" gorm:"default:false"`
}

// NewLink создание нового экземпляра ссылки.
func NewLink(url string) *Link {
	return &Link{
		Url:  url,
		Hash: generateHash(10),
	}
}

// generateShortID ф-я для генерации короткого идентификатора.
func generateHash(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Println("Ошибка генерации случайных данных:", err)
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:n]
}
