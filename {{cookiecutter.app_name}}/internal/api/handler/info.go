package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/config"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/internal/models"
)

// Info about app
func Info(c *fiber.Ctx) error {
	info := models.Info{
		AppName: config.Conf.Core.Lable,
		Env: config.Conf.Core.Environment,
	}
	return c.JSON(fiber.Map{"info": info})
}