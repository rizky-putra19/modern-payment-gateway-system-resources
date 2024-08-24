package psql

import "github.com/jmoiron/sqlx"

type ProviderWrites struct {
	db *sqlx.DB
}

func NewProviderWrites(db *sqlx.DB) *ProviderWrites {
	return &ProviderWrites{
		db: db,
	}
}

func (pw *ProviderWrites) CreateProviderConfirmationDetail(source string, paymentId string, status string) (int, error) {
	var id int

	query := `
	INSERT INTO provider_transaction_confirmation_details (payment_id, type, confirmation_result, created_at, updated_at)
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := pw.db.QueryRow(query, paymentId, source, status)
	err := row.Scan(&id)
	if err != nil || id == 0 {
		return id, err
	}

	return id, nil
}
