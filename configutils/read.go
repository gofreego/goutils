package configutils

import (
	"context"
	"fmt"

	"github.com/gofreego/goutils/configutils/consul"
	"github.com/gofreego/goutils/configutils/zookeeper"
	"github.com/spf13/viper"
)

const (
	CONSUL    = "CONSUL"
	ZOOKEEPER = "ZOOKEEPER"
)

type Config struct {
	Name      string
	Consul    consul.Config
	Zookeeper zookeeper.Config
}

type config interface {
	GetReaderConfig() *Config
	GetServiceName() string
	GetEnv() string
}

func ReadConfig(ctx context.Context, filename string, conf any) error {
	// Read the YAML file
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(conf)
	if err != nil {
		return err
	}

	// check if config implements Config
	cfg, ok := conf.(config)
	if !ok {
		return nil
	}
	readerCfg := cfg.GetReaderConfig()

	readerCfg.Consul.Path += fmt.Sprintf("/%s/%s", cfg.GetEnv(), cfg.GetServiceName())
	readerCfg.Zookeeper.Path += fmt.Sprintf("/%s/%s", cfg.GetEnv(), cfg.GetServiceName())
	// Read from the agent
	return ReadFromAgent(ctx, readerCfg, conf)
}

func ReadFromAgent(ctx context.Context, cfg *Config, conf any) error {

	switch cfg.Name {
	case CONSUL:
		agent, err := consul.NewConsulReader(ctx, &cfg.Consul)
		if err != nil {
			return err
		}
		return agent.Read(ctx, conf)
	case ZOOKEEPER:
		agent, err := zookeeper.NewZookeeperReader(ctx, &cfg.Zookeeper)
		if err != nil {
			return err
		}
		return agent.Read(ctx, conf)
	}
	return nil
}
