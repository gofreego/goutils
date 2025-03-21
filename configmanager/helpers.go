package configmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/gofreego/goutils/logger"
)

// saveConfig saves/update the config in cache and repository
func (manager *configManager) saveConfig(ctx context.Context, cfg *Config) error {

	if err := manager.cache.SetWithTimeout(ctx, cfg.Key, cfg, time.Minute*time.Duration(manager.config.CacheTimeoutMinutes)); err != nil {
		return fmt.Errorf("failed to save config in cache: %w", err)
	}

	if err := manager.repository.SaveConfig(ctx, cfg); err != nil {
		return err
	}
	return nil
}

// getConfig gets the config from cache or repository. returns ErrConfigNotFound if config is not found
func (manager *configManager) getConfig(ctx context.Context, key string) (*Config, error) {
	var cfg Config
	err := manager.cache.GetV(ctx, key, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from cache: %w", err)
	}
	if cfg.Key != "" {
		return &cfg, nil
	}

	repoCfg, err := manager.repository.GetConfig(ctx, key)
	if err != nil {
		return nil, err
	}
	if repoCfg == nil {
		return nil, ErrConfigNotFound
	}

	if err := manager.cache.SetWithTimeout(ctx, key, repoCfg, time.Minute*time.Duration(manager.config.CacheTimeoutMinutes)); err != nil {
		logger.Error(ctx, "failed to save config in cache: %v", err)
		return nil, fmt.Errorf("failed to save config in cache: %w", err)
	}
	return repoCfg, nil
}

func (manager *configManager) addConfigToMap(_ context.Context, cfg config) {
	manager.registeredConfigs[cfg.Key()] = cfg
}
