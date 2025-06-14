package trx

import "github.com/jmoiron/sqlx"

type Trx interface {
	Commit() error
	Rollback() error
}

type Manager interface {
	Begin() (Trx, error)
}

type ManagerDB struct {
	db *sqlx.DB
}

func NewDBManager(db *sqlx.DB) *ManagerDB {
	return &ManagerDB{db: db}
}

func (m *ManagerDB) Begin() (Trx, error) {
	return m.db.Begin()
}