package logger

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

type Env string

const (
	Development Env = "development"
	Production  Env = "production"
)

// InitLogger initializes the global zap logger with rotation, sampling, etc.
func InitLogger(env Env) {
	var cfg zap.Config
	switch env {
	case Production:
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(getLogLevel())
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.CallerKey = "caller"
		cfg.EncoderConfig.StacktraceKey = "stacktrace"
	case Development:
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(getLogLevel())
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.CallerKey = "caller"
	default:
		fmt.Fprintf(os.Stderr, "Unknown environment for logger: %s\n", env)
		os.Exit(1)
	}

	// Build cores: always stdout
	encoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	stdout := zapcore.Lock(os.Stdout)
	cores := []zapcore.Core{zapcore.NewCore(encoder, stdout, cfg.Level)}

	// In production also rotate to file
	if env == Production {
		if err := os.MkdirAll("logs", 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️ Failed to create logs dir: %v\n", err)
		}
		lj := &lumberjack.Logger{
			Filename:   "logs/vacancy-aggregator.log",
			MaxSize:    100, // MB
			MaxBackups: 7,
			MaxAge:     28, // days
			Compress:   true,
		}
		fileSink := zapcore.AddSync(lj)
		cores = append(cores, zapcore.NewCore(encoder, fileSink, cfg.Level))
	}

	combined := zapcore.NewTee(cores...)
	core := combined
	if env == Production {
		core = zapcore.NewSamplerWithOptions(
			combined,
			time.Second,
			100,
			100,
		)
	}

	logg := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "dev"
	}
	logg = logg.With(zap.String("service", "vacancy-aggregator"), zap.String("version", version))

	Logger = logg
	zap.ReplaceGlobals(logg)
}

// Sync flushes any buffered logs.
func Sync() {
	if Logger != nil {
		if err := Logger.Sync(); err != nil && !strings.Contains(err.Error(), "invalid argument") {
			fmt.Fprintf(os.Stderr, "⚠️ Error syncing logger: %v\n", err)
		}
	}
}

// getLogLevel reads LOG_LEVEL or defaults.
func getLogLevel() zapcore.Level {
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		var l zapcore.Level
		if err := l.UnmarshalText([]byte(lvl)); err == nil {
			return l
		}
		fmt.Fprintf(os.Stderr, "⚠️ Invalid LOG_LEVEL %q, defaulting to Info\n", lvl)
	}
	if os.Getenv("APP_ENV") == string(Production) {
		return zapcore.InfoLevel
	}
	return zapcore.DebugLevel
}

// RequestIDMiddleware injects a unique request ID into context and response header.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), "requestID", id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext retrieves the request ID from context, if any.
func FromContext(ctx context.Context) string {
	if id, ok := ctx.Value("requestID").(string); ok {
		return id
	}
	return ""
}

// ---- Convenience wrappers ----

// Info logs at InfoLevel.
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Debug logs at DebugLevel.
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Warn logs at WarnLevel.
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error logs at ErrorLevel.
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal logs at FatalLevel then os.Exit(1).
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
