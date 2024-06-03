package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofreego/goutils/logger"
	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	ReadFromConsul bool
	Address        string
	Token          string
	Path           string
	LookForChange  bool
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

func insertMap(m map[string]any, keys []string, len int, value []byte) {
	if len == 1 {
		var v any
		err := json.Unmarshal(value, &v)
		if err != nil {
			m[keys[0]] = string(value)
			return
		}
		m[keys[0]] = v
		return
	}
	if t, ok := m[keys[0]].(map[string]any); ok {
		insertMap(t, keys[1:], len-1, value)
	} else {
		m[keys[0]] = make(map[string]any)
		insertMap(m[keys[0]].(map[string]any), keys[1:], len-1, value)
	}
}

// this function will read the tags and fetch values from consul
func (a *ConsulReader) Read(conf any) error {
	pair, _, err := a.kv.List(a.conf.Path, nil)
	if err != nil {
		return err
	}
	m := make(map[string]string)
	cfg := make(map[string]any)

	for _, kv := range pair {
		if len(kv.Value) == 0 {
			continue
		}
		_key := strings.ReplaceAll(kv.Key, a.conf.Path+"/", "")
		if _key != "" {
			m[_key] = strings.TrimSpace(string(kv.Value))
			insertMap(cfg, strings.Split(_key, "/"), len(strings.Split(_key, "/")), kv.Value)
		}
	}
	err = mapstructure.Decode(cfg, conf)
	if err != nil {
		logger.Error(context.Background(), fmt.Sprintln("Error decoding config : ", err))
		return err
	}

	go a.lookForChange(conf, m, cfg)
	return nil
}

func (a *ConsulReader) lookForChange(conf any, memory map[string]string, cfg map[string]any) error {
	if a.conf.RefreshInSecs == 0 {
		a.conf.RefreshInSecs = 30
	}
	for {
		time.Sleep(time.Duration(a.conf.RefreshInSecs) * time.Second)
		var changeDetected bool = false
		pair, _, err := a.kv.List(a.conf.Path, nil)
		if err != nil {
			logger.Error(context.Background(), fmt.Sprintln("Error reading consul : ", err))
			continue
		}
		for _, kv := range pair {
			if len(kv.Value) == 0 {
				continue
			}
			_key := strings.ReplaceAll(kv.Key, a.conf.Path+"/", "")
			if _key != "" && memory[_key] != strings.TrimSpace(string(kv.Value)) {
				memory[_key] = strings.TrimSpace(string(kv.Value))
				insertMap(cfg, strings.Split(_key, "/"), len(strings.Split(_key, "/")), kv.Value)
				logger.Info(context.Background(), "config change detected for %s", _key)
				changeDetected = true
			}
		}
		if changeDetected {
			err = mapstructure.Decode(cfg, conf)
			if err != nil {
				logger.Error(context.Background(), fmt.Sprintln("Error decoding config : ", err))
				continue
			}
		}
	}
}
