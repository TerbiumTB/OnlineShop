package model

import "github.com/google/uuid"

type Account struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	FullName string    `json:"full_name" db:"full_name"`
	Balance  float64   `json:"balance" db:"balance"`
}

func NewAccount(userID uuid.UUID, fullname string, balance float64) *Account {
	return &Account{
		UserID:   userID,
		FullName: fullname,
		Balance:  balance,
	}
}
