package main

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/zanz1n/ws-messaging-app/routes"
	"github.com/zanz1n/ws-messaging-app/services"
)

func init() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(err)
		}
	}

	log.SetPrefix(fmt.Sprintf("%v - ", os.Getpid()))

	services.GenerateConfigsFromEnv()
}

func main() {
	config := services.ConfigProvider()

	app := fiber.New(fiber.Config{
		Prefork:       config.AppFork,
		CaseSensitive: true,
		StrictRouting: false,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ServerHeader:  "Fiber",
		AppName:       "Ws Messaging App",
	})

	app.Use(logger.New(logger.Config{
		Format:     "${pid} - ${time} [${ip}]:${port} ${method} ${path} ${status} ${latency}\n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	}))

	app.Use(recover.New())

	routes.NewRouter(app)

	if config.UseTls {
		log.Fatal(app.ListenTLS(config.BindAddr, config.TlsCertPath, config.TlsKeyPath))
	} else {
		log.Fatal(app.Listen(config.BindAddr))
	}
}
