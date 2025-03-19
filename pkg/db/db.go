package db

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shorty/internal/config"
	"shorty/pkg/logger"
)

type DB struct {
	*gorm.DB
}

func NewDatabase(cfg *config.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Db.Dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Ошибка подключения к базе данных", zap.Error(err))
		return nil, fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	// Проверка подключения к БД
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Ошибка при получении объекта базы данных для проверки состояния", zap.Error(err))
		return nil, fmt.Errorf("не удалось получить доступ к объекту базы данных: %w", err)
	}

	// Проверка состояния соединения с базой данных
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Не удалось подключиться к базе данных", zap.Error(err))
		return nil, fmt.Errorf("не удалось установить соединение с БД: %w", err)
	}
	fmt.Println("Подключение к БД успешно!")
	logger.Info("Подключение к БД успешно!")

	return &DB{db}, nil
}
