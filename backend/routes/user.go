package routes

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/dba"
	"github.com/zanz1n/ws-messaging-app/services"
	"github.com/zanz1n/ws-messaging-app/utils"
)

func GetUserById(db *dba.Queries) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if len(id) != utils.RandomIdLen {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("id length must be %b", utils.RandomIdLen),
			})
		}

		var (
			user dba.GetUserByIdSafeRow
			err error
		)

		if user, err = db.GetUserByIdSafe(context.Background(), id); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("could not find user with id %s", id),
			})
		}

		return c.JSON(fiber.Map{
			"data": user,
		})
	}
}

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
