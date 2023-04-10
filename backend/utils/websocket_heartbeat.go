package utils

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

type HeartbeatPayload struct {
	Type string `json:"type" validate:"required"`
}

type WebsocketHeartbeat struct {
	SInterval int
	lastPing  time.Time
}

func NewWebsocketHeartbeat(interval int) *WebsocketHeartbeat {
	return &WebsocketHeartbeat{
		SInterval: interval,
	}
}

func (h *WebsocketHeartbeat) Ping() {
	h.lastPing = time.Now()
}

func (h *WebsocketHeartbeat) Start(conn *websocket.Conn) {
	for {
		time.Sleep(time.Duration(h.SInterval) * time.Second)

		if time.Since(h.lastPing) > time.Duration(h.SInterval)*time.Second {
			break
		}
	}
}
