package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/config"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/model"
)

// Info about app
func Info(c *fiber.Ctx) error {
	info := model.Info{
		AppName: config.Conf.Core.Lable,
		Env: config.Conf.Core.Environment,
	}
	return c.JSON(fiber.Map{info})
}