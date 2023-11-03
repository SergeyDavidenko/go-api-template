package rest

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"{{cookiecutter.app_name}}/internal/repository"
)

func TestHandler_Version(t *testing.T) {
	handler := New(&repository.DB{})
	tt := fiber.New()
	tt.Get("/version", handler.Version)

	req := httptest.NewRequest("GET", "http://localhost:8080/version", nil)
	resp, err := tt.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	var version map[string]string
	err = json.NewDecoder(resp.Body).Decode(&version)
	assert.Nil(t, err)
	if val, ok := version["version"]; ok {
		assert.NotEqual(t, "", val)
	}
}
