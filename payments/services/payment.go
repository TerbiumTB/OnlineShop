package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"payments/infrastructure/inoutbox"
	"payments/infrastructure/storage"
	"payments/infrastructure/trx"
	"payments/models"
	"time"
)

type PaymentWorker struct {
	accounts storage.AccountStorer
	inbox    inoutbox.Inboxer
	outbox   inoutbox.Outboxer
	manager  trx.Manager
	lg       *log.Logger
}

func NewPaymentWorker(accounts storage.AccountStorer, inbox inoutbox.Inboxer, outbox inoutbox.Outboxer, manager trx.Manager, lg *log.Logger) *PaymentWorker {
	return &PaymentWorker{
		accounts: accounts,
		inbox:    inbox,
		outbox:   outbox,
		manager:  manager,
		lg:       lg,
	}
}

func (s *PaymentWorker) Add(a *model.Account) error {
	return s.accounts.Add(a)
}

type paymentStatus struct {
	Id     uuid.UUID    `json:"id" db:"id"`
	Status model.Status `json:"status" db:"status"`
}

type paymentEvent struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Price  float64   `json:"price" db:"amount"`
}

func (s *PaymentWorker) pay() (ok bool) {
	event, err := s.inbox.Get()
	if err != nil || event == nil {
		return false
	}
	pe := &paymentEvent{}
	err = json.Unmarshal(event.Payload, pe)

	if err != nil {
		s.lg.Printf("Error unmarshalling event: %s", err)
		return false
	}
	payment := model.NewPayment(event.ID, pe.UserID, pe.Price)
	s.lg.Println("Payment: ", payment)

	tr, err := s.manager.Begin()
	if err != nil {
		s.lg.Printf("Error begining transaction: %s", err)
		return false
	}
	defer tr.Rollback()

	status := &paymentStatus{payment.ID, model.SUCCESS}

	err = s.accounts.PayWith(tr, payment)
	if err != nil {
		status.Status = model.FAIL
		s.lg.Println(err)
	}

	err = s.inbox.CompleteWith(tr, event)
	if err != nil {
		s.lg.Println(err)
		return false
	}

	paymentEvent, _ := model.NewEventWithJson(status)

	err = s.outbox.AddWith(tr, paymentEvent)
	if err != nil {
		s.lg.Println(err)
		return false
	}

	return tr.Commit() == nil
}

func (s *PaymentWorker) StartPaying(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if s.pay() {
					s.lg.Println("successful payment")
				}
			}
		}
	}()

}
