package routes

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/zanz1n/ws-messaging-app/services"
)

func NewRouter(app *fiber.App) {
	if os.Getenv("APP_ENV") == "development" {
		app.Static("/", "./frontend/dist", fiber.Static{
			ByteRange:     true,
			Index:         "index.html",
			Compress:      true,
			MaxAge:        16,
			CacheDuration: 16 * time.Second,
		})
	}

	dbctx := context.Background()
	db, conn := services.NewDbProvider()

	publisher := services.NewRedisProvider()
	subscriber := services.NewRedisProvider()

	_ = dbctx
	_ = conn

	messagingService := services.NewMessagingService(publisher, subscriber, db)

	app.Use("/api/gateway", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return c.Status(fiber.ErrUpgradeRequired.Code).JSON(fiber.Map{
			"error": "Websocket upgrade required",
		})
	})
	
	app.Get("/api/gateway", websocket.New(ChatGateway(messagingService)))
}
