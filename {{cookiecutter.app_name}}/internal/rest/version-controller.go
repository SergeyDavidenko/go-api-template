package rest

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Version(c *fiber.Ctx) error {
	version := h.conf.GetCustom("version")
	return c.JSON(fiber.Map{"version": version})
}
