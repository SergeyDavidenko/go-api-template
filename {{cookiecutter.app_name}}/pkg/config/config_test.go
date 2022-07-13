package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	err := os.Setenv("GO_AUTH_USERNAME", "test")
	assert.Nil(t, err)
	serviceName := "stripo-security-service"
	cfg, err := New("../../configs/", serviceName)
	assert.Nil(t, err)
	assert.Equal(t, ":8080", cfg.HTTP["api"].HostString)
	assert.Equal(t, "secret", cfg.Auth.Password)
	test := os.Getenv("GO_AUTH_USERNAME")
	assert.Equal(t, "test", test)
	assert.Equal(t, "test", cfg.Auth.Username)
	assert.Equal(t, "host=localhost user=go password=go dbname=securitydb port=5432 sslmode=disable", cfg.BuildDSNPostgres())
	assert.Equal(t, "localhost:6379", cfg.GetRedis().HostString)
}
