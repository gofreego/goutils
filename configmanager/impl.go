package configmanager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofreego/goutils/cache/memory"
)

type configManagerConfig struct {
	CacheTimeoutMinutes int `name:"cache_timeout_minutes" type:"number" description:"cache timeout in minutes for config manager", required:"true"`
}

func (c *configManagerConfig) Key() string {
	return "config-manager-config"
}

type configManager struct {
	repository Repository
	cache      *memory.Cache
	config     *configManagerConfig
}

func newConfigManager(ctx context.Context, repository Repository) (*configManager, error) {
	manager := &configManager{
		repository: repository,
		cache:      memory.NewCache(),
		config: &configManagerConfig{
			CacheTimeoutMinutes: 5,
		},
	}

	err := manager.RegisterConfig(ctx, manager.config)
	if err != nil {
		return nil, err
	}
	return manager, nil
}

// RegisterConfig will register config and setup a UI for it. It will also validate the config.
func (manager *configManager) RegisterConfig(ctx context.Context, cfg config) error {
	// validate config
	cfgStr, err := marshal(ctx, cfg)
	if err != nil {
		return err
	}

	// check if config is already present in the repository
	value, err := manager.repository.GetConfig(ctx, cfg.Key())
	if err != nil {
		return err
	}

	// if config is not present in the repository, save it
	if value == "" {
		if err := manager.repository.SaveConfig(ctx, cfg.Key(), cfgStr); err != nil {
			return err
		}
		// save the config in manager

	}

	// if config is present in the repository, unmarshal it
	err = unmarshal(ctx, value, cfg)
	if err != nil {
		return fmt.Errorf("config already present in repository is not compatible with given object,Err: %s", err.Error())
	}
	// save the config in manager
	return nil
}

// RegisterRoute registers routes for the configuration manager.
// register routes with /configs/* endpoints
func (c *configManager) RegisterRoute(ctx context.Context, registerFunc RouteRegistrar) error {
	// setup ui
	if err := registerFunc(http.MethodGet, "/configs/ui", c.handleUI); err != nil {
		return err
	}
	// setup get config
	if err := registerFunc(http.MethodGet, "/configs/config/{key}", c.handleGetConfig); err != nil {
		return err
	}

	//setup save config
	if err := registerFunc(http.MethodPost, "/configs/config", c.handleSaveConfig); err != nil {
		return err
	}

	// setup get all configs
	if err := registerFunc(http.MethodGet, "/configs/metadata", c.handleGetConfigMetadata); err != nil {
		return err
	}
	return nil
}

// Get reads configs from the repository and sets them in the config.
func (c *configManager) Get(ctx context.Context, cfg config) error {

	return nil
}

func (manager *configManager) saveConfig(ctx context.Context, cfg config) error {
	cfgStr, err := marshal(ctx, cfg)
	if err != nil {
		return err
	}

	if err := manager.cache.SetWithTimeout(ctx, cfg.Key(), cfgStr, time.Minute*time.Duration(manager.config.CacheTimeoutMinutes)); err != nil {
		return fmt.Errorf("failed to save config in cache: %w", err)
	}

	if err := manager.repository.SaveConfig(ctx, cfg.Key(), cfgStr); err != nil {
		return err
	}
	return nil
}

func (manager *configManager) getConfig(ctx context.Context, key string) (string, error) {
	value, err := manager.cache.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get config from cache: %w", err)
	}
	if value != "" {
		return value, nil
	}

	value, err = manager.repository.GetConfig(ctx, key)
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("config not found")
	}

	if err := manager.cache.SetWithTimeout(ctx, key, value, time.Minute*time.Duration(manager.config.CacheTimeoutMinutes)); err != nil {
		return "", fmt.Errorf("failed to save config in cache: %w", err)
	}
	return value, nil
}
