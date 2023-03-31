package routes

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/zanz1n/ws-messaging-app/middlewares"
	"github.com/zanz1n/ws-messaging-app/services"
)

var validate = validator.New()

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

	jwtService := services.NewJwtService()

	authService := services.NewAuthService(db, jwtService)

	publisher := services.NewRedisProvider()
	subscriber := services.NewRedisProvider()

	_ = conn

	messagingService := services.NewMessagingService(publisher, subscriber, db)

	app.Use("/api/gateway", middlewares.NewWebsocketMiddleware())
	app.Use("/api/gateway", middlewares.NewAuthMiddleware(authService))

	app.Post("/api/auth/signin", PostSignIn(authService))
	app.Post("/api/auth/signup", PostSignUp(authService, db))
	app.Get("/api/gateway", websocket.New(ChatGateway(messagingService)))
}
