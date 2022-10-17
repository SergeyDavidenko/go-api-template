package repository

import (
	"{{cookiecutter.app_name}}/pkg/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	cfg *config.Config
	pg  *gorm.DB
}

func New(cfg *config.Config) (*DB, error) {
	dsn := cfg.BuildDSNPostgres()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(logrus.New(), logger.Config{
			Colorful: true,
		}),
	})
	if err != nil {
		return &DB{}, err
	}
	return &DB{
		cfg: cfg,
		pg:  db,
	}, nil
}

func (db *DB) Migrations(path string) error {
	err := db.pg.AutoMigrate()
	if err != nil {
		return err
	}
	logrus.Info("migrations done")
	return nil
}
