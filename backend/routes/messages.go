package routes

import (
	"log"
	"strconv"

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

func GetMessages(ms *services.MessagesService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tParam := c.Query("t")

		if tParam == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "'t' query param is required",
			})
		}

		tStamp, err := strconv.Atoi(tParam)

		if err != nil {
			log.Println(tStamp)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "'t' query param must be a valid interger",
			})
		}

		if tStamp < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "'t' query param must be greater than 0",
			})
		}

		limit := 128

		if lParam := c.Query("l"); lParam != "" {
			if lInt, err := strconv.Atoi(lParam); err == nil {
				if lInt > 256 || lInt < 0 {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "'t' query param must be less than 256 and greater than 0",
					})
				}
				limit = lInt
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "'t' query param must be a valid interger",
				})
			}
		}

		result, err := ms.GetUntilTimestamp(int64(tStamp), int32(limit))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data": result,
		})
	}
}
