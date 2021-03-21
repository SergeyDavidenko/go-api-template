package router

import (
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/config"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/api/handler"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	app.Get("/info", handler.Info)
}

// SetupRoutesHealtz setup router
func SetupRoutesHealtz(app *fiber.App) {
	app.Get(config.Conf.API.HealthURI, handler.Healtz)
	app.Get(config.Conf.API.MetricURI, adaptor.HTTPHandler(promhttp.Handler()))
}