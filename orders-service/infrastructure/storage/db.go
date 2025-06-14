package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"orders/infrastructure/trx"
	"orders/models"
)

type OrderDB struct {
	db *sqlx.DB
}

func NewOrderDB(db *sqlx.DB) (*OrderDB, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			status TEXT NOT NULL,
			price NUMERIC NOT NULL,
			descr TEXT
		);
		
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

	return &OrderDB{db}, nil
}

func (odb *OrderDB) Get(id uuid.UUID) (order *models.Order, err error) {
	order = &models.Order{}
	err = odb.db.Get(order, `SELECT * FROM orders WHERE id = $1`, id)

	if err != nil {
		return nil, err
	}

	return
}

const (
	insertOrderQuery = `INSERT INTO orders (id, user_id, status, price, descr) VALUES (:id, :user_id, :status, :price, :descr)`
)

func (odb *OrderDB) AddWith(trx trx.Trx, order *models.Order) (err error) {
	tx, ok := trx.(*sqlx.Tx)
	if !ok {
		return fmt.Errorf("tx is not a transaction")
	}

	_, err = tx.NamedExec(insertOrderQuery, order)

	return
}
func (odb *OrderDB) Add(order *models.Order) (err error) {
	_, err = odb.db.NamedExec(insertOrderQuery, order)

	return
}
func (odb *OrderDB) All() (f []*models.Order, err error) {
	err = odb.db.Select(&f, `SELECT * FROM orders`)

	if err != nil {
		return nil, err
	}

	return
}

//func (odb *OrderDB) addEvent(tx *sqlx.Tx, order *models.Order) (err error) {
//	event := newOrderEvent(order)
//	_, err = tx.NamedExec(`INSERT INTO outbox (id, created_at, processed, type, payload) VALUES (:id, :created_at, :processed, :type, :payload)`, event)
//	return
//}
//
//
//func (odb *OrderDB) GetEvent() (e *models.Event, err error) {
//	e = &models.Event{}
//	err = odb.db.Get(e, `SELECT * FROM outbox WHERE processed = false LIMIT 1`)
//	return
//}
//
//func (odb *OrderDB) Complete(e *models.Event) (err error) {
//	_, err = odb.db.NamedExec(`UPDATE outbox SET processed = true WHERE id = :id`, e)
//	return
//}
