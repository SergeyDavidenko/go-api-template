package postgresql

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// DB ...
var DB *sqlx.DB


// New func implements the storage interface
func New() *Storage {
	return &Storage{}
}


// Init client storage.
func (s *Storage) Init() error {
	// First set up the pgx connection pool
	connConfig := pgx.ConnConfig{
		Host:     "localhost",
		Database: "test",
		Password: "secret",
		User:     "go",
		Port: 5432,
	}
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AfterConnect:   nil,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		return errors.Wrap(err, "Call to pgx.NewConnPool failed")
	}

	// Then set up sqlx and return the created DB reference
	nativeDB := stdlib.OpenDBFromPool(connPool)

	DB = sqlx.NewDb(nativeDB, "pgx")
	return nil
}

// ShowVersion postgersql
func (s *Storage) ShowVersion() string{
	db := DB
	var version string
	err := db.Select(&version, sqlShowPostgresqlVersion)
	if err != nil {
		log.Error(err)
		return ""
	}
	return version
}


// Close the storage connection
func (s *Storage) Close() error {
	err := DB.Close()
	return err
}