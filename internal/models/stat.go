package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Stat represents the entity model for click statistics associated with a shortened link.
type Stat struct {
	gorm.Model
	LinkID    uint           `json:"link_id" gorm:"index"`
	Clicks    int            `json:"clicks"`
	Date      datatypes.Date `json:"date" gorm:"index"`
	IP        string         `json:"ip"`
	Referrer  string         `json:"referrer"`
	UserAgent string         `json:"user_agent"`
}
