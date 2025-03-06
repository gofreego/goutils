package zookeeper

import (
	"context"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/gofreego/goutils/configutils/common"
	"github.com/gofreego/goutils/logger"
)

// Config : configuration for zookeeper
// Address : address of the zookeeper server
// Path : path in zookeeper to read the configuration from
// Username : username for authentication
// Password : password for authentication
type Config struct {
	Address  string
	Path     string
	Username string
	Password string
}

type ZookeeperReader struct {
	conf *Config
	conn *zk.Conn
}

// NewZookeeperReader creates a new zookeeper configuration reader
// if username and password are provided, it adds authentication to the connection else it connects without authentication
func NewZookeeperReader(ctx context.Context, config *Config) (*ZookeeperReader, error) {
	conn, _, err := zk.Connect([]string{config.Address}, time.Second)
	if err != nil {
		logger.Error(ctx, "Error connecting to zookeeper : %v", err)
		return nil, err
	}

	if config.Username != "" && config.Password != "" {
		err = conn.AddAuth("digest", []byte(config.Username+":"+config.Password))
		if err != nil {
			logger.Error(ctx, "Error adding authentication to zookeeper : %v", err)
			return nil, err
		}
	}

	return &ZookeeperReader{conf: config, conn: conn}, nil
}

// Read reads the configuration from zookeeper
// path : path in zookeeper to read the configuration from
// conf : configuration object to unmarshal the data into
// configFormat : format of the configuration data
// returns error if any
// returns nil if successful
func (a *ZookeeperReader) Read(ctx context.Context, path string, conf any, configFormat ...common.ConfigFormatType) error {
	data, _, err := a.conn.Get(path)
	if err != nil {
		logger.Error(ctx, "Error reading from zookeeper : %v", err)
		return err
	}
	err = common.Unmarshal(data, conf, configFormat...)
	if err != nil {
		logger.Error(ctx, "Error unmarshalling for path: %s, data : %v", path, err)
		return err
	}
	return nil
}

// Update updates the configuration in zookeeper
// path : path in zookeeper to update the configuration
// conf : configuration object to marshal the data from
// configFormat : format of the configuration data
// returns error if any
// returns nil if successful
func (a *ZookeeperReader) Update(ctx context.Context, path string, conf any, configFormat ...common.ConfigFormatType) error {
	bytes, err := common.Marshal(conf, configFormat...)
	if err != nil {
		logger.Error(ctx, "Error marshalling for path: %s, err: %v", path, err)
		return err
	}

	_, err = a.conn.Set(path, bytes, -1)
	if err != nil {
		logger.Error(ctx, "Error writing to zookeeper : %v", err)
		return err
	}
	return nil
}
