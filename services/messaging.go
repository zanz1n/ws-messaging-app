package services

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/goccy/go-json"
	"github.com/gofiber/websocket/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
)

type ChatMessage struct {
	Content *string `json:"content"`
	Image   *string `json:"image"`
}

type MessagingService struct {
	db          *dba.Queries
	connections map[string]*websocket.Conn
	connMutex   *sync.Mutex
}

func NewMessagingService(ch *amqp.Channel, db *dba.Queries) *MessagingService {
	return &MessagingService{
		db:          db,
		connections: make(map[string]*websocket.Conn),
		connMutex:   &sync.Mutex{},
	}
}

func (s *MessagingService) GetConnections() map[string]*websocket.Conn {
	return s.connections
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

	go s.BroadcastMessage(&message)

	return &message, nil
}

func (s *MessagingService) AddConn(conn *websocket.Conn) string {
	addr := conn.RemoteAddr().String()

	log.Printf("[%s]  \x1b[35mWS\x1b[0m\t  %s\t%s\x1b[0m",
		strings.Split(addr, ":")[0],
		"/api/gateway",
		"\x1b[32mOPENED",
	)

	s.connMutex.Lock()

	defer s.connMutex.Unlock()

	s.connections[conn.RemoteAddr().String()] = conn

	return addr
}

func (s *MessagingService) RemoveConn(addr string) {
	log.Printf("[%s]  \x1b[35mWS\x1b[0m\t  %s\t%s\x1b[0m",
		strings.Split(addr, ":")[0],
		"/api/gateway",
		"\x1b[31mCLOSED",
	)

	s.connMutex.Lock()

	defer s.connMutex.Unlock()

	delete(s.connections, addr)
}

func (s *MessagingService) BroadcastMessage(message *ChatMessage) {
	s.connMutex.Lock()

	defer s.connMutex.Unlock()

	msgBytes, err := json.Marshal(&message)

	if err != nil {
		return
	}

	for _, conn := range s.connections {
		err = conn.WriteMessage(1, msgBytes)
		if err != nil {
			log.Printf("error writing message to connection: %s", err.Error())
		}
	}
}
