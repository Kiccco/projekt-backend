package events

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type GetMessagesEvent struct {
	Connection *websocket.Conn
	UUID       uuid.UUID
	UserID     int
}

type SendMessageEvent struct {
	Connection *websocket.Conn
	UUID       uuid.UUID
	Message    string
}

type LogoutEvent struct {
	UserID int
}
