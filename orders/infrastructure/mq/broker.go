package mq

import (
	"orders/models"
)

type Broker interface {
	Send(*model.Event) error
	Receive() (*model.Event, error)
	Close() error
	Register() error
}
