package repository

import (
	"{{cookiecutter.app_name}}/pkg/config"
	"github.com/sirupsen/logrus"
	{% if cookiecutter.db_type == "postgres" %}
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	{% endif %}
	{% if cookiecutter.db_type == "mongodb" %}
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	// Get underlying sql.DB from GORM
	sqlDB, err := db.pg.DB()
	if err != nil {
		logrus.Error("Failed to get underlying sql.DB")
		return err
	}
	
	// Create postgres driver instance
	driver, err := mpostgres.WithInstance(sqlDB, &mpostgres.Config{})
	if err != nil {
		logrus.Error("Failed to create postgres driver")
		return err
	}
	
	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"postgres", 
		driver,
	)
	if err != nil {
		logrus.Error("Failed to create migrate instance")
		return err
	}
	defer m.Close()
	
	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Error("Failed to run migrations")
		return err
	}
	
	logrus.Info("Database migrations completed successfully")
	{% endif %}
	logrus.Info("migrations done")
	return nil
}
