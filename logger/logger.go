package logger

import (
	"context"
	"fmt"
)

var (
	internalLogger Logger
)

func init() {
	// Initialize the logger with default configuration
	var err error
	internalLogger, err = NewLogger(&Config{AppName: "default", Build: "dev"})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
}

func (c Config) InitiateLogger() error {
	var err error
	internalLogger, err = NewLogger(&c)
	return err
}

func AddMiddleLayers(middlelayers ...MiddleLayer) {
	internalLogger.AddMiddleLayers(middlelayers...)
}

func Info(ctx context.Context, format string, a ...any) {
	internalLogger.Info(ctx, format, a...)
}

func Infof(ctx context.Context, format string, fields *Fields) {
	internalLogger.Infof(ctx, format, fields)
}

func Infow(ctx context.Context, message string, fs *Fields) {
	internalLogger.Infof(ctx, message, fs)
}

func Error(ctx context.Context, format string, a ...any) {
	internalLogger.Error(ctx, format, a...)
}

func Warn(ctx context.Context, format string, a ...any) {
	internalLogger.Warn(ctx, format, a...)
}

func Debug(ctx context.Context, format string, a ...any) {
	internalLogger.Debug(ctx, format, a...)
}

func Panic(ctx context.Context, format string, a ...any) {
	internalLogger.Panic(ctx, format, a...)
}

func Fatal(ctx context.Context, format string, a ...any) {
	internalLogger.Fatal(ctx, format, a...)
}

type BaseLogger interface {
	Info(ctx context.Context, format string, a ...any)
	Error(ctx context.Context, format string, a ...any)
	Warn(ctx context.Context, format string, a ...any)
	Debug(ctx context.Context, format string, a ...any)
	Panic(ctx context.Context, format string, a ...any)
	Fatal(ctx context.Context, format string, a ...any)
}

type Logger interface {
	BaseLogger
	Infof(ctx context.Context, format string, fields *Fields)
	Debugf(ctx context.Context, format string, fields *Fields)
	Errorf(ctx context.Context, format string, fields *Fields)
	Warnf(ctx context.Context, format string, fields *Fields)
	Panicf(ctx context.Context, format string, fields *Fields)
	Fatalf(ctx context.Context, format string, fields *Fields)
	AddMiddleLayers(middlelayers ...MiddleLayer)
	ReplaceMiddleLayers(middlelayers ...MiddleLayer)
	ChangeConfig(config *Config) error
}
