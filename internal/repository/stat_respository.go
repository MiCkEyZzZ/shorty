package repository

import "shorty/pkg/db"

type StatRepository struct {
	Database *db.DB
}

func NewStatRepository(db *db.DB) *StatRepository {
	return &StatRepository{Database: db}
}
