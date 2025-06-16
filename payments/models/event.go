package model

import (
	"encoding/json"
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

func NewEventWithType(typ string, payload []byte) *Event {
	return &Event{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Processed: false,
		Type:      typ,
		Payload:   payload,
	}
}

func NewEventWithJson(payload any) (*Event, error) {
	json, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Event{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Processed: false,
		Type:      "",
		Payload:   json,
	}, nil
}

func NewEvent(payload []byte) *Event {
	return &Event{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Processed: false,
		Type:      "",
		Payload:   payload,
	}
}

type PaymentEvent struct {
	Id      uuid.UUID `json:"id" db:"id"`
	Succeed bool      `json:"succeed" db:"succeed"`
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
