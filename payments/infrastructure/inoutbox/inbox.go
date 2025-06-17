package inoutbox

import (
	"github.com/jmoiron/sqlx"
	"payments/infrastructure/trx"
	"payments/models"
)

type Inboxer interface {
	Add(*model.Event) error
	Get() (*model.Event, error)
	CompleteWith(trx.Transaction, *model.Event) error
}

type Inbox struct {
	db *sqlx.DB
}

const (
	createQuery = `
		CREATE TABLE IF NOT EXISTS inbox (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP,
			processed BOOLEAN NOT NULL,
			type TEXT NOT NULL,
			payload JSONB
		);
		
	`
	addQuery      = `INSERT INTO inbox (id, created_at, processed, type, payload) VALUES (:id, :created_at, :processed, :type, :payload)`
	getQuery      = `SELECT * FROM inbox WHERE processed = false LIMIT 1`
	completeQuery = `UPDATE inbox SET processed = true WHERE id = :id`
)

func NewInbox(db *sqlx.DB) (*Inbox, error) {
	_, err := db.Exec(createQuery)

	if err != nil {
		return nil, err
	}

	return &Inbox{db}, nil
}

func (i *Inbox) Add(e *model.Event) (err error) {
	_, err = i.db.NamedExec(addQuery, e)

	return
}
func (i *Inbox) Get() (event *model.Event, err error) {
	event = &model.Event{}
	err = i.db.Get(event, getQuery)
	return
}

func (i *Inbox) CompleteWith(t trx.Transaction, event *model.Event) (err error) {
	tx, ok := t.(*sqlx.Tx)

	if !ok {
		return trx.SqlxError
	}

	_, err = tx.NamedExec(completeQuery, event)
	return
}
