package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/zanz1n/ws-messaging-app/routes"
)

var (
	useTls   bool
	tlsCert  string
	tlsKey   string
	bindAddr string
)

func main() {
	SetupEnv()

	app := fiber.New(fiber.Config{
		Prefork:       os.Getenv("APP_ENV") == "development",
		CaseSensitive: true,
		StrictRouting: false,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ServerHeader:  "Fiber",
		AppName:       "Ws Messaging App",
	})

	app.Use(logger.New(logger.Config{
		Format:     "${time} [${ip}] ${method} ${path} ${status} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	}))

	app.Use(recover.New())

	routes.NewRouter(app)

	if useTls {
		app.ListenTLS(bindAddr, tlsCert, tlsKey)
	} else {
		app.Listen(bindAddr)
	}
}

func SetupEnv() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(err)
		}
	}

	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "development")
	}

	if os.Getenv("BIND_ADDR") == "" {
		bindAddr = ":3333"
	} else {
		bindAddr = os.Getenv("BIND_ADDR")
	}

	if os.Getenv("DATABASE_URI") == "" {
		panic("DB_URI is not set")
	}

	if os.Getenv("TLS_CERT") != "" && os.Getenv("TLS_KEY") != "" {
		useTls = true
		tlsCert = os.Getenv("TLS_CERT")
		tlsKey = os.Getenv("TLS_KEY")
	}
}
