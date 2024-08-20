package psql

import (
	"fmt"
	"strings"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/jmoiron/sqlx"
)

type TransactionsWrites struct {
	db *sqlx.DB
}

func NewTransactionsWrites(db *sqlx.DB) *TransactionsWrites {
	return &TransactionsWrites{
		db: db,
	}
}

func (tr *TransactionsWrites) UpdateStatus(status string, paymentId string) error {
	query := `
	UPDATE transactions
	SET status = $1,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE payment_id = $2;
	`
	_, err := tr.db.Exec(query, status, paymentId)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TransactionsWrites) CreateTransactionStatusLog(paymentId string, statusLog string, changeBy string, notes string) (int, error) {
	var transactionStatusLogsId int
	query := `
	INSERT INTO transaction_status_logs (payment_id, status_log, change_by, notes, real_notes, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	// check if notes didn't input by user
	if notes == "" {
		query = `
	INSERT INTO transaction_status_logs (payment_id, status_log, change_by, created_at, updated_at)
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`
		row := tr.db.QueryRow(query, paymentId, statusLog, changeBy)
		err := row.Scan(&transactionStatusLogsId)
		if err != nil || transactionStatusLogsId == 0 {
			return transactionStatusLogsId, err
		}
	}

	if notes != "" {
		row := tr.db.QueryRow(query, paymentId, statusLog, changeBy, notes, notes)
		err := row.Scan(&transactionStatusLogsId)
		if err != nil || transactionStatusLogsId == 0 {
			return transactionStatusLogsId, err
		}
	}

	return transactionStatusLogsId, nil
}

func (tr *TransactionsWrites) CreateMerchantExportCapitalFlowRepo(payload dto.CreateMerchantExportReqDto) ([]entity.MerchantExportCapitalFlowEntity, error) {
	var listData []entity.MerchantExportCapitalFlowEntity
	var query string

	if payload.ExportType == constant.ExportTypeCapitalFlow {
		query = `
		SELECT
			mcf.ID,
			mcf.payment_id,
			m.merchant_id,
			m.merchant_name,
			mcf.amount,
			rl.reason_name,
			pm.name AS payment_method,
			mp.fee,
			mp.fee_type,
			mcf.temp_balance,
			p.provider_name,
			pp.paychannel_name,
			pp.fee AS provider_fee,
			pp.fee_type AS provider_fee_type,
			mcf.status,
			mcf.notes,
			mcf.capital_type,
			mcf.reverse_from,
			mcf.created_at
		FROM
			merchant_capital_flows mcf
		LEFT JOIN transactions t ON t.payment_id = mcf.payment_id
		LEFT JOIN merchant_accounts ma ON ma.ID = mcf.merchant_account_id
		LEFT JOIN merchants m ON ma.merchant_id = m.merchant_id
		LEFT JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		LEFT JOIN provider_paychannels pp ON t.provider_paychannel_id = pp.ID
		LEFT JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		LEFT JOIN payment_methods pm ON ppm.payment_method_id = pm.ID
		LEFT JOIN providers p ON ppm.provider_id = p.ID
		LEFT JOIN reason_lists rl ON mcf.reason_id = rl.ID
		`

		var conditions []string

		if payload.MerchantId != "" {
			conditions = append(conditions, fmt.Sprintf("m.merchant_id = '%v'", payload.MerchantId))
		}

		if payload.MinDate != "" {
			conditions = append(conditions, fmt.Sprintf("mcf.created_at >= '%v'", payload.MinDate))
		}

		if payload.MaxDate != "" {
			conditions = append(conditions, fmt.Sprintf("mcf.created_at <= '%v'", payload.MaxDate))
		}

		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}

		query += " ORDER BY mcf.created_at DESC"
	}

	if payload.ExportType == constant.ExportTypeIn {
		query = `
		SELECT
			mcf.ID,
			mcf.payment_id,
			m.merchant_id,
			m.merchant_name,
			mcf.amount,
			rl.reason_name,
			pm.name AS payment_method,
			mp.fee,
			mp.fee_type,
			mcf.temp_balance,
			p.provider_name,
			pp.paychannel_name,
			pp.fee AS provider_fee,
			pp.fee_type AS provider_fee_type,
			mcf.status,
			mcf.notes,
			mcf.capital_type,
			mcf.reverse_from,
			mcf.created_at
		FROM
			merchant_capital_flows mcf
		LEFT JOIN transactions t ON t.payment_id = mcf.payment_id
		LEFT JOIN merchant_accounts ma ON ma.ID = mcf.merchant_account_id
		LEFT JOIN merchants m ON ma.merchant_id = m.merchant_id
		LEFT JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		LEFT JOIN provider_paychannels pp ON t.provider_paychannel_id = pp.ID
		LEFT JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		LEFT JOIN payment_methods pm ON ppm.payment_method_id = pm.ID
		LEFT JOIN providers p ON ppm.provider_id = p.ID
		LEFT JOIN reason_lists rl ON mcf.reason_id = rl.ID
		WHERE
			mcf.reason_id IN (5, 7)
			AND pm.name NOT IN ('Disbursement')
		`

		var conditions []string

		if payload.MerchantId != "" {
			conditions = append(conditions, fmt.Sprintf("m.merchant_id = '%v'", payload.MerchantId))
		}

		if payload.MinDate != "" {
			conditions = append(conditions, fmt.Sprintf("mcf.created_at >= '%v'", payload.MinDate))
		}

		if payload.MaxDate != "" {
			conditions = append(conditions, fmt.Sprintf("mcf.created_at <= '%v'", payload.MaxDate))
		}

		if len(conditions) > 0 {
			query += " AND " + strings.Join(conditions, " AND ")
		}

		query += " ORDER BY mcf.created_at DESC"
	}

	if payload.ExportType == constant.ExportTypeOut {
		query = `
		SELECT
			mcf.ID,
			mcf.payment_id,
			m.merchant_id,
			m.merchant_name,
			mcf.amount,
			rl.reason_name,
			pm.name AS payment_method,
			mp.fee,
			mp.fee_type,
			mcf.temp_balance,
			p.provider_name,
			pp.paychannel_name,
			pp.fee AS provider_fee,
			pp.fee_type AS provider_fee_type,
			mcf.status,
			mcf.notes,
			mcf.capital_type,
			mcf.reverse_from,
			mcf.created_at
		FROM
			merchant_capital_flows mcf
		LEFT JOIN transactions t ON t.payment_id = mcf.payment_id
		LEFT JOIN merchant_accounts ma ON ma.ID = mcf.merchant_account_id
		LEFT JOIN merchants m ON ma.merchant_id = m.merchant_id
		LEFT JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		LEFT JOIN provider_paychannels pp ON t.provider_paychannel_id = pp.ID
		LEFT JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		LEFT JOIN payment_methods pm ON ppm.payment_method_id = pm.ID
		LEFT JOIN providers p ON ppm.provider_id = p.ID
		LEFT JOIN reason_lists rl ON mcf.reason_id = rl.ID
		WHERE
			mcf.reason_id IN (6, 7)
			AND pm.ID NOT IN (1, 2, 3)
		`

		var conditions []string

		if payload.MerchantId != "" {
			conditions = append(conditions, fmt.Sprintf("m.merchant_id = '%v'", payload.MerchantId))
		}

		if payload.MinDate != "" {
			conditions = append(conditions, fmt.Sprintf("mcf.created_at >= '%v'", payload.MinDate))
		}

		if payload.MaxDate != "" {
			conditions = append(conditions, fmt.Sprintf("mcf.created_at <= '%v'", payload.MaxDate))
		}

		if len(conditions) > 0 {
			query += " AND " + strings.Join(conditions, " AND ")
		}

		query += " ORDER BY mcf.created_at DESC"
	}

	err := tr.db.Select(&listData, query)
	if err != nil {
		return listData, err
	}

	return listData, nil
}

func (tr *TransactionsWrites) UpdateReportStoragesByFileName(publicUrl string, fileName string, status string) error {
	query := `
	UPDATE report_storages
	SET report_url = $1,
		status = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE file_name = $3;
	`

	_, err := tr.db.Exec(query, publicUrl, status, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TransactionsWrites) CreateListReportStoragesRepo(payload dto.CreateReportStorageDto) (int, error) {
	var idReport int

	query := `
	INSERT INTO report_storages (
    merchant_id,
    period,
    export_type,
    status,
    report_url,
    created_by_user,
    file_name,
    created_at,
    updated_at
	) VALUES (
		$1, -- merchant_id
		$2, -- period
		$3, -- export_type
		$4, -- status
		$5, -- report_url
		$6, -- created_by_user
		$7, -- file_name 
		CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta',
		CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	)
	RETURNING id
	`

	row := tr.db.QueryRow(query, payload.MerchantId, payload.Period, payload.ExportType, payload.Status, payload.ReportUrl, payload.CreatedByUser, payload.FileName)
	err := row.Scan(&idReport)
	if err != nil || idReport == 0 {
		return idReport, err
	}

	return idReport, nil
}
