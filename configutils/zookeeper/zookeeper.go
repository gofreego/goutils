package zookeeper

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/gofreego/goutils/logger"
)

type Config struct {
	ReadFromZookeeper bool
	Address           string
	Path              string
	LookForChange     bool
	RefreshInSecs     int
}

type ZookeeperReader struct {
	conf *Config
	conn *zk.Conn
}

func NewZookeeperReader(ctx context.Context, config *Config) (*ZookeeperReader, error) {
	conn, _, err := zk.Connect([]string{config.Address}, time.Second)
	if err != nil {
		logger.Error(ctx, "Error connecting to zookeeper : %v", err)
		return nil, err
	}
	return &ZookeeperReader{conf: config, conn: conn}, nil
}

func (a *ZookeeperReader) Read(ctx context.Context, conf any) error {
	data, _, err := a.conn.Get(a.conf.Path)
	if err != nil {
		logger.Error(ctx, "Error reading from zookeeper : %v", err)
		return err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		logger.Error(ctx, "Error unmarshalling data : %v", err)
		return err
	}
	return nil
}
