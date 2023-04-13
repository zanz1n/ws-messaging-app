package services

import (
	"time"
)

type CreateMessageDto struct {
	Content string         `json:"content"`
	Image   string         `json:"image"`
	User    UserJwtPayload `json:"user"`
}

type MessageCreatePayload struct {
	Type    string         `json:"type"`
	ID      string         `json:"id"`
	User    UserJwtPayload `json:"user"`
	Image   *string        `json:"image"`
	Content *string        `json:"content"`
}

type UserReturnedOnMessage struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type MessageCreateReturnedData struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	User      UserReturnedOnMessage
	Image     *string `json:"image"`
	Content   *string `json:"content"`
}

type MessageDeletePayload struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
