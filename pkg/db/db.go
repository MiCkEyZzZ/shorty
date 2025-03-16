package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shorty/internal/config"
)

type DB struct {
	*gorm.DB
}

func NewDatabase(cfg *config.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Db.Dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %w", err)
	}
	fmt.Println("Подключение к БД успешно!")

	return &DB{db}, nil
}
