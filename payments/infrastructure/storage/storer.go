package storage

import (
	"github.com/google/uuid"
	"payments/infrastructure/trx"
	"payments/models"
)

type AccountStorer interface {
	Add(*model.Account) error
	Get(uuid.UUID) (*model.Account, error)
	All() ([]*model.Account, error)
	Update(id uuid.UUID, amount float64) (err error)
	PayWith(trx.Trx, *model.Payment) error
}
