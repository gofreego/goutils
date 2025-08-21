package configutils

import (
	"context"

	"github.com/gofreego/goutils/configutils/common"
	"github.com/gofreego/goutils/configutils/impls/consul"
	"github.com/gofreego/goutils/configutils/impls/file"
	"github.com/gofreego/goutils/configutils/impls/zookeeper"
	"github.com/gofreego/goutils/logger"
)

// Config represents the configuration for the config reader.
// Name is the type of the config reader, Expect one of consul, zookeeper, database, file
// Format is the format of the configuration data, Expect one of json, yaml
// Consul is the configuration for consul reader
// Zookeeper is the configuration for zookeeper reader
// Database is the configuration for database reader
// File is the configuration for file reader
type Config struct {
	Name      common.ConfigReaderName
	Format    common.ConfigFormatType
	Consul    consul.Config
	Zookeeper zookeeper.Config
	File      file.Config
}

type ConfigReader interface {
	// Read reads the configuration from the given path and unmarshals it into the given conf.
	// path : path in the configuration store to read the configuration from
	// conf : configuration object to unmarshal the data into
	// configFormat : format of the configuration data
	// returns error if any
	// returns nil if successful
	Read(ctx context.Context, path string, conf any, configFormat ...common.ConfigFormatType) error
	Update(ctx context.Context, path string, conf any, configFormat ...common.ConfigFormatType) error
}

// NewConfigReader creates a new config reader based on the given configuration.
// it is recommended to use config reader for reading configuration from consul, zookeeper, database on production and not use file.
func NewConfigReader(ctx context.Context, conf *Config) (ConfigReader, error) {
	if conf == nil {
		return nil, common.ErrInvalidConfig
	}
	switch conf.Name {
	case common.ConsulConfigReader:
		return consul.NewConsulConfigReader(ctx, &conf.Consul)
	case common.ZookeeperConfigReader:
		return zookeeper.NewZookeeperReader(ctx, &conf.Zookeeper)
	case common.FileConfigReader:
		return file.NewFileConfigReader(&conf.File), nil
	default:
		return nil, common.ErrInvalidConfigReaderName
	}
}

type tConfig interface {
	GetReaderConfig() *Config
}

// ReadConfig reads the configuration from the given path and unmarshals it into the given conf.
// It reads the configuration from the file
// If conf implements tConfig, it reads the configuration from the reader specified in the configuration.
// It supports json and yaml formats.
// If no format is provided, it defaults to yaml.
// It returns error if any.
// It returns nil if successful.
// it is recommended to use config reader for reading configuration from consul, zookeeper, database on production and not use file.
func ReadConfig(ctx context.Context, path string, conf any) error {
	if conf == nil {
		return common.ErrInvalidConfig
	}

	err := file.NewFileConfigReader(&file.Config{Path: path}).Read(ctx, "", conf, common.ConfigFormatYAML)
	if err != nil {
		return err
	}
	if conf, ok := conf.(tConfig); ok {
		reader, err := NewConfigReader(ctx, conf.GetReaderConfig())
		if err != nil {
			return err
		}
		frmt := conf.GetReaderConfig().Format
		if frmt == "" {
			frmt = common.ConfigFormatYAML
		}
		return reader.Read(ctx, "", conf, frmt)
	}
	return nil
}

// LogConfig logs the configuration.
// It logs the configuration in yaml format if config does not implement tConfig or the format is not provided.
func LogConfig(ctx context.Context, conf any) {
	if conf1, ok := conf.(tConfig); ok {
		frmt := conf1.GetReaderConfig().Format
		if frmt == "" {
			frmt = common.ConfigFormatYAML
		}
		bytes, err := common.Marshal(conf, frmt)
		if err != nil {
			return
		}
		logger.Info(ctx, "config \n %s", string(bytes))
	} else {
		bytes, err := common.Marshal(conf, common.ConfigFormatYAML)
		if err != nil {
			return
		}
		logger.Info(ctx, "config \n %s", string(bytes))
	}
}
