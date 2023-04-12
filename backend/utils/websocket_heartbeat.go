package utils

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

type HeartbeatPayload struct {
	Type string `json:"type" validate:"required"`
}

type WebsocketHeartbeat struct {
	SInterval uint
	lastPing  time.Time
	keepAlive chan bool
}

func NewWebsocketHeartbeat(interval uint, keepAlive chan bool) *WebsocketHeartbeat {
	return &WebsocketHeartbeat{
		SInterval: interval,
		keepAlive: keepAlive,
	}
}

func (h *WebsocketHeartbeat) Ping() {
	h.lastPing = time.Now()
}

func (h *WebsocketHeartbeat) Start(conn *websocket.Conn) {
	for {
		time.Sleep(time.Duration(h.SInterval) * time.Second)

		if time.Since(h.lastPing) > time.Duration(h.SInterval)*time.Second {
			h.keepAlive <- true
			break
		}
	}
}
