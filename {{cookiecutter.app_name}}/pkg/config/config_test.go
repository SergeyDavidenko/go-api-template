package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				HTTP: map[string]HTTPConfig{
					"api": {
						HostString: ":8080",
					},
				},
				{% if cookiecutter.db_type == "postgres" %}
				Postgres: PostgresConfig{
					HostString: "localhost",
					DBName:     "testdb",
				},
				{% endif %}
				{% if cookiecutter.db_type == "mongodb" %}
				MongoDB: MongoDBConfig{
					HostString: "localhost",
					DBName:     "testdb",
				},
				{% endif %}
			},
			wantErr: false,
		},
		{
			name: "missing HTTP config",
			config: &Config{
				HTTP: nil,
			},
			wantErr: true,
		},
		{
			name: "missing API host string",
			config: &Config{
				HTTP: map[string]HTTPConfig{
					"api": {
						HostString: "",
					},
				},
			},
			wantErr: true,
		},
		{% if cookiecutter.db_type == "postgres" %}
		{
			name: "missing PostgreSQL host",
			config: &Config{
				HTTP: map[string]HTTPConfig{
					"api": {
						HostString: ":8080",
					},
				},
				Postgres: PostgresConfig{
					HostString: "",
					DBName:     "testdb",
				},
			},
			wantErr: true,
		},
		{
			name: "missing PostgreSQL database name",
			config: &Config{
				HTTP: map[string]HTTPConfig{
					"api": {
						HostString: ":8080",
					},
				},
				Postgres: PostgresConfig{
					HostString: "localhost",
					DBName:     "",
				},
			},
			wantErr: true,
		},
		{% endif %}
		{% if cookiecutter.db_type == "mongodb" %}
		{
			name: "missing MongoDB host",
			config: &Config{
				HTTP: map[string]HTTPConfig{
					"api": {
						HostString: ":8080",
					},
				},
				MongoDB: MongoDBConfig{
					HostString: "",
					DBName:     "testdb",
				},
			},
			wantErr: true,
		},
		{
			name: "missing MongoDB database name",
			config: &Config{
				HTTP: map[string]HTTPConfig{
					"api": {
						HostString: ":8080",
					},
				},
				MongoDB: MongoDBConfig{
					HostString: "localhost",
					DBName:     "",
				},
			},
			wantErr: true,
		},
		{% endif %}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_GetCustom(t *testing.T) {
	config := &Config{
		Custom: map[string]string{
			"version": "1.0.0",
			"env":     "production",
		},
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "existing key",
			key:      "version",
			expected: "1.0.0",
		},
		{
			name:     "non-existing key",
			key:      "nonexistent",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GetCustom(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_GetHTTP(t *testing.T) {
	config := &Config{
		HTTP: map[string]HTTPConfig{
			"api": {
				HostString: ":8080",
			},
			"health": {
				HostString: ":1499",
			},
		},
	}

	tests := []struct {
		name     string
		key      string
		expected *HTTPConfig
	}{
		{
			name: "existing key",
			key:  "api",
			expected: &HTTPConfig{
				HostString: ":8080",
			},
		},
		{
			name:     "non-existing key",
			key:      "nonexistent",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GetHTTP(tt.key)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected.HostString, result.HostString)
			}
		})
	}
}

func TestConfig_Integration(t *testing.T) {
	// Set up environment variables
	err := os.Setenv("GO_AUTH_USERNAME", "testuser")
	require.NoError(t, err)
	defer os.Unsetenv("GO_AUTH_USERNAME")

	// Test configuration loading
	serviceName := "{{cookiecutter.app_name}}"
	cfg, err := New("../../configs/", serviceName)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Test HTTP configuration
	apiConfig := cfg.GetHTTP("api")
	assert.NotNil(t, apiConfig)
	assert.Equal(t, ":8080", apiConfig.HostString)

	// Test custom configuration
	version := cfg.GetCustom("version")
	assert.NotEmpty(t, version)

	// Test environment variable override
	assert.Equal(t, "testuser", cfg.Auth.Username)

	{% if cookiecutter.db_type == "postgres" %}
	// Test PostgreSQL configuration
	pgConfig := cfg.GetPostgres()
	assert.NotNil(t, pgConfig)
	assert.Equal(t, "localhost", pgConfig.HostString)
	assert.Equal(t, "{{cookiecutter.app_name}}", pgConfig.DBName)

	// Test DSN building
	dsn := cfg.BuildDSNPostgres()
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "dbname={{cookiecutter.app_name}}")
	{% endif %}

	{% if cookiecutter.db_type == "mongodb" %}
	// Test MongoDB configuration
	mongoConfig := cfg.GetMongoDB()
	assert.NotNil(t, mongoConfig)
	assert.Equal(t, "localhost", mongoConfig.HostString)
	assert.Equal(t, "{{cookiecutter.app_name}}", mongoConfig.DBName)

	// Test DSN building
	clientOpts := cfg.BuildDSNMongoDB()
	assert.NotNil(t, clientOpts)
	{% endif %}

	// Test Redis configuration
	redisConfig := cfg.GetRedis()
	assert.NotNil(t, redisConfig)
	assert.Equal(t, "localhost:6379", redisConfig.HostString)
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "test.field",
		Message: "test message",
	}

	expected := "config validation error: test.field - test message"
	assert.Equal(t, expected, err.Error())
}

func TestMakeServerList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single server",
			input:    "localhost:8080",
			expected: []string{"localhost:8080"},
		},
		{
			name:     "multiple servers",
			input:    "server1:8080,server2:8080,server3:8080",
			expected: []string{"server1:8080", "server2:8080", "server3:8080"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MakeServerList(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
