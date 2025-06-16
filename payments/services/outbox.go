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

func (w *OutboxWorker) try() (ok bool) {
	e, err := w.outbox.Get()
	if err != nil {
		return false
	}

	err = w.broker.Send(e)
	if err != nil {
		return false
	}
	_ = w.outbox.Complete(e)
	return true
}

func (w *OutboxWorker) Start(ctx context.Context, handlePeriod time.Duration) {
	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if w.try() {
					w.lg.Printf("Event completed successfully")
				}
			}
		}
	}()
}
