package zookeeper

import (
	"context"
	"reflect"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/gofreego/goutils/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ReadFromZookeeper bool
	Address           string
	Path              string
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

func (a *ZookeeperReader) read(ctx context.Context, conf any, path string) error {
	data, _, err := a.conn.Get(path)
	if err != nil {
		logger.Error(ctx, "Error reading from zookeeper : %v", err)
		return err
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		logger.Error(ctx, "Error unmarshalling yaml for path: %s, data : %v", path, err)
		return err
	}
	return nil
}

func (a *ZookeeperReader) readRecursively(ctx context.Context, conf any, path string) error {
	v := reflect.ValueOf(conf)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {

			field := t.Field(i)
			fieldValue := v.Field(i)
			if !fieldValue.CanInterface() {
				continue
			}
			childrenTag := field.Tag.Get("children")
			configPath := field.Tag.Get("path")

			if childrenTag == "true" {
				ptr := fieldValue.Addr()
				if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
					ptr = fieldValue
				}
				if err := a.readRecursively(ctx, ptr.Interface(), path+"/"+configPath); err != nil {
					return err
				}
			}

			if configPath != "" {
				var ptr reflect.Value
				if fieldValue.CanAddr() {
					ptr = fieldValue.Addr()
				} else {
					logger.Warn(ctx, "filed %s has `path` but its not exported path : %s", fieldValue.Kind().String(), configPath)
					continue
				}
				a.read(ctx, ptr.Interface(), path+"/"+configPath)
			}
		}
	}

	return a.read(ctx, conf, path)
}

func (a *ZookeeperReader) Read(ctx context.Context, conf any) error {
	err := a.readRecursively(ctx, conf, a.conf.Path)
	if err != nil {
		return err
	}
	go a.refresh(ctx, conf)
	return nil
}

func (a *ZookeeperReader) refresh(ctx context.Context, conf any) {
	if a.conf.RefreshInSecs == 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(a.conf.RefreshInSecs) * time.Second)
	for {
		<-ticker.C
		err := a.readRecursively(ctx, conf, a.conf.Path)
		if err != nil {
			logger.Error(ctx, "failed to refresh config , %v", err)
		}
	}
}
