package configutils

import (
	"context"

	"github.com/gofreego/goutils/configutils/consul"
	"github.com/spf13/viper"
)

type Config interface {
	GetConsulConfig() *consul.Config
}

func ReadConfig(ctx context.Context, filename string, config any) error {
	// Read the YAML file
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(config)
	if err != nil {
		return err
	}

	// check if config implements Config
	cfg, ok := config.(Config)
	if !ok {
		return nil
	}
	if cfg.GetConsulConfig().ReadFromConsul {
		agent, err := consul.NewConsulReader(ctx, cfg.GetConsulConfig())
		if err != nil {
			return err
		}
		err = agent.Read(config)
		if err != nil {
			return err
		}
	}
	return nil
}
