package storage

import (
	"github.com/google/uuid"
	"orders/infrastructure/trx"
	"orders/models"
)

type OrdersStorer interface {
	Add(*models.Order) error
	AddWith(trx.Trx, *models.Order) error
	Get(uuid.UUID) (*models.Order, error)
	All() ([]*models.Order, error)
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
