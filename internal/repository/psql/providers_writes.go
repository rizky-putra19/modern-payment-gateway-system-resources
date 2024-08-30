package psql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/jmoiron/sqlx"
)

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

func (pw *ProviderWrites) UpdateProviderPaychannelByIdRepo(payload dto.AdjustLimitOrFeeProviderPayload) error {
	query := "UPDATE provider_paychannels SET "
	var conditions []string
	var args []interface{}

	if payload.Fee != nil {
		conditions = append(conditions, fmt.Sprintf("fee = $%d", len(args)+1))
		args = append(args, *payload.Fee)
	}

	if payload.FeeType != nil {
		conditions = append(conditions, fmt.Sprintf("fee_type = $%d", len(args)+1))
		args = append(args, *payload.FeeType)
	}

	if payload.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("min_transaction = $%d", len(args)+1))
		args = append(args, *payload.MinAmount)
	}

	if payload.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("max_transaction = $%d", len(args)+1))
		args = append(args, *payload.MaxAmount)
	}

	if payload.MaxDailyLimit != nil {
		conditions = append(conditions, fmt.Sprintf("max_daily_transaction = $%d", len(args)+1))
		args = append(args, *payload.MaxDailyLimit)
	}

	if payload.InterfaceSetting != nil {
		conditions = append(conditions, fmt.Sprintf("interface_setting = $%d", len(args)+1))
		args = append(args, *payload.InterfaceSetting)
	}

	if len(conditions) == 0 {
		return errors.New("no fields to update")
	}

	query += strings.Join(conditions, ", ")
	query += ", updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta' "
	query += "WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, payload.ProviderChannelId)

	_, err := pw.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pw *ProviderWrites) AddOperatorProviderChannelRepo(providerChannelId int, bankListId int) (int, error) {
	var id int

	query := `
	INSERT INTO provider_paychannel_bank_lists (provider_paychannel_id, bank_list_id, created_at, updated_at)
	VALUES ($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := pw.db.QueryRow(query, providerChannelId, bankListId)
	err := row.Scan(&id)
	if err != nil && id == 0 {
		return id, err
	}
	return id, nil
}

func (pw *ProviderWrites) DeleteOperatorProviderChannelRepo(providerChannelId int, bankListId int) error {
	query := `
	DELETE FROM provider_paychannel_bank_lists
	WHERE provider_paychannel_id = $1
		AND bank_list_id = $2;
	`

	_, err := pw.db.Exec(query, providerChannelId, bankListId)
	if err != nil {
		return err
	}
	return nil
}

func (pw *ProviderWrites) UpdateStatusProviderPaychannelRepo(id int, status string) error {
	query := `
	UPDATE provider_paychannels
	SET
		status = $1,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE id = $2;
	`

	_, err := pw.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}
