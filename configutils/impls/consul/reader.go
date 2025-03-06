package consul

import (
	"context"

	"github.com/gofreego/goutils/configutils/common"
	"github.com/gofreego/goutils/logger"
	"github.com/hashicorp/consul/api"
)

// Config : configuration for consul
// Address : address of the consul server
// Token : token for authentication
// Path : path in consul to read the configuration from
type Config struct {
	Address string
	Token   string
	Path    string
}

type ConsulConfigReader struct {
	kv  *api.KV
	cfg *Config
}

// NewConsulConfigReader creates a new consul configuration reader
func NewConsulConfigReader(ctx context.Context, config *Config) (*ConsulConfigReader, error) {
	client, err := api.NewClient(&api.Config{
		Address: config.Address,
		Token:   config.Token,
	})
	if err != nil {
		logger.Error(ctx, "Error creating consul client : %v", err)
		return nil, err
	}
	return &ConsulConfigReader{kv: client.KV(), cfg: config}, nil
}

// Read reads the configuration from consul
// path : path in consul to read the configuration from
// conf : configuration object to unmarshal the data into
// configFormat : format of the configuration data
// returns error if any
// returns nil if successful
func (a *ConsulConfigReader) Read(ctx context.Context, path string, conf any, configFormat ...common.ConfigFormatType) error {
	path = a.cfg.Path + path
	logger.Debug(ctx, "Reading from consul path : %s", path)
	data, _, err := a.kv.Get(path, nil)
	if err != nil {
		logger.Error(ctx, "Error reading from zookeeper : %v", err)
		return err
	}
	if data == nil {
		return nil
	}

	err = common.Unmarshal(data.Value, conf, configFormat...)
	if err != nil {
		logger.Error(ctx, "Error unmarshalling yaml for path: %s, data : %v", path, err)
		return err
	}
	return nil
}
