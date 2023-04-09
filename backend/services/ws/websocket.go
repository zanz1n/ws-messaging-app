package ws

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zanz1n/ws-messaging-app/services"
)

type BaseWebsocketEvent struct {
	Type string      `json:"type" validate:"required"`
	Data interface{} `json:"data" validate:"required"`
}

type WebsocketService struct {
	conns  map[string]*websocket.Conn
	ConnsM *sync.Mutex
	Pub    *redis.Client
	Sub    *redis.Client
	msgH   *MessageEventHandler
}

func NewWebsocketService(pubClient *redis.Client, subClient *redis.Client) *WebsocketService {
	ws := WebsocketService{
		conns:  make(map[string]*websocket.Conn),
		ConnsM: &sync.Mutex{},
		Pub:    pubClient,
		Sub:    subClient,
	}

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

func (ws *WebsocketService) HandleRawPayload(p *[]byte, user *services.UserJwtPayload) error {
	var parsed = BaseWebsocketEvent{}

	if err := json.Unmarshal(*p, &parsed); err != nil {
		return err
	}

	if parsed.Type == "message" {
		msgParsed := MessageWebsocketEvent{}
		if err := json.Unmarshal(*p, &msgParsed); err != nil {
			return err
		}

		if err := validator.New().Struct(msgParsed); err != nil {
			return err
		}

		log.Println(msgParsed.Data)
		return nil
		// return ws.msgH.Handle(&msgParsed.Data, user)
	}

	return errors.New("invalid payload type")
}

/*
Must be launched inside a goroutine and payload must be
a pointer to a json resolvable interface

	go w.Broadcast(&map[string]interface{}{
		"prop": "value",
	})
*/
func (ws *WebsocketService) Broadcast(payload interface{}) {
	ws.ConnsM.Lock()

	defer ws.ConnsM.Unlock()

	msgBytes, err := json.Marshal(payload)

	if err != nil {
		return
	}

	for _, conn := range ws.conns {
		err = conn.WriteMessage(1, msgBytes)
		if err != nil {
			log.Printf("error writing message to connection: %s", err.Error())
		}
	}
}
