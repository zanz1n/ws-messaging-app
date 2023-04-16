package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
	"github.com/zanz1n/ws-messaging-app/services"
	"github.com/zanz1n/ws-messaging-app/utils"
)

type UserCreatePayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func PostSignUp(as *services.AuthService, db *dba.Queries) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var (
			body UserCreatePayload
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

		err := db.CreateUser(ctx, dba.CreateUserParams{
			Username:  body.Username,
			Password:  as.GenerateHash(body.Password),
			UpdatedAt: time.Now(),
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
			"user": user,
			"token":   token,
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
