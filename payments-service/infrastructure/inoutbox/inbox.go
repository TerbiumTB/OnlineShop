package inoutbox

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"payments/infrastructure/trx"
	"payments/models"
)

type Inboxer interface {
	Add(*models.Event) error
	Get() (*models.Event, error)
	CompleteWith(trx.Trx, *models.Event) error
}

type Inbox struct {
	db *sqlx.DB
}

func NewInbox(db *sqlx.DB) (*Inbox, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS inbox (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP,
			processed BOOLEAN NOT NULL,
			type TEXT NOT NULL,
			payload JSONB
		);
		
	`)

	if err != nil {
		return nil, err
	}

	return &Inbox{db}, nil
}

func (i *Inbox) Add(e *models.Event) (err error) {
	_, err = i.db.NamedExec(`INSERT INTO outbox (id, created_at, processed, type, payload) VALUES (:id, :created_at, :processed, :type, :payload)`, e)

	return
}
func (i *Inbox) Get() (*models.Event, error) {
	return nil, fmt.Errorf("not implemented")
}
func (i *Inbox) CompleteWith(trx.Trx, *models.Event) error {
	return fmt.Errorf("not implemented")
}
