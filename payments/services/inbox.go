package service

import (
	"context"
	"log"

	//"encoding/json"
	"github.com/google/uuid"
	"payments/infrastructure/inoutbox"
	"payments/infrastructure/mq"
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
	w.lg.Println("received event:", event)

	if err != nil {
		return false
	}

	w.lg.Println("adding event...")
	err = w.inbox.Add(event)
	if err != nil {
		return false
	}
	w.lg.Println("compliting event...")
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
				//noop
			}
			if w.try() {
				w.lg.Println("added event to inbox")
			} else {
				w.lg.Println("tried")
			}
		}
	}()
}
