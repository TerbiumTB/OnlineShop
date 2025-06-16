package service

import (
	"context"
	"log"

	//"encoding/json"
	"github.com/google/uuid"
	"orders/infrastructure/inoutbox"
	"orders/infrastructure/mq"
	"time"
)

type InboxWorker struct {
	broker mq.Broker
	inbox  inoutbox.Inboxer
	lg     *log.Logger
	//storage storage.AccountStorer
}

func NewInboxWorker(broker mq.Broker, inbox inoutbox.Inboxer, lg *log.Logger) *InboxWorker {
	return &InboxWorker{broker: broker, inbox: inbox, lg: lg}
}

type orderEvent struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Price  float64   `json:"price"`
}

func (w *InboxWorker) try() (ok bool) {
	event, err := w.broker.Receive()

	if err != nil {
		return false
	}

	err = w.inbox.Add(event)
	if err != nil {
		return false
	}

	return w.broker.Register() == nil
}

func (w *InboxWorker) Start(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if w.try() {
					w.lg.Println("added event to inbox")
				}
			}
		}
	}()
}
