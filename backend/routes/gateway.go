package routes

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/zanz1n/ws-messaging-app/services"
	"github.com/zanz1n/ws-messaging-app/utils"
)

const interval uint = 32

func ChatGateway(s *services.WebsocketService) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		connId := s.AddConn(c)

		defer s.RemoveConn(connId)
		defer c.Close()

		if hostname, err := os.Hostname(); err == nil {
			c.WriteJSON(fiber.Map{
				"instanceId": hostname,
				"heartbeat":  fmt.Sprintf("%v", interval),
			})
		}

		keepAlive := make(chan bool)

		heartbeat := utils.NewWebsocketHeartbeat(interval, keepAlive)

		c.SetCloseHandler(func(code int, text string) error {
			keepAlive <- true
			return nil
		})

		go heartbeat.Start(c)

		go func() {
			var (
				err     error
				payload []byte = make([]byte, 1024)
				data    utils.HeartbeatPayload
				size    int
			)

			for {
				if _, payload, err = c.ReadMessage(); err != nil {
					break
				}

				if size = len(payload); size == 0 || size > 1024 {
					if size > 1024 {
						if err = c.WriteJSON(fiber.Map{
							"error": "payload too large",
						}); err != nil {
							break
						}
					} else if size == 0 {
						if err = c.WriteJSON(fiber.Map{
							"error": "no data was sent",
						}); err != nil {
							break
						}
					}
					continue
				}

				data = utils.HeartbeatPayload{}

				if err = json.Unmarshal(payload, &data); err != nil {
					if err = c.WriteJSON(fiber.Map{
						"error": err.Error(),
					}); err != nil {
						break
					}
					continue
				}

				if data.Type == "ping" {
					heartbeat.Ping()
					if err = c.WriteJSON(fiber.Map{
						"type": "pong",
					}); err != nil {
						break
					}
				}
			}
		}()

		<-keepAlive
	}
}
