# WebSocket Library

A Go library for building WebSocket applications with session validation, message handling, and persistence capabilities.

## Features

- Session-based authentication and validation
- Structured message handling with custom processors
- Message persistence and delivery confirmation
- Client connection management with unique identities
- Built on top of the reliable Gorilla WebSocket library
- Comprehensive test coverage (80.8%)

## Installation

```bash
go get github.com/oduortoni/websocket
```

## Quick Start

### Basic Usage

```go
package main

import (
    "net/http"
    "github.com/oduortoni/websocket/ws"
)

func main() {
    // Create your implementations
    validator := &MySessionValidator{}
    messageHandler := &MyMessageHandler{}
    persister := &MyEnvelopePersister{}
    
    // Create WebSocket handler
    wsHandler := ws.NewWebSocketHandler(validator, messageHandler, persister)
    
    // Register the handler
    http.Handle("/ws", &wsHandler)
    
    // Start server
    http.ListenAndServe(":8080", nil)
}
```

## Core Concepts

### 1. Session Validation

Implement the `SessionValidator` interface to control who can connect:

```go
type SessionValidator interface {
    Validate(r *http.Request) (SessionInfo, error)
}
```

Example implementation:

```go
type MySessionValidator struct{}

func (v *MySessionValidator) Validate(r *http.Request) (ws.SessionInfo, error) {
    // Extract token from header, cookie, or query parameter
    token := r.Header.Get("Authorization")
    
    if token == "" {
        return ws.SessionInfo{}, errors.New("missing authorization token")
    }
    
    // Validate token and get user info
    userID, err := validateToken(token)
    if err != nil {
        return ws.SessionInfo{}, err
    }
    
    return ws.SessionInfo{
        ClientID: ws.Identity(userID),
        Metadata: map[string]string{
            "user_id": userID.String(),
            "role":    "user",
        },
    }, nil
}
```

### 2. Message Handling

Implement the `MessageHandler` interface to process incoming messages:

```go
type MessageHandler interface {
    Handle(client *Client, data []byte) error
}
```

Example implementation:

```go
type MyMessageHandler struct{}

func (h *MyMessageHandler) Handle(client *ws.Client, data []byte) error {
    // Parse the message
    var msg struct {
        Type    string      `json:"type"`
        Payload interface{} `json:"payload"`
    }
    
    if err := json.Unmarshal(data, &msg); err != nil {
        return err
    }
    
    // Handle different message types
    switch msg.Type {
    case "chat":
        return h.handleChatMessage(client, msg.Payload)
    case "ping":
        return h.handlePing(client)
    default:
        return fmt.Errorf("unknown message type: %s", msg.Type)
    }
}

func (h *MyMessageHandler) handleChatMessage(client *ws.Client, payload interface{}) error {
    // Process chat message
    // Broadcast to other clients, save to database, etc.
    return nil
}
```

### 3. Message Persistence

Implement the `EnvelopePersister` interface for message persistence:

```go
type EnvelopePersister interface {
    SaveEnvelope(e Envelope) error
    ConfirmDelivery(envelopeID Identity, clientID Identity) error
}
```

Example implementation:

```go
type MyEnvelopePersister struct {
    db *sql.DB
}

func (p *MyEnvelopePersister) SaveEnvelope(e ws.Envelope) error {
    query := `
        INSERT INTO envelopes (id, client_id, type, payload, timestamp)
        VALUES ($1, $2, $3, $4, $5)
    `
    
    payloadJSON, _ := json.Marshal(e.Payload)
    _, err := p.db.Exec(query, e.ID, e.ClientID, e.Type, payloadJSON, e.Timestamp)
    return err
}

func (p *MyEnvelopePersister) ConfirmDelivery(envelopeID, clientID ws.Identity) error {
    query := `
        UPDATE envelopes 
        SET delivered = NOW() 
        WHERE id = $1 AND client_id = $2
    `
    
    _, err := p.db.Exec(query, envelopeID, clientID)
    return err
}
```

## Advanced Usage

### Custom Client Management

Access client information in your message handler:

```go
func (h *MyMessageHandler) Handle(client *ws.Client, data []byte) error {
    // Access client properties
    clientID := client.ID
    connection := client.Conn
    sendChannel := client.Send
    connectedAt := client.Connected
    
    // Send message to specific client
    select {
    case client.Send <- []byte("Hello, client!"):
        // Message queued successfully
    default:
        // Client's send channel is full or closed
        return errors.New("failed to send message to client")
    }
    
    return nil
}
```

