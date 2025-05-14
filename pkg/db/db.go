package db

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shorty/internal/config"
	"shorty/pkg/logger"
)

// DB is a wrapper around gorm.DB
type DB struct {
	*gorm.DB
}

// NewDatabase initializes and returns a new database connection
func NewDatabase(cfg *config.Config) (*DB, error) {
	// Attempt to open a connection to the database using GORM and the provided DSN
	db, err := gorm.Open(postgres.Open(cfg.Db.Dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to the database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Retrieve the underlying sql.DB object to perform connection checks
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to obtain database object for connection check", zap.Error(err))
		return nil, fmt.Errorf("unable to access database object: %w", err)
	}

	// Ping the database to ensure the connection is alive
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Database connection test failed", zap.Error(err))
		return nil, fmt.Errorf("failed to establish a connection to the database: %w", err)
	}

	fmt.Println("Successfully connected to the database!")
	logger.Info("Successfully connected to the database")

	return &DB{db}, nil
}
