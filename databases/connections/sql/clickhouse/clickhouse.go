package clickhouse

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2" // registers "clickhouse" database/sql driver
	"github.com/gofreego/goutils/customerrors"
)

type Config struct {
	Host         string `yaml:"Host"`
	Port         int    `yaml:"Port"`
	Username     string `yaml:"Username"`
	Password     string `yaml:"Password"`
	Database     string `yaml:"Database"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
}

func (c *Config) WithDefaults() {
	if c.Port == 0 {
		c.Port = 9000
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 10
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 5
	}
}

// GetConnection opens a ClickHouse database/sql connection.
// Uses the v1 TCP DSN format required by golang-migrate's clickhouse driver.
func GetConnection(cfg *Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_INVALID_CONFIGURATION, "configuration cannot be nil")
	}
	cfg.WithDefaults()

	dsn := fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_CONNECTION_FAILED, "failed to open clickhouse connection, Err: %s", err.Error())
	}
	if err := db.Ping(); err != nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_PING_FAILED, "failed to ping clickhouse, Err: %s", err.Error())
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return db, nil
}
