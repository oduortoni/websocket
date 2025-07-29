package tests

import (
	"errors"
	"net/http"

	"github.com/oduortoni/websocket/ws"
)

type mockSessionValidator struct {
	shouldFail bool
	session    ws.SessionInfo
}

func (m *mockSessionValidator) Validate(r *http.Request) (ws.SessionInfo, error) {
	if m.shouldFail {
		return ws.SessionInfo{}, errors.New("unauthorized")
	}
	return m.session, nil
}

type mockMessageHandler struct {
	shouldFail bool
	messages   [][]byte
}

func (m *mockMessageHandler) Handle(client *ws.Client, data []byte) error {
	if m.shouldFail {
		return errors.New("message handling failed")
	}
	m.messages = append(m.messages, data)
	return nil
}

type mockEnvelopePersister struct {
	shouldFail bool
	envelopes  []ws.Envelope
}

func (m *mockEnvelopePersister) SaveEnvelope(e ws.Envelope) error {
	if m.shouldFail {
		return errors.New("save failed")
	}
	m.envelopes = append(m.envelopes, e)
	return nil
}

func (m *mockEnvelopePersister) ConfirmDelivery(envelopeID ws.Identity, clientID ws.Identity) error {
	if m.shouldFail {
		return errors.New("confirm delivery failed")
	}
	return nil
}
