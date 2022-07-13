package rest

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Version(c *fiber.Ctx) error {
	version := "0.0.1"
	return c.JSON(fiber.Map{"version": version})
}
