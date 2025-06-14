package mq

import (
	"orders/models"
)

//type Inboxer interface {
//}
//
//type Outboxer interface {
//}

//type Event struct {
//	ID        uuid.UUID
//	Timestamp time.Time
//	Processed bool
//	Type      string
//	Payload   []byte
//}

//type InOutBoxer interface {
//	AddWith(*modelsEvent) error
//	Get() (*Event, error)
//	Complete(*Event) error
//}

type Broker interface {
	Send(*models.Event) error
	Receive(*models.Event) error
	Close() error
}
