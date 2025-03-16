package configmanager

import (
	"context"
	"net/http"

	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/datastructure"
)

type configManager struct {
	repository     Repository
	configsNameSet datastructure.Set[string]
}

func New(ctx context.Context) ConfigManager {
	return &configManager{}
}

// RegisterConfig will register config and setup a UI for it. It will also validate the config.
func (c *configManager) RegisterConfig(ctx context.Context, cfg config) error {
	value, err := c.repository.GetConfig(ctx, cfg.Key())
	if err != nil {
		return err
	}

	if value.Key != cfg.Key() {
		return customerrors.New(0, "invalid config returned from repository for key:%s", cfg.Key())
	}

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
	// todo complete this

	return nil
}
