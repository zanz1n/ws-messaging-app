package services

import (
	"context"
	"log"
)

func SubscribeOnGlobalWs(w *WebsocketService) {
	subCtx := context.Background()
	sub := w.Sub.Subscribe(subCtx, "ws_global")

	go func() {
		for {
			msg, err := sub.ReceiveMessage(subCtx)

			if err != nil {
				log.Printf("Error receiving message: %s", err.Error())
				continue
			}

			go w.BroadcastLocal([]byte(msg.Payload))
		}
	}()
}
