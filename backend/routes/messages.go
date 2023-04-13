package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/services"
)

type MessageIncomingData struct {
	Content string `json:"content"`
	Image   string `json:"image"`
}

func DeleteMessage(ms *services.MessagesService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		msgId := c.Params("id")

		if len(msgId) > 36 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "uuid params must be shorter than 36 characters",
			})
		}

		user := c.Locals("user").(*services.UserJwtPayload)

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
			"id":      msgId,
		})
	}
}

func PostMessage(ms *services.MessagesService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var (
			user = c.Locals("user").(*services.UserJwtPayload)
			data MessageIncomingData
			err  error
		)

		if err = c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err = validate.Struct(data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		pubMsg := &services.CreateMessageDto{
			User:    *user,
			Content: data.Content,
			Image:   data.Image,
		}

		msg, status, err := ms.Publish(pubMsg)

		if err != nil {
			return c.Status(status).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"data": msg,
		})
	}
}
