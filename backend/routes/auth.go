package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/dba"
	"github.com/zanz1n/ws-messaging-app/services"
	"github.com/zanz1n/ws-messaging-app/utils"
)

type SignInBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpBody struct {
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

func PostSignUp(as *services.AuthService, db *dba.Queries) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var (
			body SignUpBody
			ctx  = context.Background()
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

		now := time.Now().UnixMilli()

		err := db.CreateUser(ctx, dba.CreateUserParams{
			Username:  body.Username,
			Password:  as.GenerateHash(body.Password),
			CreatedAt: now,
			UpdatedAt: now,
			Role:      "USER",
			ID:        utils.RandomId(),
		})

		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "username already taken",
			})
		}

		token, _ := as.AuthenticateUser(body.Username, body.Password)

		user, err := as.ValidateJwtToken(token)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "something went wrong",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"user":  user,
			"token": token,
		})
	}
}
