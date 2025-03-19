package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Stat стурктура представляет сущность статистики.
type Stat struct {
	gorm.Model
	LinkID uint           `json:"link_id" gorm:"index"`
	Clicks int            `json:"clicks"`
	Date   datatypes.Date `json:"date" gorm:"index"`
}
