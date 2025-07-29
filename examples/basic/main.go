package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oduortoni/websocket/ws"
)

type SessionValidator struct{}

func (n *SessionValidator) Validate(r *http.Request) (ws.SessionInfo, error) {
	return ws.SessionInfo{}, nil
}

type MessageHandler struct{}

func (n *MessageHandler) Handle(client *ws.Client, data []byte) error {
	return nil
}

type EnvelopePersister struct{}

func (n *EnvelopePersister) SaveEnvelope(e ws.Envelope) error {
	return nil
}

func (n *EnvelopePersister) ConfirmDelivery(envelopeID ws.Identity, clientID ws.Identity) error {
	return nil
}

func main() {
	sessValidator := &SessionValidator{}
	messageHandler := &MessageHandler{}
	envelopePersister := &EnvelopePersister{}
	ws := ws.NewWebSocketHandler(sessValidator, messageHandler, envelopePersister)
	http.Handle("/ws", ws)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>WebSocket Example</title>
		</head>
		<body>
			<h1>WebSocket Example</h1>
			<script>
				const socket = new WebSocket("ws://localhost:9000/ws");
				socket.addEventListener("open", (event) => {
					console.log("Connected to WebSocket server");
				});
				socket.addEventListener("message", (event) => {
					console.log("Received message from server:", event.data);
				});
				socket.addEventListener("close", (event) => {
					console.log("Disconnected from WebSocket server");
				});
				socket.addEventListener("error", (event) => {
					console.error("WebSocket error:", event);
				});
			</script>
		</body>
		</html>`
		w.Write([]byte(page))
	})

	fmt.Println("Server started on :9000")
	err := http.ListenAndServe(":9000", nil)
	log.Fatal(err)
}
