package main

import (
	"context"

	"github.com/gofreego/goutils/configutils"
	"github.com/gofreego/goutils/configutils/consul"
	"github.com/gofreego/goutils/logger"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Name    string `yaml:"Name" path:"Name"`
	Age     int    `yaml:"Age"`
	Student struct {
		Name   string `yaml:"Name"`
		Age    int    `yaml:"Age"`
		School struct {
			Name     string `yaml:"Name"`
			Location string `yaml:"Location"`
		} `yaml:"School" path:"school"`
	} `yaml:"Student" children:"true" path:"student"`
}

func (c *Configuration) GetReaderConfig() *configutils.Config {
	return &configutils.Config{
		Name: "CONSUL",
		Consul: consul.Config{
			Address:       "localhost:8500",
			Path:          "/configs/test",
			RefreshInSecs: 1,
		},
	}
}

func main() {
	k := Configuration{}
	ctx := context.Background()
	err := configutils.ReadFromAgent(ctx, k.GetReaderConfig(), &k)
	if err != nil {
		panic(err)
	}
	bytes, _ := yaml.Marshal(k)
	logger.Info(ctx, "\n%s", bytes)
}
