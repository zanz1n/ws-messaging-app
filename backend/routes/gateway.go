package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/zanz1n/ws-messaging-app/services"
	"github.com/zanz1n/ws-messaging-app/services/ws"
)

func ChatGateway(s *ws.WebsocketService) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()
		var (
			rawMsg []byte
			err    error
			user = c.Locals("user").(*services.UserJwtPayload)
		)

		connId := s.AddConn(c)
		defer s.RemoveConn(connId)

		for {
			_, rawMsg, err = c.ReadMessage()

			if err != nil {
				err = c.WriteJSON(fiber.Map{
					"error": err.Error(),
				})
				if err != nil {
					break
				}
				continue
			}

			err = s.HandleRawPayload(&rawMsg, user)

			if err != nil {
				err = c.WriteJSON(fiber.Map{
					"error": err.Error(),
				})
				if err != nil {
					break
				}
				continue
			}
		}
	}
}
