package ws

import (
	"net/http"
)

type SessionInfo struct {
	ClientID Identity
	Metadata map[string]string
}

type SessionValidator interface {
	Validate(r *http.Request) (SessionInfo, error)
}
