package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"orders/infrastructure/inoutbox"
	"orders/infrastructure/storage"
	"orders/infrastructure/trx"
	"orders/models"
)

type OrderServicing interface {
	Add(string, float64, string) error
	Get(userID string) (*model.Order, error)
	All() ([]*model.Order, error)
	UpdateStatus(string, string) error
}

type OrderService struct {
	storage storage.OrdersStorer
	outbox  inoutbox.Outboxer
	manager trx.Manager
}

func NewOrderService(storage storage.OrdersStorer, outbox inoutbox.Outboxer, manager trx.Manager) *OrderService {
	return &OrderService{storage: storage, outbox: outbox, manager: manager}
}

type orderPayload struct {
	//ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Price  float64   `json:"price"`
}

func (s *OrderService) newOrderEvent(order *model.Order) *model.Event {
	op := orderPayload{
		UserID: order.UserID,
		Price:  order.Price,
	}

	payload, _ := json.Marshal(op)
	return model.NewEventWithID(order.ID, payload)
}

func (s *OrderService) Add(userID string, price float64, descr string) (err error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	if price < 0 {
		return fmt.Errorf("price must be greater than zero")
	}

	order := model.NewOrder(uid, price, descr)

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

func (s *OrderService) Get(userID string) (order *model.Order, err error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return
	}
	order, err = s.storage.Get(uid)
	return
}

func (s *OrderService) All() ([]*model.Order, error) {
	return s.storage.All()
}

func (s *OrderService) UpdateStatus(id string, status string) (err error) {
	sstatus, err := model.ParseStatus(status)
	if err != nil {
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return
	}

	return s.storage.UpdateStatus(uid, sstatus)
}
