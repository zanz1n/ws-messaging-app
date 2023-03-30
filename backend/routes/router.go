package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/zanz1n/ws-messaging-app/middlewares"
	"github.com/zanz1n/ws-messaging-app/services"
)

func NewRouter(app *fiber.App) {
	if services.ConfigProvider().AppEnv == "development" {
		app.Static("/", "./frontend/dist", fiber.Static{
			ByteRange:     true,
			Index:         "index.html",
			Compress:      true,
			MaxAge:        16,
			CacheDuration: 16 * time.Second,
		})
	}

	db, conn := services.NewDbProvider()

	publisher := services.NewRedisProvider()
	subscriber := services.NewRedisProvider()

	_ = conn

	messagingService := services.NewMessagingService(publisher, subscriber, db)

	app.Use("/api/gateway", middlewares.NewWebsocketMiddleware())

	app.Get("/api/gateway", websocket.New(ChatGateway(messagingService)))
}
