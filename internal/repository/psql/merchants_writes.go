package psql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/jmoiron/sqlx"
)

type MerchantWrites struct {
	db *sqlx.DB
}

func NewMerchantWrites(db *sqlx.DB) *MerchantWrites {
	return &MerchantWrites{
		db: db,
	}
}

func (mw *MerchantWrites) UpdateMerchantCapitalAndNotSettleBalance(notSettleBalance float64, balanceCapitalFlow float64, merchantId string) error {
	query := `
	UPDATE merchant_accounts
	SET not_settle_balance = $1,
		balance_capital_flow = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE merchant_id = $3;
	`

	_, err := mw.db.Exec(query, notSettleBalance, balanceCapitalFlow, merchantId)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MerchantWrites) UpdateMerchantSecretKeyRepo(secretKey string, merchantId string) error {
	query := `
	UPDATE merchants
	SET
		merchant_secret = $1,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE merchant_id = $2;
	`

	_, err := mw.db.Exec(query, secretKey, merchantId)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MerchantWrites) UpdateMerchantCapitalAndSettleBalance(settleBalance float64, balanceCapitalFlow float64, merchantId string) error {
	query := `
	UPDATE merchant_accounts
	SET settle_balance = $1,
		balance_capital_flow = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE merchant_id = $3;
	`

	_, err := mw.db.Exec(query, settleBalance, balanceCapitalFlow, merchantId)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MerchantWrites) UpdateMerchantCapitalPendingOut(pendingAmount float64, balanceCapitalFlow float64, merchantId string) error {
	query := `
	UPDATE merchant_accounts
	SET balance_capital_flow = $1,
		pending_transaction_out = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE merchant_id = $3;
	`

	_, err := mw.db.Exec(query, balanceCapitalFlow, pendingAmount, merchantId)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MerchantWrites) UpdateMerchantHoldBalanceAndSettleBalance(settleBalance float64, holdBalance float64, merchantId string) error {
	query := `
	UPDATE merchant_accounts
	SET settle_balance = $1,
		hold_balance = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE merchant_id = $3;
	`

	_, err := mw.db.Exec(query, settleBalance, holdBalance, merchantId)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MerchantWrites) UpdateMerchantHoldBalanceAndNotSettleBalance(notSettleBalance float64, holdBalance float64, merchantId string) error {
	query := `
	UPDATE merchant_accounts
	SET not_settle_balance = $1,
		hold_balance = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE merchant_id = $3;
	`

	_, err := mw.db.Exec(query, notSettleBalance, holdBalance, merchantId)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MerchantWrites) CreateMerchantCapitalFlow(payload dto.CreateMerchantCapitalFlowPayload) (int, error) {
	var merchantCapitalFlowId int
	query := `
	INSERT INTO merchant_capital_flows (payment_id, merchant_account_id, temp_balance, amount, reason_id, status, created_by, capital_type, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	if payload.Notes != "" {
		query = `
		INSERT INTO merchant_capital_flows (payment_id, merchant_account_id, temp_balance, amount, reason_id, status, notes, created_by, capital_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
		RETURNING id
		`

		if payload.ReverseFrom != "" {
			query = `
			INSERT INTO merchant_capital_flows (payment_id, merchant_account_id, temp_balance, amount, reason_id, status, notes, created_by, capital_type, reverse_from, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
			RETURNING id
			`

			row := mw.db.QueryRow(query,
				payload.PaymentId,
				payload.MerchantAccountId,
				payload.TempBalance,
				payload.Amount,
				payload.ReasonId,
				payload.Status,
				payload.Notes,
				payload.CreateBy,
				payload.CapitalType,
				payload.ReverseFrom)
			err := row.Scan(&merchantCapitalFlowId)
			if err != nil || merchantCapitalFlowId == 0 {
				return 0, err
			}

			return merchantCapitalFlowId, nil
		}

		row := mw.db.QueryRow(query,
			payload.PaymentId,
			payload.MerchantAccountId,
			payload.TempBalance,
			payload.Amount,
			payload.ReasonId,
			payload.Status,
			payload.Notes,
			payload.CreateBy,
			payload.CapitalType)
		err := row.Scan(&merchantCapitalFlowId)
		if err != nil || merchantCapitalFlowId == 0 {
			return 0, err
		}
	}

	if payload.Notes == "" {
		row := mw.db.QueryRow(query,
			payload.PaymentId,
			payload.MerchantAccountId,
			payload.TempBalance,
			payload.Amount,
			payload.ReasonId,
			payload.Status,
			payload.CreateBy,
			payload.CapitalType)
		err := row.Scan(&merchantCapitalFlowId)
		if err != nil || merchantCapitalFlowId == 0 {
			return 0, err
		}
	}

	return merchantCapitalFlowId, nil
}

func (mw *MerchantWrites) UpdateMerchantSettlement(settleBalance float64, notSettleBalance float64, merchantId string) error {
	query := `
	UPDATE merchant_accounts
	SET settle_balance = $1,
		not_settle_balance = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE merchant_id = $3;
	`

	_, err := mw.db.Exec(query, settleBalance, notSettleBalance, merchantId)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MerchantWrites) CreateMerchantCallback(paymentId string, callbackStatus string, paymentStatusInCallback string, callbackResult string, triggerBy string) (int, error) {
	var merchantCallbackId int

	query := `
	INSERT INTO merchant_callbacks (payment_id, callback_status, payment_status_in_callback, callback_result, triggered_by, created_at)
	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := mw.db.QueryRow(query, paymentId, callbackStatus, paymentStatusInCallback, callbackResult, triggerBy)
	err := row.Scan(&merchantCallbackId)
	if err != nil || merchantCallbackId == 0 {
		return merchantCallbackId, err
	}

	return merchantCallbackId, nil
}

func (mw *MerchantWrites) CreateMerchantRepo(merchantName string, merchantId string, merchantSecret string) (int, error) {
	var createMerchantId int

	query := `
	INSERT INTO merchants (merchant_id, merchant_name, merchant_secret, currency, status, created_at, updated_at)
	VALUES ($1, $2, $3, 'IDR', 'INACTIVE', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := mw.db.QueryRow(query, merchantId, merchantName, merchantSecret)
	err := row.Scan(&createMerchantId)
	if err != nil || createMerchantId == 0 {
		return createMerchantId, err
	}

	return createMerchantId, nil
}

func (mw *MerchantWrites) CreateMerchantPaymentMethodRepo(merchantId int, paymentMethodId int) (int, error) {
	var createPaymentMethodId int

	query := `
	INSERT INTO merchant_payment_methods (merchant_id, payment_method_id, created_at, updated_at)
	VALUES ($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := mw.db.QueryRow(query, merchantId, paymentMethodId)
	err := row.Scan(&createPaymentMethodId)
	if err != nil || createPaymentMethodId == 0 {
		return createPaymentMethodId, err
	}

	return createPaymentMethodId, nil
}

func (mw *MerchantWrites) CreateMerchantPaychannelRepo(merchantPaymentMethodId int, segment string, fee float64, feeType string, minAmount float64, maxAmount float64, dailyLimit float64, merchantPaychannelCode string) (int, error) {
	var merchantPaychannelId int

	query := `
	INSERT INTO merchant_paychannels (merchant_payment_method_id, segment, fee, fee_type, min_transaction, max_transaction, max_daily_transaction, merchant_paychannel_code, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'INACTIVE', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := mw.db.QueryRow(query, merchantPaymentMethodId, segment, fee, feeType, minAmount, maxAmount, dailyLimit, merchantPaychannelCode)
	err := row.Scan(&merchantPaychannelId)
	if err != nil || merchantPaychannelId == 0 {
		return merchantPaychannelId, err
	}

	return merchantPaychannelId, nil
}

func (mw *MerchantWrites) CreateMerchantAccountsRepo(merchantId string) (int, error) {
	var merchantAccountId int

	query := `
	INSERT INTO merchant_accounts (merchant_id, created_at, updated_at)
	VALUES ($1, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := mw.db.QueryRow(query, merchantId)
	err := row.Scan(&merchantAccountId)
	if err != nil || merchantAccountId == 0 {
		return merchantAccountId, err
	}

	return merchantAccountId, nil
}

func (mw *MerchantWrites) UpdateMerchantStatusRepo(merchantId string, status string) error {
	query := `
	UPDATE merchants
	SET 
		status = $1,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE merchant_id = $2;
	`

	_, err := mw.db.Exec(query, status, merchantId)
	if err != nil {
		return err
	}

	return nil
}

func (mw *MerchantWrites) UpdateMerchantPaychannelByIdRepo(payload dto.AdjustLimitOrFeePayload) error {
	query := "UPDATE merchant_paychannels SET "
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

	if len(conditions) == 0 {
		return errors.New("no fields to update")
	}

	query += strings.Join(conditions, ", ")
	query += ", updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta' "
	query += "WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, payload.MerchantPaychannelId)

	_, err := mw.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (mw *MerchantWrites) UpdateStatusMerchantPaychannelById(id int, status string) error {
	query := `
	UPDATE merchant_paychannels
	SET
		status = $1,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE id = $2;
	`

	_, err := mw.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (mw *MerchantWrites) DeleteRoutingPaychannelByMerchantPaychannelId(id int) error {
	query := `
	DELETE FROM paychannel_routings
	WHERE merchant_paychannel_id = $1;
	`
	_, err := mw.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (mw *MerchantWrites) AddRoutingPaychannelRepo(merchantPaychannelId int, providerPaychannelId int) (int, error) {
	var routingPaychannelId int

	query := `
	INSERT INTO paychannel_routings (provider_paychannel_id, merchant_paychannel_id, created_at, updated_at)
	VALUES ($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := mw.db.QueryRow(query, providerPaychannelId, merchantPaychannelId)
	err := row.Scan(&routingPaychannelId)
	if err != nil || routingPaychannelId == 0 {
		return routingPaychannelId, err
	}

	return routingPaychannelId, nil
}
