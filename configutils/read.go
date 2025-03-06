package configutils

import (
	"context"

	"github.com/gofreego/goutils/configutils/common"
	"github.com/gofreego/goutils/configutils/impls/consul"
	"github.com/gofreego/goutils/configutils/impls/database"
	"github.com/gofreego/goutils/configutils/impls/file"
	"github.com/gofreego/goutils/configutils/impls/zookeeper"
)

// Config represents the configuration for the config reader.
// Name is the type of the config reader, Expect one of consul, zookeeper, database, file
// Consul is the configuration for consul reader
// Zookeeper is the configuration for zookeeper reader
// Database is the configuration for database reader
// File is the configuration for file reader
type Config struct {
	Name      common.ConfigReaderName
	Consul    consul.Config
	Zookeeper zookeeper.Config
	Database  database.Config
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
}

// NewConfigReader creates a new config reader based on the given configuration.
func NewConfigReader(ctx context.Context, conf *Config) (ConfigReader, error) {
	if conf == nil {
		return nil, common.ErrInvalidConfig
	}
	switch conf.Name {
	case common.ConsulConfigReader:
		return consul.NewConsulConfigReader(ctx, &conf.Consul)
	case common.ZookeeperConfigReader:
		return zookeeper.NewZookeeperReader(ctx, &conf.Zookeeper)
	case common.DatabaseConfigReader:
		return database.NewDatabaseReader(ctx, &conf.Database)
	case common.FileConfigReader:
		return file.NewFileConfigReader(&conf.File), nil
	default:
		return nil, common.ErrInvalidConfigReaderName
	}
}
