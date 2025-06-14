package models

import "github.com/google/uuid"

type Account struct {
	//ID uuid.UUID	`json:"id"`
	UserID  uuid.UUID `json:"user_id" db:"user_id"`
	Balance float64   `json:"balance" db:"balance"`
}

func NewAccount(userID uuid.UUID, amount float64) *Account {
	return &Account{
		UserID:  userID,
		Balance: amount,
	}
}
