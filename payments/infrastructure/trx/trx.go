package trx

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

var SqlxError = errors.New("transaction is not a sqlx.Tx object")

type Transaction interface {
	Commit() error
	Rollback() error
}

type Manager interface {
	Begin() (Transaction, error)
}

type ManagerDB struct {
	db *sqlx.DB
}

func NewDBManager(db *sqlx.DB) *ManagerDB {
	return &ManagerDB{db: db}
}

func (m *ManagerDB) Begin() (Transaction, error) {
	return m.db.Beginx()
}
