package consul

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/gofreego/goutils/logger"
	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ReadFromConsul bool
	Address        string
	Token          string
	Path           string
	RefreshInSecs  int
}

type ConsulReader struct {
	conf *Config
	kv   *api.KV
}

func NewConsulReader(ctx context.Context, config *Config) (*ConsulReader, error) {
	client, err := api.NewClient(&api.Config{
		Address: config.Address,
		Token:   config.Token,
	})
	if err != nil {
		logger.Error(ctx, fmt.Sprintln("Error creating consul client : ", err))
		return nil, err
	}
	return &ConsulReader{kv: client.KV(), conf: config}, nil
}

func (a *ConsulReader) read(ctx context.Context, conf any, path string) error {
	logger.Debug(ctx, "Reading from consul path : %s", path)
	data, _, err := a.kv.Get(path, nil)
	if err != nil {
		logger.Error(ctx, "Error reading from zookeeper : %v", err)
		return err
	}
	if data == nil {
		return nil
	}
	err = yaml.Unmarshal(data.Value, conf)
	if err != nil {
		logger.Error(ctx, "Error unmarshalling yaml for path: %s, data : %v", path, err)
		return err
	}
	return nil
}

func (a *ConsulReader) readRecursively(ctx context.Context, conf any, path string) error {
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

			if !fieldValue.CanAddr() {
				logger.Warn(ctx, "filed %s has `path` but its not exported path : %s", fieldValue.Kind().String(), configPath)
				continue
			}
			ptr := fieldValue.Addr()
			var err error
			if configPath != "" {
				err = a.read(ctx, ptr.Interface(), path+"/"+configPath)
			} else {
				err = a.read(ctx, ptr.Interface(), path+"/"+field.Name)
			}

			if err != nil {
				return err
			}
		}
	}

	return a.read(ctx, conf, path)
}

func (a *ConsulReader) Read(ctx context.Context, conf any) error {
	err := a.readRecursively(ctx, conf, a.conf.Path)
	if err != nil {
		return err
	}
	go a.refresh(ctx, conf)
	return nil
}

func (a *ConsulReader) refresh(ctx context.Context, conf any) {
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
