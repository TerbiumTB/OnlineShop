package model

import "github.com/google/uuid"

type Order struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Price       float64   `json:"price" db:"price"`
	Status      Status    `json:"status" db:"status"`
	Description string    `json:"descr" db:"descr, omitempty"`
}

func NewOrder(userID uuid.UUID, price float64, descr string) *Order {
	return &Order{
		ID:          uuid.New(),
		UserID:      userID,
		Price:       price,
		Status:      CREATED,
		Description: descr,
	}
}
