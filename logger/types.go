package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Fields struct {
	fields []zapcore.Field
}

func NewFields() *Fields {
	return &Fields{
		fields: []zapcore.Field{},
	}
}

func (f *Fields) AddField(key string, value any) *Fields {
	f.fields = append(f.fields, zap.Any(key, value))
	return f
}

type Field struct {
	Key   string
	Value any
}

func NewField(key string, value any) *Field {
	return &Field{
		Key:   key,
		Value: value,
	}
}

type MiddleLayer func(ctx context.Context, msg string, fields *Fields) (context.Context, string, *Fields)

// Build : if prod it will set to prod else dev
type Config struct {
	AppName string   `yaml:"AppName" json:"AppName" name:"AppName" type:"string" description:"Application Name" required:"true"`
	Build   Build    `yaml:"Build" json:"Build" name:"Build" type:"choice" description:"Build Type" choices:"prod,dev"`
	Level   LogLevel `yaml:"Level" json:"Level" name:"Level" type:"choice" description:"Log Level" choices:"debug,info,warn,error,panic,fatal"`
}
