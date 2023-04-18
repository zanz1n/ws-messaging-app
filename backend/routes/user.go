package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/services"
)

/*
 * GetSelf returns the user's information
 * Needs auth middleware
 */
func GetSelf() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*services.UserJwtPayload)
		return c.JSON(user)
	}
}
