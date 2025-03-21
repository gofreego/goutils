package configmanager

import (
	"encoding/json"
	"net/http"

	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/response"
)

// @title Config Manager API
// @version 1.0
// @description This is a sample API for managing configurations.
// @host localhost:8080
// @BasePath /
func (c *configManager) SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("not implemented"))
}

// Swagger doc
// @Summary UI
// @Description UI
// @Tags Config
// @Accept json
// @Produce html
// @Success 200 {string} string "UI"
// @Failure 400 {object} ErrorResponse
// @Router /configs/ui [get]
func (c *configManager) handleUI(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("not implemented"))
}

//Swagger doc
// @Summary Get config
// @Description Get config by key
// @Tags Config
// @Accept json
// @Produce json
// @Param key query string true "config key"
// @Success 200 {object} Config
// @Failure 400 {object} ErrorResponse
// @Router /configs/config/{key} [get]

func (c *configManager) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("key is required in query params"))
		return
	}
	cfg, err := c.getConfig(r.Context(), key)
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}
	response.WriteSuccessV2(r.Context(), w, cfg)
}

// Swagger doc
// @Summary Save config
// @Description Save config
// @Tags Config
// @Accept json
// @Produce json
// @Param config body Config true "config object"
// @Success 200 {string} string "config saved successfully"
// @Failure 400 {object} ErrorResponse
// @Router /configs/config [post]
func (c *configManager) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	var cfg Config
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		response.WriteErrorV2(r.Context(), w, customerrors.BAD_REQUEST_ERROR("failed to decode request body, Err: %s", err.Error()))
		return
	}
	err = c.saveConfig(r.Context(), &cfg)
	if err != nil {
		response.WriteErrorV2(r.Context(), w, err)
		return
	}
	response.WriteSuccessV2(r.Context(), w, "config saved successfully")
}

type configMetadata struct {
	Keys []string `json:"keys"`
}

// Swagger doc
// @Summary Get all config keys
// @Description Get all config keys
// @Tags Config
// @Accept json
// @Produce json
// @Success 200 {object} configMetadata
// @Failure 400 {object} ErrorResponse
// @Router /configs/metadata [get]
func (c *configManager) handleGetConfigMetadata(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(c.registeredConfigs))
	for k := range c.registeredConfigs {
		keys = append(keys, k)
	}
	response.WriteSuccessV2(r.Context(), w, configMetadata{Keys: keys})
}
