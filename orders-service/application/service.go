package application

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"orders/infrastructure/inoutbox"
	"orders/infrastructure/storage"
	"orders/infrastructure/trx"
	"orders/models"
)

type Service struct {
	storage storage.OrdersStorer
	outbox  inoutbox.Outboxer
	manager trx.Manager
}

func NewService(storage storage.OrdersStorer, outbox inoutbox.Outboxer, manager trx.Manager) *Service {
	return &Service{storage: storage, outbox: outbox, manager: manager}
}

type orderPayload struct {
	//ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Price  float64   `json:"price"`
}

func (s *Service) newOrderEvent(order *models.Order) *models.Event {
	op := orderPayload{
		UserID: order.UserID,
		Price:  order.Price,
	}

	payload, _ := json.Marshal(op)
	return models.NewEvent("NewOrder", payload)
}

func (s *Service) Add(userID string, price float64, descr string) (err error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	if price < 0 {
		return fmt.Errorf("price must be greater than zero")
	}

	order := models.NewOrder(uid, price, descr)

	trx, err := s.manager.Begin()

	defer trx.Rollback()
	if err != nil {
		return
	}

	err = s.storage.AddWith(trx, order)
	if err != nil {
		return
	}

	event := s.newOrderEvent(order)

	err = s.outbox.AddWith(trx, event)
	if err != nil {
		return
	}

	return trx.Commit()
}
