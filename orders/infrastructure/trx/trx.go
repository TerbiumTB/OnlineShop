package trx

import "github.com/jmoiron/sqlx"

type Trx interface {
	Commit() error
	Rollback() error
}

type Manager interface {
	Begin() (Trx, error)
}

type DBManager struct {
	db *sqlx.DB
}

func NewDBManager(db *sqlx.DB) *DBManager {
	return &DBManager{db: db}
}

func (m *DBManager) Begin() (Trx, error) {
	return m.db.Beginx()
}
