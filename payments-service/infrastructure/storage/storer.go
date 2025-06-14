package storage

import (
	"github.com/google/uuid"
	"payments/models"
)

type AccountStorer interface {
	Add(*models.Account) error
	Get(uuid.UUID) (*models.Account, error)
	All() ([]*models.Account, error)
	Update(uuid.UUID, float64) error
	//GetEvent() (*models.Event, error)
	//Complete(*models.Event) error
}
