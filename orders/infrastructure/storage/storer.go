package storage

import (
	"github.com/google/uuid"
	"orders/infrastructure/trx"
	"orders/models"
)

type OrdersStorer interface {
	Add(*model.Order) error
	AddWith(trx.Trx, *model.Order) error
	Get(uuid.UUID) (*model.Order, error)
	All() ([]*model.Order, error)
	UpdateStatus(uuid.UUID, model.Status) error
}

//type orderPayload struct {
//	UserID uuid.UUID `json:"user_id"`
//	Price  float64   `json:"price"`
//}
//
//func newOrderEvent(order *models.Order) *models.Event {
//	op := orderPayload{
//		UserID: order.UserID,
//		Price:  order.Price,
//	}
//
//	payload, _ := json.Marshal(op)
//	return models.NewEvent("NewOrder", payload)
//}
