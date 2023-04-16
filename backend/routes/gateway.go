package routes

import (
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/zanz1n/ws-messaging-app/services"
	"github.com/zanz1n/ws-messaging-app/utils"
)

const (
	wsBodyMax  int  = 128
	hbInterval uint = 32
)

func ChatGateway(s *services.WebsocketService, ap *services.AuthService) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()

		authToken := c.Query("auth_token")

		if authToken == "" {
			c.WriteJSON(fiber.Map{
				"error": "'auth_token' query param is required",
			})
			return
		}

		if _, err := ap.ValidateJwtToken(authToken); err != nil {
			c.WriteJSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}

		connId := s.AddConn(c)

		defer s.RemoveConn(connId)

		if hostname, err := os.Hostname(); err == nil {
			c.WriteJSON(fiber.Map{
				"instanceId": hostname,
				"heartbeat":  hbInterval,
			})
		}

		keepAlive := make(chan bool)

		heartbeat := utils.NewWebsocketHeartbeat(hbInterval, keepAlive)

		c.SetCloseHandler(func(code int, text string) error {
			keepAlive <- true
			return nil
		})

		go heartbeat.Start(c)

		go func() {
			var (
				err     error
				payload []byte = make([]byte, wsBodyMax)
				data    utils.HeartbeatPayload
				size    int
			)

			for {
				if _, payload, err = c.ReadMessage(); err != nil {
					break
				}

				if size = len(payload); size == 0 || size > wsBodyMax {
					if size > wsBodyMax {
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
				} else {
					if err = c.WriteJSON(fiber.Map{
						"error": "only ping events are supported",
					}); err != nil {
						break
					}
				}
			}
		}()

		<-keepAlive
	}
}
