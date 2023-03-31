package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/services"
)

type SignInBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func PostSignIn(as *services.AuthService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var (
			body SignInBody
		)

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid from body",
			})
		}

		if err := validate.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid from body",
			})
		}

		jwtToken, err := as.AuthenticateUser(body.Username, body.Password)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "username or password do not match",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": jwtToken,
		})
	}
}
