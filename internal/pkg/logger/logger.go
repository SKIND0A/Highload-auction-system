package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type ctxKey string

const (
	keyRequestID ctxKey = "request_id"
	keyUserID    ctxKey = "user_id"
	keyLotID     ctxKey = "lot_id"
	keyTxID      ctxKey = "tx_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

func WithLotID(ctx context.Context, lotID string) context.Context {
	return context.WithValue(ctx, keyLotID, lotID)
}

func WithTxID(ctx context.Context, txID string) context.Context {
	return context.WithValue(ctx, keyTxID, txID)
}

func Init() {
	env := os.Getenv("ENV")
	logLevelStr := os.Getenv("LOG_LEVEL")

	var level slog.Level
	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.ToLower(env) == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}


func extractAttrs(ctx context.Context) []any {
	var attrs []any

	if v, ok := ctx.Value(keyRequestID).(string); ok && v != "" {
		attrs = append(attrs, "request_id", v)
	}
	if v, ok := ctx.Value(keyUserID).(string); ok && v != "" {
		attrs = append(attrs, "user_id", v)
	}
	if v, ok := ctx.Value(keyLotID).(string); ok && v != "" {
		attrs = append(attrs, "lot_id", v)
	}
	if v, ok := ctx.Value(keyTxID).(string); ok && v != "" {
		attrs = append(attrs, "tx_id", v)
	}

	return attrs
}

func FromContext(ctx context.Context) *slog.Logger {
	attrs := extractAttrs(ctx)
	if len(attrs) == 0 {
		return slog.Default()
	}
	return slog.Default().With(attrs...)
}

// helper context
func Debug(ctx context.Context, msg string, args ...any) {
	allArgs := append(extractAttrs(ctx), args...)
	slog.Default().DebugContext(ctx, msg, allArgs...)
}

func Info(ctx context.Context, msg string, args ...any) {
	allArgs := append(extractAttrs(ctx), args...)
	slog.Default().InfoContext(ctx, msg, allArgs...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	allArgs := append(extractAttrs(ctx), args...)
	slog.Default().WarnContext(ctx, msg, allArgs...)
}

func Error(ctx context.Context, msg string, args ...any) {
	allArgs := append(extractAttrs(ctx), args...)
	slog.Default().ErrorContext(ctx, msg, allArgs...)
}

// aliases
func InfoCtx(ctx context.Context, msg string, args ...any)  { Info(ctx, msg, args...) }
func WarnCtx(ctx context.Context, msg string, args ...any)  { Warn(ctx, msg, args...) }
func ErrorCtx(ctx context.Context, msg string, args ...any) { Error(ctx, msg, args...) }
func DebugCtx(ctx context.Context, msg string, args ...any) { Debug(ctx, msg, args...) }