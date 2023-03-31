package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
	"github.com/zanz1n/ws-messaging-app/services"
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
			ID:        uuid.New().String(),
		})

		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "username already taken",
			})
		}

		token, _ := as.AuthenticateUser(body.Username, body.Password)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "user created",
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
