package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/goccy/go-json"
	"github.com/gofiber/websocket/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
)

var (
	mqCtx = context.Background()
)

type ChatMessage struct {
	Content *string `json:"content"`
	Image   *string `json:"image"`
}

type MessagingService struct {
	db        *dba.Queries
	conns     map[string]*websocket.Conn
	connMutex *sync.Mutex
	ch        *amqp.Channel
	queue     *amqp.Queue
}

func NewMessagingService(ch *amqp.Channel, db *dba.Queries) *MessagingService {
	queue, err := ch.QueueDeclare("chat", false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	ms := MessagingService{
		db:        db,
		conns:     make(map[string]*websocket.Conn),
		connMutex: &sync.Mutex{},
		ch:        ch,
		queue:     &queue,
	}

	keepAlive := make(chan bool)

	go func() {
		msgs, _ := ms.ch.Consume(ms.queue.Name, "", true, false, false, false, nil)

		var bodyParsed ChatMessage

		go func() {
			for m := range msgs {

				json.Unmarshal(m.Body, &bodyParsed)

				log.Printf("[%s]:%s  \x1b[35mWS\x1b[0m\t%s\x1b[0m\t%v",
					"NONE",
					"none",
					"\x1b[31mMESSAGE",
					bodyParsed,
				)
			}
		}()
		<-keepAlive
	}()

	return &ms
}

func (s *MessagingService) GetConnections() map[string]*websocket.Conn {
	s.connMutex.Lock()
	defer s.connMutex.Unlock()
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

	s.connMutex.Lock()

	defer s.connMutex.Unlock()

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

	s.connMutex.Lock()

	defer s.connMutex.Unlock()

	delete(s.conns, addr)
}

func (s *MessagingService) BroadcastLocal(message *ChatMessage) {
	s.connMutex.Lock()

	defer s.connMutex.Unlock()

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

	s.ch.PublishWithContext(mqCtx, "", s.queue.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBytes,
		},
	)

	log.Println(s.queue)
}
