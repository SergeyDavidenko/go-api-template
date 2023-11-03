package rest

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "OK"})
}
