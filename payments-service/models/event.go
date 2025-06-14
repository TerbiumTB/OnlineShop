package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Processed bool      `json:"processed" db:"processed"`
	Type      string    `json:"type" db:"type"`
	Payload   []byte    `json:"payload" db:"payload"`
}

func NewEvent(typ string, payload []byte) *Event {
	return &Event{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Processed: false,
		Type:      typ,
		Payload:   payload,
	}
}

func NewEventWithID(id uuid.UUID, payload []byte) *Event {
	return &Event{
		ID:        id,
		CreatedAt: time.Now(),
		Processed: false,
		Type:      "",
		Payload:   payload,
	}
}
