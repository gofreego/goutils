package logger

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	zapLogger    *zap.Logger
	appNameField zap.Field
	middleLayers []MiddleLayer
}

func NewLogger(config *Config, middleLayers ...MiddleLayer) (Logger, error) {
	var err error
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = timeKey
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var zapConfig zap.Config
	if config.Build == "prod" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}
	if config.Level != "" {
		level, found := logLevelToZapLevelMap[config.Level]
		if !found {
			return nil, errors.New("invalid log level in config")
		}
		zapConfig.Level = zap.NewAtomicLevelAt(level)
	}
	zapConfig.DisableStacktrace = true

	zapConfig.EncoderConfig = encoderConfig
	appNameField := zap.Field{Key: "App", Type: zapcore.StringType, String: "default"}
	if config.skipLevels <= 0 {
		config.skipLevels = 1
	}
	zapLogger, err := zapConfig.Build(zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(config.skipLevels))
	if err != nil {
		return nil, err
	}
	return &logger{
		zapLogger:    zapLogger,
		appNameField: appNameField,
		middleLayers: middleLayers,
	}, nil
}

func (l *logger) AddMiddleLayers(middlelayers ...MiddleLayer) {
	l.middleLayers = append(l.middleLayers, middlelayers...)
}

func (l *logger) ReplaceMiddleLayers(middlelayers ...MiddleLayer) {
	l.middleLayers = middlelayers
}

func (l *logger) ChangeConfig(config *Config) error {
	var err error
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = timeKey
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var zapConfig zap.Config
	if config.Build == "prod" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}
	if config.Level != "" {
		level, found := logLevelToZapLevelMap[config.Level]
		if !found {
			return errors.New("invalid log level in config")
		}
		zapConfig.Level = zap.NewAtomicLevelAt(level)
	}
	zapConfig.DisableStacktrace = true

	zapConfig.EncoderConfig = encoderConfig
	l.appNameField.String = config.AppName
	if config.skipLevels <= 0 {
		config.skipLevels = 1
	}
	l.zapLogger, err = zapConfig.Build(zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(config.skipLevels))
	return err
}

func (l *logger) executeMiddleLayers(ctx context.Context, msg string, fields *Fields) (context.Context, string, *Fields) {
	for _, layer := range l.middleLayers {
		ctx, msg, fields = layer(ctx, msg, fields)
	}
	return ctx, msg, fields
}

func (l *logger) Info(ctx context.Context, format string, a ...any) {
	_, msg, fields := l.executeMiddleLayers(ctx, fmt.Sprintf(format, a...), &Fields{fields: []zap.Field{l.appNameField}})
	l.zapLogger.Info(msg, fields.fields...)
}

func (l *logger) Debug(ctx context.Context, format string, a ...any) {
	_, msg, fields := l.executeMiddleLayers(ctx, fmt.Sprintf(format, a...), &Fields{fields: []zap.Field{l.appNameField}})
	l.zapLogger.Debug(msg, fields.fields...)
}

func (l *logger) Error(ctx context.Context, format string, a ...any) {
	_, msg, fields := l.executeMiddleLayers(ctx, fmt.Sprintf(format, a...), &Fields{fields: []zap.Field{l.appNameField}})
	l.zapLogger.Error(msg, fields.fields...)
}

func (l *logger) Warn(ctx context.Context, format string, a ...any) {
	_, msg, fields := l.executeMiddleLayers(ctx, fmt.Sprintf(format, a...), &Fields{fields: []zap.Field{l.appNameField}})
	l.zapLogger.Warn(msg, fields.fields...)
}

func (l *logger) Panic(ctx context.Context, format string, a ...any) {
	_, msg, fields := l.executeMiddleLayers(ctx, fmt.Sprintf(format, a...), &Fields{fields: []zap.Field{l.appNameField}})
	l.zapLogger.Panic(msg, fields.fields...)
}

func (l *logger) Fatal(ctx context.Context, format string, a ...any) {
	_, msg, fields := l.executeMiddleLayers(ctx, fmt.Sprintf(format, a...), &Fields{fields: []zap.Field{l.appNameField}})
	l.zapLogger.Fatal(msg, fields.fields...)
}

func (l *logger) Infof(ctx context.Context, msg string, fs *Fields) {
	fs.fields = append(fs.fields, l.appNameField)
	_, msg, fields := l.executeMiddleLayers(ctx, msg, fs)
	l.zapLogger.Info(msg, fields.fields...)
}

func (l *logger) Debugf(ctx context.Context, msg string, fields *Fields) {
	fields.fields = append(fields.fields, l.appNameField)
	_, msg, fields = l.executeMiddleLayers(ctx, msg, fields)
	l.zapLogger.Debug(msg, fields.fields...)
}

func (l *logger) Warnf(ctx context.Context, message string, fs *Fields) {
	fs.fields = append(fs.fields, l.appNameField)
	_, msg, fields := l.executeMiddleLayers(ctx, message, fs)
	l.zapLogger.Warn(msg, fields.fields...)
}

func (l *logger) Errorf(ctx context.Context, msg string, fields *Fields) {
	fields.fields = append(fields.fields, l.appNameField)
	_, msg, fields = l.executeMiddleLayers(ctx, msg, fields)
	l.zapLogger.Error(msg, fields.fields...)
}

func (l *logger) Fatalf(ctx context.Context, msg string, fields *Fields) {
	fields.fields = append(fields.fields, l.appNameField)
	_, msg, fields = l.executeMiddleLayers(ctx, msg, fields)
	l.zapLogger.Fatal(msg, fields.fields...)
}
func (l *logger) Panicf(ctx context.Context, msg string, fields *Fields) {
	fields.fields = append(fields.fields, l.appNameField)
	_, msg, fields = l.executeMiddleLayers(ctx, msg, fields)
	l.zapLogger.Panic(msg, fields.fields...)
}
