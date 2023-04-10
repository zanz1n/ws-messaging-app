package services

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
)

type BaseWebsocketEvent struct {
	Type string                 `json:"type" validate:"required"`
	Data map[string]interface{} `json:"data" validate:"required"`
}

type WebsocketService struct {
	conns  map[string]*websocket.Conn
	ConnsM *sync.Mutex
	Pub    *redis.Client
	Sub    *redis.Client
}

func NewWebsocketService(pubClient *redis.Client, subClient *redis.Client) *WebsocketService {
	ws := WebsocketService{
		conns:  make(map[string]*websocket.Conn),
		ConnsM: &sync.Mutex{},
		Pub:    pubClient,
		Sub:    subClient,
	}

	go SubscribeOnGlobalWs(&ws)

	return &ws
}

func (ws *WebsocketService) GetConns() map[string]*websocket.Conn {
	ws.ConnsM.Lock()
	defer ws.ConnsM.Unlock()
	return ws.conns
}

func (ws *WebsocketService) AddConn(conn *websocket.Conn) string {
	addr := conn.RemoteAddr().String()

	log.Printf("[%s]:%s  \x1b[35mWS\x1b[0m\t%s\t%s\x1b[0m",
		strings.Split(addr, ":")[0],
		strings.Split(addr, ":")[1],
		"/api/gateway",
		"\x1b[32mOPENED",
	)

	ws.ConnsM.Lock()

	defer ws.ConnsM.Unlock()

	ws.conns[conn.RemoteAddr().String()] = conn

	return addr
}

func (ws *WebsocketService) RemoveConn(addr string) {
	log.Printf("[%s]:%s  \x1b[35mWS\x1b[0m\t%s\t%s\x1b[0m",
		strings.Split(addr, ":")[0],
		strings.Split(addr, ":")[1],
		"/api/gateway",
		"\x1b[31mCLOSED",
	)

	ws.ConnsM.Lock()

	defer ws.ConnsM.Unlock()

	delete(ws.conns, addr)
}

func (ws *WebsocketService) BroadcastLocal(payload []byte) {
	ws.ConnsM.Lock()

	defer ws.ConnsM.Unlock()

	for _, conn := range ws.conns {
		err := conn.WriteMessage(1, payload)
		if err != nil {
			log.Printf("error writing message to connection: %s", err.Error())
		}
	}
}

func (ws *WebsocketService) BroadcastRemote(payload []byte) {
	err := ws.Pub.Publish(context.Background(), "ws_global", payload).Err()
	if err != nil {
		log.Printf("error publishing message to redis: %s", err.Error())
	}
}
