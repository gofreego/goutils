package sql

import (
	"database/sql"
	"math/rand"

	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/databases/pgsql"
)

type Name string

const (
	MySQL    Name = "mysql"
	Postgres Name = "postgres"
)

type PostgresqlConfig struct {
	Primary pgsql.Config
	Replica []pgsql.Config
}

type Config struct {
	Name       Name `yaml:"Name"`
	Postgresql PostgresqlConfig
}

type DBManager interface {
	// Primary returns the primary database connection.
	Primary() *sql.DB
	// Replica returns the random replica database connection.
	// eventually distribution will be equal
	Replica() *sql.DB
}

type DBManagerImpl struct {
	primary      *sql.DB
	replica      []*sql.DB
	noOfReplicas int
}

func (d *DBManagerImpl) Primary() *sql.DB {
	return d.primary
}
func (d *DBManagerImpl) Replica() *sql.DB {
	if d.noOfReplicas == 0 {
		return d.primary
	}
	if d.noOfReplicas == 1 {
		return d.replica[0]
	}
	randomIndex := rand.Intn(d.noOfReplicas)
	return d.replica[randomIndex]
}

// NewDBManager creates a new DBManager based on the provided configuration.
func NewDBManager(cfg *Config) (DBManager, error) {
	if cfg.Name != Postgres {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_INVALID_CONFIGURATION, "unsupported database type: %s", cfg.Name)
	}
	primaryDB, err := pgsql.GetConnection(&cfg.Postgresql.Primary)
	if err != nil {
		return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_CONNECTION_FAILED, "failed to connect to primary database, Err: %s", err.Error())
	}

	// Create connections for each replica database
	replicaDBs := make([]*sql.DB, len(cfg.Postgresql.Replica))
	for i, replicaCfg := range cfg.Postgresql.Replica {
		replicaDBs[i], err = pgsql.GetConnection(&replicaCfg)
		if err != nil {
			return nil, customerrors.New(customerrors.ERROR_CODE_DATABASE_CONNECTION_FAILED, "failed to connect to replica database %d, Err: %s", i, err.Error())
		}
	}

	return &DBManagerImpl{
		primary:      primaryDB,
		replica:      replicaDBs,
		noOfReplicas: len(replicaDBs),
	}, nil
}
