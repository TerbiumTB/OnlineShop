package inoutbox

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"orders/infrastructure/trx"
	"orders/models"
)

type Outboxer interface {
	AddWith(trx.Trx, *models.Event) error
	Get() (*models.Event, error)
	Complete(*models.Event) error
}

type Outbox struct {
	db *sqlx.DB
}

func NewOutbox(db *sqlx.DB) (*Outbox, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS outbox (
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

	return &Outbox{db}, nil
}

func (o *Outbox) AddWith(trx trx.Trx, e *models.Event) (err error) {
	tx, ok := trx.(*sqlx.Tx)

	if !ok {
		return fmt.Errorf(`trx is not a sqlx.Tx`)
	}
	_, err = tx.NamedExec(`INSERT INTO outbox (id, created_at, processed, type, payload) VALUES (:id, :created_at, :processed, :type, :payload)`, e)

	return
}

func (o *Outbox) Get() (e *models.Event, err error) {
	e = &models.Event{}
	err = o.db.Get(e, `SELECT * FROM outbox WHERE processed = false LIMIT 1`)
	return
}

func (o *Outbox) Complete(e *models.Event) (err error) {
	_, err = o.db.NamedExec(`UPDATE outbox SET processed = true WHERE id = :id`, e)
	return
}
