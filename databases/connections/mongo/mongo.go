package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/gofreego/goutils/customerrors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Hosts                  string        `yaml:"Host"`
	Username               string        `yaml:"Username"`
	Password               string        `yaml:"Password"`
	Database               string        `yaml:"Database"`
	MaxPoolSize            uint64        `yaml:"MaxPoolSize"`
	MinPoolSize            uint64        `yaml:"MinPoolSize"`
	MaxConnIdleTime        time.Duration `yaml:"MaxConnIdleTime"`
	MaxConnecting          uint64        `yaml:"MaxConnecting"`
	ConnectTimeout         time.Duration `yaml:"ConnectTimeout"`
	ServerSelectionTimeout time.Duration `yaml:"ServerSelectionTimeout"`
	Direct                 bool          `yaml:"Direct"`
	ReplicaSet             string        `yaml:"ReplicaSet"`
}

// setDefaultPoolConfig sets default values for connection pool configuration
func (cfg *Config) withDefault() {
	if cfg.MaxPoolSize == 0 {
		cfg.MaxPoolSize = 100 // Default max pool size
	}
	if cfg.MinPoolSize == 0 {
		cfg.MinPoolSize = 5 // Default min pool size
	}
	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = 30 * time.Minute // Default max idle time
	}
	if cfg.MaxConnecting == 0 {
		cfg.MaxConnecting = 10 // Default max connecting
	}
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10 * time.Second // Default connect timeout
	}
	if cfg.ServerSelectionTimeout == 0 {
		cfg.ServerSelectionTimeout = 30 * time.Second // Default server selection timeout
	}
	if cfg.ReplicaSet == "" {
		cfg.ReplicaSet = "rs0" // Default replica set name
	}
}

func NewMongoConnection(ctx context.Context, cfg *Config) (*mongo.Client, error) {
	// Set default pool configuration if not provided
	cfg.withDefault()

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/?replicaSet=%s", cfg.Username, cfg.Password, cfg.Hosts, cfg.ReplicaSet))

	// Configure connection pool settings
	clientOptions.SetMaxPoolSize(cfg.MaxPoolSize)
	clientOptions.SetMinPoolSize(cfg.MinPoolSize)
	clientOptions.SetMaxConnIdleTime(cfg.MaxConnIdleTime)
	clientOptions.SetMaxConnecting(cfg.MaxConnecting)
	clientOptions.SetConnectTimeout(cfg.ConnectTimeout)
	clientOptions.SetServerSelectionTimeout(cfg.ServerSelectionTimeout)
	if cfg.Direct {
		clientOptions.SetDirect(true) // Use direct connection if specified
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_CONNECTION_FAILED, "failed to connect to MongoDB, Err: %s", err.Error())
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_PING_FAILED, "failed to ping MongoDB, Err: %s", err.Error())
	}

	return client, nil
}
