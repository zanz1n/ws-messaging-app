package ws

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
	"github.com/zanz1n/ws-messaging-app/services"
)

var (
	validate = validator.New()
	ctx      = context.Background()
	pubCtx   = context.Background()
)

type MessageBroadcastedData struct {
	Type string   `json:"type" validate:"required"`
	Data *Message `json:"data" validate:"required"`
}

type MessageWebsocketEvent struct {
	Type string          `json:"type" validate:"required"`
	Data struct {
		Content *string `json:"content"`
		Image   *string `json:"image"`
	} `json:"data" validate:"required"`
}

type IncomingMessage struct {
	Content *string `json:"content"`
	Image   *string `json:"image"`
}

type Message struct {
	Content *string                  `json:"content"`
	Image   *string                  `json:"image"`
	Author  *services.UserJwtPayload `json:"author"`
}

type MessageEventHandler struct {
	db *dba.Queries
	ws *WebsocketService
}

func NewMessageEventHandler(db *dba.Queries, ws *WebsocketService) *MessageEventHandler {
	return &MessageEventHandler{
		db: db,
		ws: ws,
	}
}

func (m *MessageEventHandler) Handle(payload *IncomingMessage, user *services.UserJwtPayload) error {
	if err := validate.Struct(payload); err != nil {
		return err
	}

	log.Println(*payload)

	if payload.Content != nil {
		if *payload.Content == "" {
			payload.Content = nil
		}
	}

	if payload.Image != nil {
		if !strings.HasPrefix(*payload.Image, "https://") {
			return errors.New("image must start with https://")
		}
	}

	if payload.Content == nil && payload.Image == nil {
		return errors.New("if image is empty, message content is required")
	}

	go m.SendMessage(&Message{
		Content: payload.Content,
		Image:   payload.Image,
		Author:  user,
	})

	return nil
}

func (m *MessageEventHandler) SendMessage(payload *Message) error {
	createMsgArgs := &dba.CreateMessageParams{
		ID:        uuid.New().String(),
		UserId:    payload.Author.ID,
		UpdatedAt: time.Now(),
	}

	if payload.Content != nil {
		createMsgArgs.Content = sql.NullString{
			String: *payload.Content,
			Valid:  true,
		}
	} else {
		createMsgArgs.Content = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	if payload.Image != nil {
		createMsgArgs.ImageUrl = sql.NullString{
			String: *payload.Image,
			Valid:  true,
		}
	} else {
		createMsgArgs.ImageUrl = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	// if _, err := m.db.CreateMessage(ctx, *createMsgArgs); err != nil {
	// 	return errors.New("something went wrong")
	// }

	if err := m.ws.Pub.Publish(pubCtx, "chat", &MessageBroadcastedData{
		Type: "message",
		Data: payload,
	}); err != nil {
		return errors.New("something went wrong")
	}

	return nil
}
