package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"orders/infrastructure/mq"
	"orders/infrastructure/storage"
	model "orders/models"
	"time"
)

type StatusWorker struct {
	broker mq.Broker
	orders storage.OrdersStorer
	lg     *log.Logger
}

func NewStatusWorker(broker mq.Broker, orders storage.OrdersStorer, lg *log.Logger) *StatusWorker {
	return &StatusWorker{broker, orders, lg}
}

type statusEvent struct {
	UserID uuid.UUID    `json:"user_id" db:"user_id"`
	Status model.Status `json:"status" db:"status"`
}

func (w *StatusWorker) try() bool {
	event, err := w.broker.Receive()

	if err != nil {
		return false
	}
	status := &statusEvent{}
	err = json.Unmarshal(event.Payload, status)
	if err != nil {
		return false
	}

	err = w.orders.UpdateStatus(status.UserID, status.Status)
	return err == nil
}

func (w *StatusWorker) Start(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if w.try() {
					w.lg.Println("updated status")
				}
			}
		}
	}()
}
