package repository

import (
	"{{cookiecutter.app_name}}/pkg/config"
	"github.com/sirupsen/logrus"
	{% if cookiecutter.db_type == "postgres" %}
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	{% endif %}
	{% if cookiecutter.db_type == "mongodb" %}
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	{% endif %}
)

type DB struct {
	cfg *config.Config
	{% if cookiecutter.db_type == "postgres" %}
	pg  *gorm.DB
	{% endif %}
	{% if cookiecutter.db_type == "mongodb" %}
	client *mongo.Client
	{% endif %}
}

func New(cfg *config.Config) (*DB, error) {
	{% if cookiecutter.db_type == "postgres" %}
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
	{% endif %}
	{% if cookiecutter.db_type == "mongodb" %}
	client, err := mongo.Connect(context.Background(), cfg.BuildDSNMongoDB())
	if err != nil {
		logrus.Error(err)
		return &DB{}, err
	}
	return &DB{
		cfg:  cfg,
		client: client,
	}, nil
	{% endif %}
	{% if cookiecutter.db_type == "none" %}
	return &DB{
		cfg: cfg,
	}, nil
	{% endif %}
}

func (db *DB) Migrations(path string) error {
	{% if cookiecutter.db_type == "postgres" %}
	err := db.pg.AutoMigrate()
	if err != nil {
		return err
	}
	{% endif %}
	logrus.Info("migrations done")
	return nil
}
