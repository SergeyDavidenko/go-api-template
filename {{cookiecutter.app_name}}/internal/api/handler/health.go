 
package handler

import (
	"github.com/gofiber/fiber/v2"
)

// Healtz check
func Healtz(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
