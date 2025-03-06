package file

import (
	"context"
	"os"

	"github.com/gofreego/goutils/configutils/common"
	"github.com/gofreego/goutils/logger"
)

// Config represents the configuration for file reader
// Path is the base path for the file reader
type Config struct {
	Path string
}

type FileConfigReader struct {
	cfg *Config
}

func NewFileConfigReader(config *Config) *FileConfigReader {
	return &FileConfigReader{cfg: config}
}

// Read reads the configuration from the file
// path : path in file to read the configuration from
// conf : configuration object to unmarshal the data into
// configFormat : format of the configuration data
// returns error if any
// returns nil if successful
func (a *FileConfigReader) Read(ctx context.Context, path string, conf any, configFormat ...common.ConfigFormatType) error {
	path = a.cfg.Path + path
	bytes, err := os.ReadFile(path)
	if err != nil {
		logger.Error(ctx, "error reading from file : %v", err)
		return err
	}

	err = common.Unmarshal(bytes, conf, configFormat...)
	if err != nil {
		logger.Error(ctx, "error unmarshalling for path: %s, err: %v", path, err)
		return err
	}
	return nil
}
