package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
)

type CreateMessageDto struct {
	Content string         `json:"content"`
	Image   string         `json:"image"`
	User    UserJwtPayload `json:"user"`
}

type MessageCreatePayload struct {
	Type    string  `json:"type"`
	ID      string  `json:"id"`
	Image   *string `json:"image"`
	Content *string `json:"content"`
}

type MessageDeletePayload struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type MessagesService struct {
	db *dba.Queries
	ws *WebsocketService
}

func NewMessagesService(db *dba.Queries, ws *WebsocketService) *MessagesService {
	return &MessagesService{
		db: db,
		ws: ws,
	}
}

func (s *MessagesService) Publish(data *CreateMessageDto) (*dba.Message, error) {

	message := dba.CreateMessageParams{
		ID:        uuid.New().String(),
		UserId:    data.User.ID,
		UpdatedAt: time.Now(),
	}

	broadcast := MessageCreatePayload{
		Type: "messageCreated",
		ID:   message.ID,
	}

	if data.Content != "" {
		*broadcast.Content = data.Content
		message.Content = sql.NullString{String: data.Content, Valid: true}
	} else {
		broadcast.Image = nil
		message.Content = sql.NullString{String: "", Valid: false}
	}

	if data.Image != "" {
		*broadcast.Image = data.Image
		message.ImageUrl = sql.NullString{String: data.Image, Valid: true}
	} else {
		broadcast.Image = nil
		message.ImageUrl = sql.NullString{String: "", Valid: false}
	}

	payload, err := json.Marshal(&broadcast)

	if err != nil {
		return nil, errors.New("failed to marshal message")
	}

	result, err := s.db.CreateMessage(context.Background(), message)

	if err != nil {
		return nil, errors.New("message creation failed, try again later")
	}

	s.ws.BroadcastRemote(payload)

	return &result, err
}

func (s *MessagesService) Delete(id string) error {
	if err := s.db.DeleteMessageById(context.Background(), id); err != nil {
		return err
	}

	broadcast := MessageDeletePayload{
		Type: "messageDeleted",
		ID:   id,
	}

	payload, err := json.Marshal(&broadcast)

	if err != nil {
		return err
	}

	s.ws.BroadcastRemote(payload)

	return err
}
