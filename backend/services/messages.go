package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/zanz1n/ws-messaging-app/internal/dba"
	"github.com/zanz1n/ws-messaging-app/utils"
)

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

func (s *MessagesService) Publish(data *CreateMessageDto) (*MessageCreateReturnedData, int, error) {
	if data.Content == "" && data.Image == "" {
		return nil, 400, errors.New("content and image can't be empty at the same time")
	}

	message := dba.CreateMessageParams{
		ID:        utils.RandomId(),
		UserId:    data.User.ID,
		UpdatedAt: time.Now(),
	}

	broadcast := MessageCreatePayload{
		Type:    "messageCreated",
		ID:      message.ID,
		User:    data.User,
		Content: &data.Content,
		Image:   &data.Image,
	}

	if data.Content == "" {
		broadcast.Content = nil
		message.Content = sql.NullString{String: "", Valid: false}
	} else {
		message.Content = sql.NullString{String: data.Content, Valid: true}
	}

	if data.Image == "" {
		broadcast.Image = nil
		message.ImageUrl = sql.NullString{String: "", Valid: false}
	} else {
		message.ImageUrl = sql.NullString{String: data.Image, Valid: true}
	}

	payload, err := json.Marshal(&broadcast)

	if err != nil {
		return nil, 500, errors.New("failed to marshal message")
	}

	result, err := s.db.CreateMessage(context.Background(), message)

	if err != nil {
		log.Printf("error while creating message: %s", err.Error())
		return nil, 500, errors.New("message creation failed, try again later")
	}

	s.ws.BroadcastRemote(payload)

	return &MessageCreateReturnedData{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		User: UserReturnedOnMessage{
			ID:       data.User.ID,
			Username: data.User.Username,
		},
		Image:   broadcast.Image,
		Content: broadcast.Content,
	}, 200, err
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

func (s *MessagesService) IsAllowed(userId string, msgId string) (bool, error) {
	msg, err := s.db.GetMessageById(context.Background(), msgId)

	if err != nil {
		return false, err
	}

	if msg.UserId == userId {
		return true, nil
	}

	return false, nil
}
