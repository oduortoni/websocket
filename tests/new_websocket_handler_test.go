package tests

import (
	"testing"

	"github.com/oduortoni/websocket/ws"
)

func TestNewWebSocketHandler(t *testing.T) {
	validator := &mockSessionValidator{}
	messageHandler := &mockMessageHandler{}
	persister := &mockEnvelopePersister{}

	handler := ws.NewWebSocketHandler(validator, messageHandler, persister)

	if handler.SessionValidator != validator {
		t.Error("Expected validator to be set correctly")
	}
	if handler.MessageHandler != messageHandler {
		t.Error("Expected message handler to be set correctly")
	}
	if handler.EnvelopePersister != persister {
		t.Error("Expected persister to be set correctly")
	}
}
