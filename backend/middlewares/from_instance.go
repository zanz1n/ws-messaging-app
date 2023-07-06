package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func NewFromInstanceMiddleware() func(c *fiber.Ctx) error {
	hostname, err := os.Hostname()

	if err != nil {
		return func(c *fiber.Ctx) error {
			return nil
		}
	}

	return func(c *fiber.Ctx) error {
		c.Append("x-cluster-machine-id", hostname)

		return c.Next()
	}
}
