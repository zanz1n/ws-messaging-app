package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/goccy/go-json"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
)

var (
	pubCtx = context.Background()
	subCtx = context.Background()
)

type ChatMessage struct {
	Content *string `json:"content"`
	Image   *string `json:"image"`
}

type MessagingService struct {
	db     *dba.Queries
	conns  map[string]*websocket.Conn
	connsM *sync.Mutex
	pub    *redis.Client
	sub    *redis.Client
}

func NewMessagingService(pubClient *redis.Client, subClient *redis.Client, db *dba.Queries) *MessagingService {
	ms := MessagingService{
		db:     db,
		conns:  make(map[string]*websocket.Conn),
		connsM: &sync.Mutex{},
		pub:    pubClient,
		sub:    subClient,
	}

	go func() {

		sub := ms.sub.Subscribe(subCtx, "chat")

		var (
			bodyParsed ChatMessage
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
				go ms.BroadcastLocal(&bodyParsed)
			}
		}()
	}()

	return &ms
}

func (s *MessagingService) GetConnections() map[string]*websocket.Conn {
	s.connsM.Lock()
	defer s.connsM.Unlock()
	return s.conns
}

func (s *MessagingService) HanleIncomingMessage(rawPayload *[]byte) (*ChatMessage, error) {
	var (
		message ChatMessage
		err     error
	)

	err = json.Unmarshal(*rawPayload, &message)

	if err != nil {
		return nil, err
	}

	if message.Content != nil {
		if *message.Content == "" {
			message.Content = nil
		}
	}

	if message.Image != nil {
		if !strings.HasPrefix(*message.Image, "https://") {
			return nil, errors.New("image must start with https://")
		}
	}

	if message.Content == nil && message.Image == nil {
		return nil, errors.New("if image is empty, message content is required")
	}

	go s.BroadcastGlobal(&message)

	return &message, nil
}

func (s *MessagingService) AddConn(conn *websocket.Conn) string {
	addr := conn.RemoteAddr().String()

	log.Printf("[%s]:%s  \x1b[35mWS\x1b[0m\t%s\t%s\x1b[0m",
		strings.Split(addr, ":")[0],
		strings.Split(addr, ":")[1],
		"/api/gateway",
		"\x1b[32mOPENED",
	)

	s.connsM.Lock()

	defer s.connsM.Unlock()

	s.conns[conn.RemoteAddr().String()] = conn

	return addr
}

func (s *MessagingService) RemoveConn(addr string) {
	log.Printf("[%s]:%s  \x1b[35mWS\x1b[0m\t%s\t%s\x1b[0m",
		strings.Split(addr, ":")[0],
		strings.Split(addr, ":")[1],
		"/api/gateway",
		"\x1b[31mCLOSED",
	)

	s.connsM.Lock()

	defer s.connsM.Unlock()

	delete(s.conns, addr)
}

func (s *MessagingService) BroadcastLocal(message *ChatMessage) {
	s.connsM.Lock()

	defer s.connsM.Unlock()

	msgBytes, err := json.Marshal(&message)

	if err != nil {
		return
	}

	for _, conn := range s.conns {
		err = conn.WriteMessage(1, msgBytes)
		if err != nil {
			log.Printf("error writing message to connection: %s", err.Error())
		}
	}
}

func (s *MessagingService) BroadcastGlobal(message *ChatMessage) {
	msgBytes, err := json.Marshal(&message)

	if err != nil {
		log.Printf("Error marshalling message: %s", err.Error())
		return
	}

	s.pub.Publish(pubCtx, "chat", msgBytes)
}
