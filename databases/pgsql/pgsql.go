package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofreego/goutils/logger"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
	SSLMode  string `yaml:"SSLMode"` //
}

func getPGSQLConn(ctx context.Context, cfg *Config) *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Panic(ctx, "failed to connect to postgresql , Err: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		logger.Panic(ctx, "failed to ping to postgresql , Err: %s", err.Error())
	}

	return db
}
