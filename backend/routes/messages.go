package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/services"
)

type MessageIncomingData struct {
	Content string `json:"content"`
	Image   string `json:"image"`
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
