package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// DbConfig представляет конфигурацию базы данных.
type DbConfig struct {
	Dsn string
}

// AuthConfig представляет конфигурацию аутентификации.
type AuthConfig struct {
	Secret string
}

// Config представляет конфигурацию приложения.
type Config struct {
	Db           DbConfig
	Auth         AuthConfig
	Env          string
	DefaultLimit int
}

// NewConfig создаёт новый экземпляр конфигурации.
func NewConfig() *Config {
	// Загружаем .env (если есть)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: The .env file was not found. Default values are being used.")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
		Env: getEnv("APP_ENV", "development"),
	}
}

// getEnv возвращает значение переменной окружения или дефолтное значение
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
