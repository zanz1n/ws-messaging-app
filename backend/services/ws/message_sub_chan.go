package ws

import (
	"context"
	"log"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

func SubscribeForMessages(w *WebsocketService) {
	subCtx := context.Background()
	sub := w.Sub.Subscribe(subCtx, "chat")

	var (
		bodyParsed BaseWebsocketEvent
		err        error
		msg        *redis.Message
	)

	go func() {
		for {
			msg, err = sub.ReceiveMessage(subCtx)

			if err != nil {
				log.Printf("Error receiving message: %s", err.Error())
				continue
			}

			err = json.Unmarshal([]byte(msg.Payload), &bodyParsed)

			if err != nil {
				log.Printf("Error unmarshalling message: %s", err.Error())
				continue
			}
			go w.Broadcast(&bodyParsed)
		}
	}()
}
