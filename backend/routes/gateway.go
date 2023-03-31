package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/zanz1n/ws-messaging-app/services"
)

func ChatGateway(s *services.MessagingService) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()
		var (
			rawMsg []byte
			err    error
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

			_, err = s.HanleIncomingMessage(&rawMsg)

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
