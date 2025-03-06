package database

import (
	"context"

	"github.com/gofreego/goutils/configutils/common"
)

type Config struct {
}

type DatabaseConfigReader struct {
}

func NewDatabaseReader(ctx context.Context, cfg *Config) (*DatabaseConfigReader, error) {
	return &DatabaseConfigReader{}, nil
}

func (d *DatabaseConfigReader) Read(ctx context.Context, path string, conf any, configFormat ...common.ConfigFormatType) error {
	panic("method unimplemented")
}
