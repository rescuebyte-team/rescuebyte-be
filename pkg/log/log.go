package log

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Fields map[string]any

type reqKey struct{}

type userKey struct{}

type UserInfo struct {
	ID   string
	Role string
	NIM  string
}

type CtxLogger struct {
	logger *zerolog.Logger
}

var RequestIDKey = reqKey{}

var UserKey = userKey{}

func WithUser(ctx context.Context, user UserInfo) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

var logger zerolog.Logger

func GetLogger() *zerolog.Logger {
	return &logger
}

func init() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}

	multi := zerolog.MultiLevelWriter(consoleWriter)

	logger = zerolog.New(multi).With().Timestamp().Logger()
}

func Ctx(ctx context.Context) *CtxLogger {
	if ctx == nil {
		return &CtxLogger{logger: &logger}
	}

	reqID, _ := ctx.Value(RequestIDKey).(string)
	user, hasUser := ctx.Value(UserKey).(UserInfo)
	if reqID == "" && !hasUser {
		return &CtxLogger{logger: &logger}
	}

	builder := logger.With()
	if reqID != "" {
		builder = builder.Str("request_id", reqID)
	}

	if hasUser {
		builder = builder.Str("user_id", user.ID)
		if user.Role != "" {
			builder = builder.Str("role", user.Role)
		}
		if user.NIM != "" {
			builder = builder.Str("nim", user.NIM)
		}
	}

	child := builder.Logger()
	return &CtxLogger{logger: &child}
}

func (l *CtxLogger) Trace(fields Fields, msg string) {
	event := l.logger.Trace()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (l *CtxLogger) Debug(fields Fields, msg string) {
	event := l.logger.Debug()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (l *CtxLogger) Info(fields Fields, msg string) {
	event := l.logger.Info()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (l *CtxLogger) Warn(fields Fields, msg string) {
	event := l.logger.Warn()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (l *CtxLogger) Error(fields Fields, msg string) {
	event := l.logger.Error()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (l *CtxLogger) Fatal(fields Fields, msg string) {
	event := l.logger.Fatal()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (l *CtxLogger) Panic(fields Fields, msg string) {
	event := l.logger.Panic()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Info(fields Fields, msg string) {
	event := logger.Info()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Warn(fields Fields, msg string) {
	event := logger.Warn()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Error(fields Fields, msg string) {
	event := logger.Error()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Fatal(fields Fields, msg string) {
	event := logger.Fatal()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func Panic(fields Fields, msg string) {
	event := logger.Panic()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}
