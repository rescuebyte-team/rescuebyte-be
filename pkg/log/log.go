package log

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Fields map[string]interface{}

type reqKey struct{}

var RequestIDKey = reqKey{}

var logger zerolog.Logger

func init() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
}

func GetLogger() *zerolog.Logger {
	return &logger
}

func Ctx(ctx context.Context) *zerolog.Logger {
	if ctx == nil {
		return &logger
	}
	reqID, ok := ctx.Value(RequestIDKey).(string)
	if !ok || reqID == "" {
		return &logger
	}
	child := logger.With().Str("request_id", reqID).Logger()
	return &child
}

func Debug(fields Fields, msg string) {
	event := logger.Debug()
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
