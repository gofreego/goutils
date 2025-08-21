package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofreego/goutils/customerrors"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
	SSLMode  string `yaml:"SSLMode"`

	// Connection Pool Settings
	MaxOpenConns    int           `yaml:"MaxOpenConns"`    // Maximum number of open connections
	MaxIdleConns    int           `yaml:"MaxIdleConns"`    // Maximum number of idle connections
	ConnMaxLifetime time.Duration `yaml:"ConnMaxLifetime"` // Maximum amount of time a connection may be reused
	ConnMaxIdleTime time.Duration `yaml:"ConnMaxIdleTime"` // Maximum amount of time a connection may be idle
}

func (c *Config) WithDefaults() {
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 20 // Default value
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10 // Default value
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 30 * time.Minute // Default value
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = 5 * time.Minute // Default value
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable" // Default SSL mode
	}
}

func GetConnection(ctx context.Context, cfg *Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_INVALID_CONFIGURATION, "configuration cannot be nil")
	}
	cfg.WithDefaults()
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_CONNECTION_FAILED, "failed to connect to postgresql, Err: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_PING_FAILED, "failed to ping to postgresql, Err: %s", err.Error())
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns) // Default value
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}
