package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Printf(format string, v ...any)
	Fatalf(format string, v ...any)
}

var globalLogger Logger = &DefaultLogger{}

func Init(format string) {
	switch strings.ToLower(format) {
	case "structured":
		globalLogger = NewStructuredLogger()
	case "plain":
		globalLogger = &DefaultLogger{}
	default:
		panic(fmt.Sprintf("invalid LOG_FORMAT: %q. Must be 'structured' or 'plain'", format))
	}
}

type DefaultLogger struct{}

func (l *DefaultLogger) Info(msg string, args ...any) {
	log.Printf("INFO: %s %v", msg, args)
}

func (l *DefaultLogger) Error(msg string, args ...any) {
	log.Printf("ERROR: %s %v", msg, args)
}

func (l *DefaultLogger) Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func (l *DefaultLogger) Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}

type StructuredLogger struct {
	logger *slog.Logger
}

func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (l *StructuredLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *StructuredLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *StructuredLogger) Printf(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l *StructuredLogger) Fatalf(format string, v ...any) {
	l.logger.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	globalLogger.Error(msg, args...)
}

func Printf(format string, v ...any) {
	globalLogger.Printf(format, v...)
}

func Fatalf(format string, v ...any) {
	globalLogger.Fatalf(format, v...)
}

func With(args ...any) Logger {
	if sl, ok := globalLogger.(*StructuredLogger); ok {
		return &StructuredLogger{logger: sl.logger.With(args...)}
	}
	return globalLogger
}
