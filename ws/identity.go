package ws

import (
	"github.com/google/uuid"
)

type Identity uuid.UUID

func NewIdentity() Identity {
	return Identity(uuid.New())
}

func (i Identity) String() string {
	return uuid.UUID(i).String()
}

func (i Identity) UUID() uuid.UUID {
	return uuid.UUID(i)
}
