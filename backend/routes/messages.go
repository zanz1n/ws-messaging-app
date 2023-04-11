package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/services"
)

func DeleteMessage(ms *services.MessagesService) func (c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		msgId := c.Params("id")

		if len(msgId) > 36 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "uuid params must be shorter than 36 characters",
			})
		}

		var user = c.Locals("user").(*services.UserJwtPayload)

		if allowed, _ := ms.IsAllowed(user.ID, msgId); !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "you can't edit a message you did't create",
			})
		}

		if err := ms.Delete(msgId); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "something went wrong",
			})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "message deleted",
			"id": msgId,
		})
	}
}
