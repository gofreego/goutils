package main

import (
	"context"

	"github.com/gofreego/goutils/configutils/zookeeper"
	"github.com/gofreego/goutils/logger"
)

type Level1 struct {
	V1     int    `yaml:"V1"`
	V2     string `yaml:"V2"`
	Level2 Level2 `yaml:"Level2" path:"level2"`
}

type Level2 struct {
	V1 int    `yaml:"V1"`
	V2 string `yaml:"V2"`
}

type Config struct {
	Name   string `yaml:"Name"`
	Int    int    `yaml:"Int"`
	Level1 Level1 `yaml:"Level1" path:"level1" children:"true"`
}

func main() {
	ctx := context.Background()
	zk, err := zookeeper.NewZookeeperReader(ctx, &zookeeper.Config{
		ReadFromZookeeper: true,
		Address:           "localhost:2181",
		Path:              "/configs/test",
	})
	if err != nil {
		logger.Panic(ctx, "failed to load zk , %v", err)
	}
	var conf Config
	err = zk.Read(ctx, &conf)
	if err != nil {
		logger.Error(ctx, "failed to read from zk %v", err)
	}
	logger.Info(ctx, "%v", conf)
}
