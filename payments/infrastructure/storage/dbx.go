package storage

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"payments/infrastructure/trx"
	"payments/models"
)

type AccountDB struct {
	db *sqlx.DB
}

const (
	createQuery = `
		CREATE TABLE IF NOT EXISTS accounts (
			user_id UUID NOT NULL,
			full_name TEXT NOT NULL,
			balance NUMERIC NOT NULL
		);		
	`
	addQuery     = `INSERT INTO accounts (user_id, full_name, balance) VALUES (:user_id, :full_name, :balance)`
	paymentQuery = `UPDATE accounts SET balance = balance - :amount WHERE (user_id = :user_id) AND (balance > :amount)`
	updateQuery  = `UPDATE accounts SET balance = balance + $2 WHERE user_id = $1`
	allQuery     = `SELECT * FROM accounts`
	getQuery     = `SELECT * FROM accounts WHERE id = $1`
)

func NewAccountDB(db *sqlx.DB) (*AccountDB, error) {
	_, err := db.Exec(createQuery)

	if err != nil {
		return nil, err
	}

	return &AccountDB{db}, nil
}

func (adb *AccountDB) Get(id uuid.UUID) (account *model.Account, err error) {
	account = &model.Account{}

	err = adb.db.Get(account, getQuery, id)

	if err != nil {
		return nil, err
	}

	return
}

func (adb *AccountDB) AddWith(t trx.Transaction, account *model.Account) (err error) {
	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return trx.SqlxError
	}

	_, err = tx.NamedExec(addQuery, account)

	return
}
func (adb *AccountDB) Add(account *model.Account) (err error) {
	_, err = adb.db.NamedExec(addQuery, account)

	return
}
func (adb *AccountDB) All() (f []*model.Account, err error) {
	err = adb.db.Select(&f, allQuery)

	if err != nil {
		return nil, err
	}

	return
}

func (adb *AccountDB) Update(id uuid.UUID, amount float64) (err error) {
	_, err = adb.db.Exec(updateQuery, id, amount)
	return
}

func (adb *AccountDB) PayWith(t trx.Transaction, payment *model.Payment) (err error) {
	tx, ok := t.(*sqlx.Tx)
	if !ok {
		return trx.SqlxError
	}

	res, err := tx.NamedExec(paymentQuery, payment)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil || rows == 0 {
		return errors.New("couldn't proceed payment")
	}

	return
}
