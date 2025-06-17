package service

import (
	"context"
	"log"
	"payments/infrastructure/inoutbox"
	"payments/infrastructure/mq"
	"time"
)

type OutboxWorker struct {
	broker mq.Broker
	outbox inoutbox.Outboxer
	lg     *log.Logger
}

func NewOutboxWorker(broker mq.Broker, outbox inoutbox.Outboxer, lg *log.Logger) *OutboxWorker {
	return &OutboxWorker{broker: broker, outbox: outbox, lg: lg}
}

func (w *OutboxWorker) try() error {
	e, err := w.outbox.Get()
	if err != nil || e == nil {
		//w.lg.Printf("Couldn't get event: %v", err)
		return err
	}

	err = w.broker.Send(e)
	if err != nil {
		//w.lg.Printf("Couldn't send event: %v", err)
		return err
	}
	_ = w.outbox.Complete(e)
	return nil
}

func (w *OutboxWorker) Start(ctx context.Context, handlePeriod time.Duration) {
	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := w.try()

				if err != nil {
					//w.lg.Printf("Outbox worker error: %s", err)
				} else {
					w.lg.Printf("Event completed successfully")
				}
			}
		}
	}()
}
