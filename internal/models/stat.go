package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Stat struct {
	gorm.Model
	LinkID uint           `json:"link_id" gorm:"index"`
	Clicks int            `json:"clicks"`
	Date   datatypes.Date `json:"date" gorm:"index"`
}
