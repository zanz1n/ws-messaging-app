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

	wsService := services.NewWebsocketService(publisher, subscriber)

	messagingService := services.NewMessagesService(db, wsService)

	/* Websocket Middlewares */
	app.Use("/api/gateway", middlewares.NewWebsocketMiddleware())
	/* End Websocket Middlewares */

	/* Websocket Routes */
	app.Get("/api/gateway", websocket.New(ChatGateway(wsService, authService)))
	/* End Websocket Routes */

	/* Auth Middlewares */
	app.Use("/api/auth/self", middlewares.NewAuthMiddleware(authService))
	/* End Auth Middlewares */

	/* Auth Routes */
	app.Post("/api/auth/signin", PostSignIn(authService))
	app.Post("/api/auth/signup", PostSignUp(authService, db))
	app.Get("/api/auth/self", GetSelf())
	/* End Auth Routes */

	/* Message Middlewares */
	app.Use("/api/message", middlewares.NewAuthMiddleware(authService))
	app.Use("/api/messages", middlewares.NewAuthMiddleware(authService))
	/* End Message Middlewares */

	/* Message Routes */
	app.Delete("/api/message/:id", DeleteMessage(messagingService))
	app.Post("/api/messages", PostMessage(messagingService))
	app.Get("api/messages", GetMessages(messagingService))
	/* End Message Routes */

	/* User Routes */
	app.Get("/api/user/:id", GetUserById(db))
	/* End User Routes */
}
