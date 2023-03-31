package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/services"
)

func NewAuthMiddleware(ap *services.AuthService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing Authorization header",
			})
		}

		jwtPayload, err := ap.ValidateJwtToken(authHeader)

		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid jwt token",
			})
		}

		c.Locals("user", jwtPayload)

		return c.Next()
	}
}
