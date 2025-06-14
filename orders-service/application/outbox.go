package application

import (
	"context"
	"log"
	"orders/infrastructure/inoutbox"
	"orders/infrastructure/mq"
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

func (o *OutboxWorker) Start(ctx context.Context, handlePeriod time.Duration) {
	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				//	noop
			}

			e, err := o.outbox.Get()
			if err != nil || e == nil {
				o.lg.Printf("Couldn't get event: %v", err)
				continue
			}

			err = o.broker.Send(e)
			if err != nil {
				o.lg.Printf("Couldn't send event: %v", err)
				continue
			}

			_ = o.outbox.Complete(e)
			o.lg.Printf("Event completed successfully")

		}
	}()
}
