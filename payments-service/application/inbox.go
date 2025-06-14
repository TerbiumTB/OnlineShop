package application

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

func (w *InboxWorker) try() (err error) {
	e, err := w.broker.Receive()

	if err != nil {
		return err
	}

	err = w.inbox.Add(e)
	if err != nil {
		return err
	}

	return w.broker.Register()
}

func (w *InboxWorker) Start(ctx context.Context, handlePeriod time.Duration) {
	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				//	noop
			}

			if w.try() != nil {
				w.lg.Println("failed to add event to queue")
			}
			w.lg.Println("Tried")
		}
	}()
}
