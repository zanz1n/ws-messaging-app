package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func NewWebsocketMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		
		return c.Status(fiber.ErrUpgradeRequired.Code).JSON(fiber.Map{
			"error": "websocket upgrade required",
		})
	}
}
