package logger

import (
	l "log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger глобальный логгер.
var Logger *zap.Logger

type Env string

const (
	Development Env = "development"
	Production  Env = "production"
)

// InitLogger инициализация zap логгера с поддержкой разных режимов.
func InitLogger(env Env) {
	var cfg zap.Config

	switch env {
	case Production:
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	case Development:
		cfg = zap.NewDevelopmentConfig()
	default:
		l.Fatalf("Неизвестная среда: %s", env)
	}

	log, err := cfg.Build()
	if err != nil {
		l.Fatalf("Не удалось инициализировать логгер: %v", err)
	}

	Logger = log
	zap.ReplaceGlobals(log)
}

// Sync закрываем логгер (нужно вызывать при завершении приложения).
func Sync() {
	if Logger != nil {
		if err := Logger.Sync(); err != nil {
			Logger.Error("Ошибка при завершении логгера", zap.Error(err))
		}
	}
}

// Helper ф-ии для удобства.

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}
