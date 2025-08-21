package main

import (
	"context"

	"github.com/gofreego/goutils/configutils"
)

type Config struct {
	Name        string
	Reader      configutils.Config
	Application string
}

// GetReaderConfig implements configutils.tConfig.
func (c *Config) GetReaderConfig() *configutils.Config {
	return &c.Reader
}

func main() {
	var conf Config
	configutils.LogConfig(context.TODO(), conf)
	err := configutils.ReadConfig(context.TODO(), "dev.yaml", &conf)
	if err != nil {
		panic(err)
	}
	configutils.LogConfig(context.Background(), &conf)
}
