package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID        Identity
	Conn      *websocket.Conn
	Send      chan []byte
	Connected time.Time
}

func NewClient(id Identity, conn *websocket.Conn) *Client {
	return &Client{
		ID:        Identity(id),
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Connected: time.Now(),
	}
}
