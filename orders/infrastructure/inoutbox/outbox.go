package inoutbox

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"orders/infrastructure/trx"
	"orders/models"
)

type Outboxer interface {
	AddWith(trx.Trx, *model.Event) error
	Get() (*model.Event, error)
	Complete(*model.Event) error
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

func (o *Outbox) AddWith(trx trx.Trx, e *model.Event) (err error) {
	tx, ok := trx.(*sqlx.Tx)

	if !ok {
		return fmt.Errorf(`trx is not a sqlx.Tx`)
	}
	_, err = tx.NamedExec(`INSERT INTO outbox (id, created_at, processed, type, payload) VALUES (:id, :created_at, :processed, :type, :payload)`, e)

	return
}

func (o *Outbox) Get() (e *model.Event, err error) {
	e = &model.Event{}
	err = o.db.Get(e, `SELECT * FROM outbox WHERE processed = false LIMIT 1`)
	return
}

func (o *Outbox) Complete(e *model.Event) (err error) {
	_, err = o.db.NamedExec(`UPDATE outbox SET processed = true WHERE id = :id`, e)
	return
}
