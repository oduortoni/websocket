package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	SessionValidator  SessionValidator
	MessageHandler    MessageHandler
	EnvelopePersister EnvelopePersister
}

func NewWebSocketHandler(validator SessionValidator, messeger MessageHandler, persister EnvelopePersister) *WebsocketHandler {
	return &WebsocketHandler{
		SessionValidator:  validator,
		MessageHandler:    messeger,
		EnvelopePersister: persister,
	}
}

func (h *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.SessionValidator.Validate(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	client := NewClient(session.ClientID, conn)
	defer client.Conn.Close()
	go HandleClient(client, h.MessageHandler, h.EnvelopePersister)
}

func HandleClient(client *Client, messager MessageHandler, persister EnvelopePersister) {
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		err = messager.Handle(client, message)
		if err != nil {
			continue
		}
	}

	client.Conn.Close()
}
