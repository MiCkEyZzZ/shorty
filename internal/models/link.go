package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
}

func NewLink(url string) *Link {
	return &Link{
		Url:  url,
		Hash: generateHash(6),
	}
}

// generateShortID ф-я для генерации короткого идентификатора.
func generateHash(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Println("Ошибка генерации случайных данных:", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