### Broadcasting Messages

Create a client manager to broadcast messages:

```go
type ClientManager struct {
    clients map[ws.Identity]*ws.Client
    mutex   sync.RWMutex
}

func (cm *ClientManager) AddClient(client *ws.Client) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    cm.clients[client.ID] = client
}

func (cm *ClientManager) RemoveClient(clientID ws.Identity) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    delete(cm.clients, clientID)
}

func (cm *ClientManager) Broadcast(message []byte) {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    for _, client := range cm.clients {
        select {
        case client.Send <- message:
            // Message sent successfully
        default:
            // Client disconnected, remove it
            go cm.RemoveClient(client.ID)
        }
    }
}
```

### Error Handling

The library provides several error scenarios you should handle:

1. **Session Validation Errors**: Return HTTP 401 Unauthorized
2. **WebSocket Upgrade Errors**: Return HTTP 400 Bad Request
3. **Message Handling Errors**: Continue processing other messages
4. **Connection Errors**: Automatically close and clean up

### Testing

Run the included tests:

```bash
# Run all tests
go test ./ws -v

# Run with coverage
go test ./ws -v -cover

# Run benchmarks
go test ./ws -v -bench=.
```

## Examples

### Chat Application

See `examples/basic/main.go` for a complete chat application example.

### Real-time Notifications

```go
type NotificationHandler struct {
    userSessions map[string]*ws.Client
}

func (h *NotificationHandler) Handle(client *ws.Client, data []byte) error {
    // Handle subscription requests
    var req struct {
        Action string `json:"action"`
        Topic  string `json:"topic"`
    }
    
    json.Unmarshal(data, &req)
    
    if req.Action == "subscribe" {
        h.subscribeToTopic(client, req.Topic)
    }
    
    return nil
}

func (h *NotificationHandler) SendNotification(topic string, message interface{}) {
    // Send to all subscribers of the topic
    for _, client := range h.getSubscribers(topic) {
        notification, _ := json.Marshal(map[string]interface{}{
            "topic":   topic,
            "message": message,
        })
        
        select {
        case client.Send <- notification:
        default:
            // Handle disconnected client
        }
    }
}
```

## API Reference

### Types

- `Identity`: UUID-based unique identifier
- `Client`: Represents a connected WebSocket client
- `SessionInfo`: Contains client ID and metadata
- `Envelope`: Message wrapper with persistence information
- `WebsocketHandler`: Main HTTP handler for WebSocket connections

### Interfaces

- `SessionValidator`: Validates incoming connection requests
- `MessageHandler`: Processes incoming messages from clients
- `EnvelopePersister`: Handles message persistence and delivery confirmation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Configuration

### WebSocket Upgrader Settings

The library uses sensible defaults, but you can customize the WebSocket upgrader by modifying the handler:

```go
// The library uses these default settings:
upgrader := websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
```

### Production Considerations

1. **CORS Configuration**: Add CheckOrigin for cross-origin requests:

```go
upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        // Implement your CORS policy
        origin := r.Header.Get("Origin")
        return isAllowedOrigin(origin)
    },
}
```

2. **Connection Limits**: Implement connection limiting in your SessionValidator
3. **Rate Limiting**: Add rate limiting to prevent message spam
4. **Graceful Shutdown**: Handle server shutdown gracefully:

```go
func gracefulShutdown(server *http.Server) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    <-c
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    server.Shutdown(ctx)
}
```

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure your server is running and accessible
2. **Upgrade Failed**: Check that proper WebSocket headers are sent
3. **Message Not Received**: Verify your MessageHandler implementation
4. **Memory Leaks**: Ensure proper client cleanup in your ClientManager

### Debug Mode

Enable debug logging to troubleshoot issues:

```go
import "log"

func (h *MyMessageHandler) Handle(client *ws.Client, data []byte) error {
    log.Printf("Received message from client %s: %s", client.ID, string(data))
    // Your handling logic
    return nil
}
```

## Performance Tips

1. **Buffer Sizes**: Adjust read/write buffer sizes based on your message patterns
2. **Connection Pooling**: Reuse connections when possible
3. **Message Batching**: Batch multiple small messages for better throughput
4. **Compression**: Enable WebSocket compression for large messages

## License

MIT License - see LICENSE file for details.
