package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	{% if cookiecutter.db_type == "mongodb" %}
	"go.mongodb.org/mongo-driver/mongo/options"
	{% endif %}
)

type (
	Config struct {
		HTTP     map[string]HTTPConfig
		Redis    RedisConfig
		Postgres PostgresConfig
		MongoDB  MongoDBConfig
		Custom   map[string]string
		Auth     Auth
	}

	Auth struct {
		Username string `mapstructure:"username" yaml:"username"`
		Password string `mapstructure:"password" yaml:"password"`
	}

	HTTPConfig struct {
		HostString         string        `mapstructure:"hostString"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegabytes"`
	}

	RedisConfig struct {
		HostString    string `mapstructure:"hostString"`
		Login         string `mapstructure:"login"`
		Password      string `mapstructure:"password"`
		ConnectsCount int    `mapstructure:"connectionsCount"`
		DBNum         int    `mapstructure:"dbNum"`
	}

	MongoDBConfig struct {
		HostString string `mapstructure:"hostString"`
		Port       int    `mapstructure:"port"`
		Login      string `mapstructure:"login"`
		Password   string `mapstructure:"password"`
		DBName     string `mapstructure:"dbName"`
	}

	PostgresConfig struct {
		HostString    string `mapstructure:"hostString"`
		Port          int    `mapstructure:"port"`
		Login         string `mapstructure:"login"`
		Password      string `mapstructure:"password"`
		DBName        string `mapstructure:"dbName"`
		ConnectsCount int    `mapstructure:"connectionsCount"`
	}
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("config validation error: %s - %s", e.Field, e.Message)
}

// MakeServerList splits host string into a list of servers
func MakeServerList(hostString string) []string {
	return strings.Split(hostString, ",")
}

// GetMongoDB returns MongoDB configuration
func (c *Config) GetMongoDB() *MongoDBConfig {
	return &c.MongoDB
}

// GetPostgres returns PostgreSQL configuration
func (c *Config) GetPostgres() *PostgresConfig {
	return &c.Postgres
}

// GetCustom returns custom configuration value
func (c *Config) GetCustom(key string) string {
	if config, ok := c.Custom[key]; ok {
		return config
	}
	logrus.Warnf("Custom config key %v not found", key)
	return ""
}

// GetHTTP returns HTTP configuration for the specified key
func (c *Config) GetHTTP(key string) *HTTPConfig {
	if config, ok := c.HTTP[key]; ok {
		return &config
	}
	logrus.Warnf("HTTP config key %v not found", key)
	return nil
}

// GetRedis returns Redis configuration
func (c *Config) GetRedis() *RedisConfig {
	return &c.Redis
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	var errors []ValidationError

	// Validate HTTP configurations
	if c.HTTP == nil {
		errors = append(errors, ValidationError{Field: "http", Message: "HTTP configuration is required"})
	} else {
		if apiConfig := c.HTTP["api"]; apiConfig.HostString == "" {
			errors = append(errors, ValidationError{Field: "http.api.hostString", Message: "API host string is required"})
		}
	}

	// Validate database configuration based on type
	{% if cookiecutter.db_type == "postgres" %}
	if c.Postgres.HostString == "" {
		errors = append(errors, ValidationError{Field: "postgres.hostString", Message: "PostgreSQL host is required"})
	}
	if c.Postgres.DBName == "" {
		errors = append(errors, ValidationError{Field: "postgres.dbName", Message: "PostgreSQL database name is required"})
	}
	{% endif %}

	{% if cookiecutter.db_type == "mongodb" %}
	if c.MongoDB.HostString == "" {
		errors = append(errors, ValidationError{Field: "mongodb.hostString", Message: "MongoDB host is required"})
	}
	if c.MongoDB.DBName == "" {
		errors = append(errors, ValidationError{Field: "mongodb.dbName", Message: "MongoDB database name is required"})
	}
	{% endif %}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %v", errors)
	}

	return nil
}

{% if cookiecutter.db_type == "mongodb" %}
// BuildDSNMongoDB builds MongoDB connection string
func (c *Config) BuildDSNMongoDB() *options.ClientOptions {
	mg := c.GetMongoDB()
	credential := options.Credential{
		AuthSource: mg.DBName,
		Username:   mg.Login,
		Password:   mg.Password,
	}
	if mg.HostString == "" {
		logrus.Fatal("mongodb hostname not set")
	}
	dsn := fmt.Sprintf("mongodb://%s:%d", mg.HostString, mg.Port)
	clientOpts := options.Client().ApplyURI(dsn).SetAuth(credential)
	logrus.Info("MongoDB DSN built successfully")
	return clientOpts
}
{% endif %}

// BuildDSNPostgres builds PostgreSQL connection string
func (c *Config) BuildDSNPostgres() string {
	pg := c.GetPostgres()
	if pg.HostString == "" {
		logrus.Fatal("postgres hostname not set")
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		pg.HostString, pg.Login, pg.Password, pg.DBName, pg.Port,
	)
	return dsn
}

// New creates a new configuration instance
func New(path, serviceName string) (*Config, error) {
	if err := parseConfigFile(path, serviceName); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Load from cloud config if enabled
	if viper.GetBool("CLOUD_CONFIG") {
		if viper.GetString("CLOUD_URL") == "" {
			return nil, fmt.Errorf("CLOUD_URL not set when CLOUD_CONFIG is enabled")
		}
		if serviceName == "" {
			serviceName = viper.GetString("SERVICE_NAME")
		}
		loadConfiguration(viper.GetString("CLOUD_URL"), serviceName, "default")
	} else {
		logrus.Infof("Loading config from %s", path)
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}

// unmarshal unmarshals configuration from viper
func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return fmt.Errorf("failed to unmarshal HTTP config: %w", err)
	}
	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return fmt.Errorf("failed to unmarshal Redis config: %w", err)
	}
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return fmt.Errorf("failed to unmarshal PostgreSQL config: %w", err)
	}
	if err := viper.UnmarshalKey("custom", &cfg.Custom); err != nil {
		return fmt.Errorf("failed to unmarshal custom config: %w", err)
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return fmt.Errorf("failed to unmarshal auth config: %w", err)
	}
	if err := viper.UnmarshalKey("mongodb", &cfg.MongoDB); err != nil {
		return fmt.Errorf("failed to unmarshal MongoDB config: %w", err)
	}

	// Override with environment variables
	cfg.Auth.Username = viper.GetString("auth.username")
	cfg.Auth.Password = viper.GetString("auth.password")
	
	// Override PostgreSQL config with environment variables
	if viper.GetString("postgres.hostString") != "" {
		cfg.Postgres.HostString = viper.GetString("postgres.hostString")
	}
	if viper.GetInt("postgres.port") != 0 {
		cfg.Postgres.Port = viper.GetInt("postgres.port")
	}
	if viper.GetString("postgres.login") != "" {
		cfg.Postgres.Login = viper.GetString("postgres.login")
	}
	if viper.GetString("postgres.password") != "" {
		cfg.Postgres.Password = viper.GetString("postgres.password")
	}
	if viper.GetString("postgres.dbName") != "" {
		cfg.Postgres.DBName = viper.GetString("postgres.dbName")
	}

	return nil
}

// parseConfigFile parses the configuration file
func parseConfigFile(filepath, serviceName string) error {
	path := strings.Split(filepath, "/")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("go")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(strings.Join(path[:len(path)-1], "/")) // folder
	viper.SetConfigName(serviceName)                           // config file name
	
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	return nil
}
