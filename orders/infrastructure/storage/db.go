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

const (
	createQuery = `
		CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			status TEXT NOT NULL,
			price NUMERIC NOT NULL,
			descr TEXT
		);		
	`
	addQuery          = `INSERT INTO orders (id, user_id, status, price, descr) VALUES (:id, :user_id, :status, :price, :descr)`
	updateStatusQuery = `UPDATE orders SET status=$2 WHERE id=$1`
	allQuery          = `SELECT * FROM orders`
	getQuery          = `SELECT * FROM orders WHERE id=$1`
)

func NewOrderDB(db *sqlx.DB) (*OrderDB, error) {
	_, err := db.Exec(createQuery)

	if err != nil {
		return nil, err
	}

	return &OrderDB{db}, nil
}

func (odb *OrderDB) Get(id uuid.UUID) (order *model.Order, err error) {
	order = &model.Order{}
	err = odb.db.Get(order, getQuery, id)

	if err != nil {
		return nil, err
	}

	return
}

func (odb *OrderDB) AddWith(trx trx.Trx, order *model.Order) (err error) {
	tx, ok := trx.(*sqlx.Tx)
	if !ok {
		return fmt.Errorf("transaction is not a sqlx.Tx")
	}

	_, err = tx.NamedExec(addQuery, order)

	return
}
func (odb *OrderDB) Add(order *model.Order) (err error) {
	_, err = odb.db.NamedExec(addQuery, order)

	return
}
func (odb *OrderDB) All() (f []*model.Order, err error) {
	err = odb.db.Select(&f, allQuery)

	if err != nil {
		return nil, err
	}

	return
}

func (odb *OrderDB) UpdateStatus(id uuid.UUID, status model.Status) (err error) {
	_, err = odb.db.Exec(updateStatusQuery, id, status)
	return
}
