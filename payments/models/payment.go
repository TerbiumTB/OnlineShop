package model

import "github.com/google/uuid"

type Payment struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Amount float64   `json:"amount" db:"amount"`
}

func NewPayment(id uuid.UUID, userID uuid.UUID, amount float64) *Payment {
	return &Payment{
		ID:     id,
		UserID: userID,
		Amount: amount,
	}
}
