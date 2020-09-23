package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/{{cookiecutter.github_username}}/{{cookiecutter.app_name}}/config"
)

// Info func
func Info(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"info": config.Core.Lable, "env": config.Core.Environment})
}