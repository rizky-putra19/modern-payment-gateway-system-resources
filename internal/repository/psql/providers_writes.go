package psql

import "github.com/jmoiron/sqlx"

type ProviderWrites struct {
	db *sqlx.DB
}

func NewProviderWrites(db *sqlx.DB) *TransactionsReads {
	return &TransactionsReads{
		db: db,
	}
}
