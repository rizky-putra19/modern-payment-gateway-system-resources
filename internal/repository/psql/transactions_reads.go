package psql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/jmoiron/sqlx"
)

type TransactionsReads struct {
	db *sqlx.DB
}

func NewTransactionsReads(db *sqlx.DB) *TransactionsReads {
	return &TransactionsReads{
		db: db,
	}
}

func (tr *TransactionsReads) GetCapitalFlowsWitPaymentId(paymentId string) ([]entity.CapitalFlows, error) {
	var transactionCapitalFlows []entity.CapitalFlows

	query := `
	SELECT
		mcf.ID,
		mcf.payment_id,
		mcf.merchant_account_id,
		mcf.temp_balance,
		mcf.status,
		mcf.notes,
		mcf.amount,
		mcf.capital_type,
		mcf.created_by,
		ma.merchant_id AS merchant_id,
		m.merchant_name AS merchant_name,
		r.ID AS reason_id,
		r.reason_name,
		r.reason_description,
		mcf.created_at
	FROM
		merchant_capital_flows mcf
		JOIN merchant_accounts ma ON mcf.merchant_account_id = ma.ID
		JOIN reason_lists r ON mcf.reason_id = r.ID
		JOIN merchants m ON ma.merchant_id = m.merchant_id
	WHERE
		mcf.payment_id = $1
	ORDER BY mcf.created_at DESC;
	`

	err := tr.db.Select(&transactionCapitalFlows, query, paymentId)
	if err != nil {
		return transactionCapitalFlows, err
	}

	return transactionCapitalFlows, nil
}

func (tr *TransactionsReads) GetStatusChangeLogData(paymentId string) ([]entity.TransactionStatusLogs, error) {
	var statusLogs []entity.TransactionStatusLogs

	query := `
	SELECT *
	FROM transaction_status_logs tsl
	WHERE tsl.payment_id = $1
	ORDER BY tsl.created_at DESC;
	`

	err := tr.db.Select(&statusLogs, query, paymentId)
	if err != nil {
		return statusLogs, err
	}

	return statusLogs, nil
}

func (tr *TransactionsReads) GetPaymentDetailPrvConfirmDetail(paymentId string) ([]entity.ProviderConfirmDetail, error) {
	var detailConfirmData []entity.ProviderConfirmDetail

	query := `
	SELECT *
	FROM provider_transaction_confirmation_details ptcd
	WHERE ptcd.payment_id = $1
	ORDER BY ptcd.created_at DESC;
	`

	err := tr.db.Select(&detailConfirmData, query, paymentId)
	if err != nil {
		return detailConfirmData, err
	}

	return detailConfirmData, nil
}

func (tr *TransactionsReads) GetPaymentDetailAccountInformation(paymentId string) ([]entity.AccountData, error) {
	var accountData []entity.AccountData

	query := `
	SELECT *
	FROM account_informations ai
	WHERE ai.payment_id = $1;
	`

	err := tr.db.Select(&accountData, query, paymentId)
	if err != nil {
		return accountData, err
	}

	return accountData, nil
}

func (tr *TransactionsReads) GetAccountInformationByPaymentIdAccountType(paymentId string, accountType string) (entity.AccountData, error) {
	var accountData entity.AccountData

	query := `
	SELECT *
	FROM account_informations ai
	WHERE ai.payment_id = $1
	AND ai.account_type = $2;
	`

	err := tr.db.Get(&accountData, query, paymentId, accountType)
	if err != nil {
		return accountData, err
	}

	return accountData, nil
}

func (tr *TransactionsReads) GetPaymentDetailProviderMerchant(paymentId string) (entity.PaymentDetailMerchantProvider, error) {
	var detailData entity.PaymentDetailMerchantProvider

	query := `
    SELECT
		t.ID AS transaction_id,
		t.payment_id,
		t.merchant_reference_number,
		t.provider_reference_number,
		t.transaction_amount,
		t.bank_code,
		t.status,
		t.client_ip_address,
		t.merchant_callback_url,
		t.request_method,
		pm.name AS payment_method_name,
		pm.pay_type AS pay_type,
		t.created_at AS transaction_created_at,
		t.updated_at AS transaction_updated_at,
		m.merchant_id,
		m.merchant_name,
		mp.merchant_payment_method_id,
		mp.segment,
		mp.fee AS merchant_fee,
		mp.fee_type AS merchant_fee_type,
		mp.status AS merchant_status,
		mp.min_transaction AS merchant_min_transaction,
		mp.max_transaction AS merchant_max_transaction,
		mp.max_daily_transaction AS merchant_max_daily_transaction,
		mp.merchant_paychannel_code,
		mp.created_at AS merchant_created_at,
		mp.updated_at AS merchant_updated_at,
		p.provider_name,
		pp.provider_payment_method_id,
		pp.bank_code AS provider_bank_code,
		pp.paychannel_name,
		pp.fee AS provider_fee,
		pp.fee_type AS provider_fee_type,
		pp.status AS provider_status,
		pp.min_transaction AS provider_min_transaction,
		pp.max_transaction AS provider_max_transaction,
		pp.max_daily_transaction AS provider_max_daily_transaction,
		pp.interface_setting,
		pp.created_at AS provider_created_at,
		pp.updated_at AS provider_updated_at
	FROM
		transactions t
		LEFT JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		LEFT JOIN provider_paychannels pp ON t.provider_paychannel_id = pp.ID
		JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		JOIN providers p ON ppm.provider_id = p.ID
	WHERE
		t.merchant_paychannel_id IS NOT NULL
		AND t.provider_paychannel_id IS NOT NULL
		AND t.payment_id = $1;
    `

	err := tr.db.Get(&detailData, query, paymentId)
	if err != nil {
		return detailData, err
	}

	return detailData, nil
}

func (tr *TransactionsReads) GetTransactionList(params dto.QueryParams) ([]entity.Transaction, dto.PaginatedResponse, error) {
	var transactions []entity.Transaction
	var pagination dto.PaginatedResponse
	pageInt := converter.ToInt(params.Page)
	pageSizeInt := converter.ToInt(params.PageSize)

	// Set default values for pagination if not provided
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 {
		pageSizeInt = 50 // Default page size
	}

	query, err := buildTransactionQuery(params)
	if err != nil {
		return nil, pagination, err
	}

	err = tr.db.Select(&transactions, query)
	if err != nil {
		return nil, pagination, err
	}

	// Get the total number of items for the given filters
	countQuery := buildCountQuery(params)
	var totalItems int
	err = tr.db.Get(&totalItems, countQuery)
	if err != nil {
		return nil, pagination, err
	}

	// Calculate total pages
	totalPages := (totalItems + pageSizeInt - 1) / pageSizeInt
	pagination = dto.PaginatedResponse{
		CurrentPage: pageInt,
		PageSize:    pageSizeInt,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNextPage: pageInt < totalPages,
		HasPrevPage: pageInt > 1,
	}

	return transactions, pagination, nil
}

func (tr *TransactionsReads) GetTransactionInListRepo(params dto.QueryParams) ([]entity.MerchantTransactionList, dto.PaginatedResponse, error) {
	var transactionInList []entity.MerchantTransactionList
	var pagination dto.PaginatedResponse
	pageInt := converter.ToInt(params.Page)
	pageSizeInt := converter.ToInt(params.PageSize)

	if pageInt < 1 {
		pageInt = 1
	}

	if pageSizeInt < 1 {
		pageSizeInt = 50
	}

	query, err := buildListInMerchantQuery(params)
	if err != nil {
		return nil, pagination, err
	}

	err = tr.db.Select(&transactionInList, query)
	if err != nil {
		return nil, pagination, err
	}

	// Get the total number of items for the given filters
	countQuery := buildCountListInQuery(params)
	var totalItems int
	err = tr.db.Get(&totalItems, countQuery)
	if err != nil {
		return nil, pagination, err
	}

	// Calculate total pages
	totalPages := (totalItems + pageSizeInt - 1) / pageSizeInt
	pagination = dto.PaginatedResponse{
		CurrentPage: pageInt,
		PageSize:    pageSizeInt,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNextPage: pageInt < totalPages,
		HasPrevPage: pageInt > 1,
	}

	return transactionInList, pagination, nil
}

func (tr *TransactionsReads) GetTransactionOutListRepo(params dto.QueryParams) ([]entity.MerchantTransactionList, dto.PaginatedResponse, error) {
	var transactionInList []entity.MerchantTransactionList
	var pagination dto.PaginatedResponse
	pageInt := converter.ToInt(params.Page)
	pageSizeInt := converter.ToInt(params.PageSize)

	if pageInt < 1 {
		pageInt = 1
	}

	if pageSizeInt < 1 {
		pageSizeInt = 50
	}

	query := buildListOutMerchantQuery(params)

	err := tr.db.Select(&transactionInList, query)
	if err != nil {
		return nil, pagination, err
	}

	// Get the total number of items for the given filters
	countQuery := buildCountListOutMerchantQuery(params)
	var totalItems int
	err = tr.db.Get(&totalItems, countQuery)
	if err != nil {
		return nil, pagination, err
	}

	// Calculate total pages
	totalPages := (totalItems + pageSizeInt - 1) / pageSizeInt
	pagination = dto.PaginatedResponse{
		CurrentPage: pageInt,
		PageSize:    pageSizeInt,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNextPage: pageInt < totalPages,
		HasPrevPage: pageInt > 1,
	}

	return transactionInList, pagination, nil
}

func (tr *TransactionsReads) GetTransactionAnalyticsRepo(params dto.GetMerchantAnalyticsDtoReq) ([]entity.PaymentDetailMerchantProvider, error) {
	var transactionDatas []entity.PaymentDetailMerchantProvider

	query := `
		SELECT
		t.ID AS transaction_id,
		t.payment_id,
		t.merchant_reference_number,
		t.provider_reference_number,
		t.transaction_amount,
		t.bank_code,
		t.status,
		t.client_ip_address,
		t.merchant_callback_url,
		t.request_method,
		pm.name AS payment_method_name,
		pm.pay_type AS pay_type,
		t.created_at AS transaction_created_at,
		t.updated_at AS transaction_updated_at,
		m.merchant_id,
		m.merchant_name,
		mp.merchant_payment_method_id,
		mp.segment,
		mp.fee AS merchant_fee,
		mp.fee_type AS merchant_fee_type,
		mp.status AS merchant_status,
		mp.min_transaction AS merchant_min_transaction,
		mp.max_transaction AS merchant_max_transaction,
		mp.max_daily_transaction AS merchant_max_daily_transaction,
		mp.merchant_paychannel_code,
		mp.created_at AS merchant_created_at,
		mp.updated_at AS merchant_updated_at,
		p.provider_name,
		pp.provider_payment_method_id,
		pp.bank_code AS provider_bank_code,
		pp.paychannel_name,
		pp.fee AS provider_fee,
		pp.fee_type AS provider_fee_type,
		pp.status AS provider_status,
		pp.min_transaction AS provider_min_transaction,
		pp.max_transaction AS provider_max_transaction,
		pp.max_daily_transaction AS provider_max_daily_transaction,
		pp.interface_setting,
		pp.created_at AS provider_created_at,
		pp.updated_at AS provider_updated_at
	FROM
		transactions t
		LEFT JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		LEFT JOIN provider_paychannels pp ON t.provider_paychannel_id = pp.ID
		JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		JOIN providers p ON ppm.provider_id = p.ID
	WHERE
		t.merchant_paychannel_id IS NOT NULL
		AND t.provider_paychannel_id IS NOT NULL
		AND m.merchant_id = $1
	`

	var conditions []string

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", params.MaxDate))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.created_at DESC"

	err := tr.db.Select(&transactionDatas, query, params.MerchantId)
	if err != nil {
		return transactionDatas, err
	}

	return transactionDatas, nil
}

func (tr *TransactionsReads) GetTransactionDataForProviderAnalyticsRepo(payload dto.GetProviderAnalyticsDtoReq) ([]entity.PaymentDetailMerchantProvider, error) {
	var transactionDatas []entity.PaymentDetailMerchantProvider

	query := `
		SELECT
		t.ID AS transaction_id,
		t.payment_id,
		t.merchant_reference_number,
		t.provider_reference_number,
		t.transaction_amount,
		t.bank_code,
		t.status,
		t.client_ip_address,
		t.merchant_callback_url,
		t.request_method,
		pm.name AS payment_method_name,
		pm.pay_type AS pay_type,
		t.created_at AS transaction_created_at,
		t.updated_at AS transaction_updated_at,
		m.merchant_id,
		m.merchant_name,
		mp.merchant_payment_method_id,
		mp.segment,
		mp.fee AS merchant_fee,
		mp.fee_type AS merchant_fee_type,
		mp.status AS merchant_status,
		mp.min_transaction AS merchant_min_transaction,
		mp.max_transaction AS merchant_max_transaction,
		mp.max_daily_transaction AS merchant_max_daily_transaction,
		mp.merchant_paychannel_code,
		mp.created_at AS merchant_created_at,
		mp.updated_at AS merchant_updated_at,
		p.provider_name,
		pp.provider_payment_method_id,
		pp.bank_code AS provider_bank_code,
		pp.paychannel_name,
		pp.fee AS provider_fee,
		pp.fee_type AS provider_fee_type,
		pp.status AS provider_status,
		pp.min_transaction AS provider_min_transaction,
		pp.max_transaction AS provider_max_transaction,
		pp.max_daily_transaction AS provider_max_daily_transaction,
		pp.interface_setting,
		pp.created_at AS provider_created_at,
		pp.updated_at AS provider_updated_at
	FROM
		transactions t
		LEFT JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		LEFT JOIN provider_paychannels pp ON t.provider_paychannel_id = pp.ID
		JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		JOIN providers p ON ppm.provider_id = p.ID
	WHERE
		t.merchant_paychannel_id IS NOT NULL
		AND t.provider_paychannel_id IS NOT NULL
		AND p.ID = $1
	`

	var conditions []string

	if payload.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", payload.MinDate))
	}

	if payload.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", payload.MaxDate))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.created_at DESC"

	err := tr.db.Select(&transactionDatas, query, payload.ProviderId)
	if err != nil {
		return transactionDatas, err
	}

	return transactionDatas, nil
}

func (tr *TransactionsReads) GetTransactionCapitalFlowRepo(params dto.QueryParams) ([]entity.TransactionCapitalFlows, error) {
	var listTransactionCapital []entity.TransactionCapitalFlows

	query := `
	SELECT
		mcf.ID,
		mcf.payment_id,
		mcf.amount,
		mcf.created_at,
		ma.merchant_id,
		mcf.status,
		r.reason_name,
		mcf.temp_balance,
		mcf.capital_type
	FROM
		merchant_capital_flows mcf
		JOIN merchant_accounts ma ON mcf.merchant_account_id = ma.ID
		JOIN reason_lists r ON mcf.reason_id = r.ID
		JOIN merchants m ON ma.merchant_id = m.merchant_id
		JOIN transactions t ON t.payment_id = mcf.payment_id
		JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN payment_methods pm ON pm.ID = mpm.payment_method_id
	WHERE
		mcf.reason_id IN (7, 6, 5)
	`

	var conditions []string

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.created_at <= '%v'", params.MaxDate))
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("mcf.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(mcf.payment_id LIKE '%v')", searchStr))
	}

	if params.MerchantId != "" {
		conditions = append(conditions, fmt.Sprintf("ma.merchant_id = '%v'", params.MerchantId))
	}

	if params.PayType != "" {
		slicePayType := helper.SplitString(params.PayType)
		var payTypeConditions []string
		for _, payType := range slicePayType {
			payTypeConditions = append(payTypeConditions, fmt.Sprintf("pm.pay_type = '%v'", payType))
		}
		conditions = append(conditions, "("+strings.Join(payTypeConditions, " OR ")+")")
	}

	if params.PaymentMethod != "" {
		slicePaymentMethod := helper.SplitString(params.PaymentMethod)
		var paymentMethodConditions []string
		for _, paymentMethod := range slicePaymentMethod {
			paymentMethodConditions = append(paymentMethodConditions, fmt.Sprintf("pm.name = '%v'", paymentMethod))
		}
		conditions = append(conditions, "("+strings.Join(paymentMethodConditions, " OR ")+")")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY mcf.created_at DESC"

	err := tr.db.Select(&listTransactionCapital, query)
	if err != nil && err != sql.ErrNoRows {
		return listTransactionCapital, err
	}

	return listTransactionCapital, nil
}

func (tr *TransactionsReads) GetTransactionListByMerchantPaychannelRepo(payload dto.GetMerchantAnalyticsDtoReq) ([]entity.PaymentDetailMerchantProvider, error) {
	var transactionData []entity.PaymentDetailMerchantProvider

	query := `
	SELECT
		t.ID AS transaction_id,
		t.payment_id,
		t.merchant_reference_number,
		t.provider_reference_number,
		t.transaction_amount,
		t.bank_code,
		t.status,
		t.client_ip_address,
		t.merchant_callback_url,
		t.request_method,
		pm.name AS payment_method_name,
		pm.pay_type AS pay_type,
		t.created_at AS transaction_created_at,
		t.updated_at AS transaction_updated_at,
		m.merchant_id,
		m.merchant_name,
		mp.merchant_payment_method_id,
		mp.segment,
		mp.fee AS merchant_fee,
		mp.fee_type AS merchant_fee_type,
		mp.status AS merchant_status,
		mp.min_transaction AS merchant_min_transaction,
		mp.max_transaction AS merchant_max_transaction,
		mp.max_daily_transaction AS merchant_max_daily_transaction,
		mp.merchant_paychannel_code,
		mp.created_at AS merchant_created_at,
		mp.updated_at AS merchant_updated_at,
		p.provider_name,
		pp.provider_payment_method_id,
		pp.bank_code AS provider_bank_code,
		pp.paychannel_name,
		pp.fee AS provider_fee,
		pp.fee_type AS provider_fee_type,
		pp.status AS provider_status,
		pp.min_transaction AS provider_min_transaction,
		pp.max_transaction AS provider_max_transaction,
		pp.max_daily_transaction AS provider_max_daily_transaction,
		pp.interface_setting,
		pp.created_at AS provider_created_at,
		pp.updated_at AS provider_updated_at
	FROM
		transactions t
		LEFT JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		LEFT JOIN provider_paychannels pp ON t.provider_paychannel_id = pp.ID
		JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		JOIN providers p ON ppm.provider_id = p.ID
	WHERE
		t.merchant_paychannel_id IS NOT NULL
		AND t.provider_paychannel_id IS NOT NULL
		AND t.merchant_paychannel_id = $1
	`

	var conditions []string

	if payload.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", payload.MinDate))
	}

	if payload.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", payload.MaxDate))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.created_at DESC"

	err := tr.db.Select(&transactionData, query, payload.MerchantPaychannelId)
	if err != nil {
		return transactionData, err
	}

	return transactionData, nil
}

func (tr *TransactionsReads) GetListMerchantReportRepo(params dto.GetListMerchantExportFilter, merchantId string) ([]entity.ReportStoragesEntity, error) {
	var listMerchantExport []entity.ReportStoragesEntity

	query := `
	SELECT
		rs.ID,
		rs.created_at,
		rs.merchant_id,
		m.merchant_name,
		m.currency,
		rs.export_type,
		rs.period,
		rs.status,
		rs.report_url,
		rs.created_by_user
	FROM
		report_storages rs
		JOIN merchants m ON m.merchant_id = rs.merchant_id
	`

	var conditions []string

	conditions = append(conditions, fmt.Sprintf("rs.merchant_id = '%v'", merchantId))
	conditions = append(conditions, fmt.Sprintf("rs.created_by_user = '%v'", "USER_MERCHANT"))

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("rs.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("rs.created_at <= '%v'", params.MaxDate))
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(rs.merchant_id LIKE '%v')", searchStr))
	}

	if params.ExportStatus != "" {
		sliceStatus := helper.SplitString(params.ExportStatus)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("rs.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.ExportType != "" {
		sliceExportType := helper.SplitString(params.ExportType)
		var exportTypeConditions []string
		for _, exportType := range sliceExportType {
			exportTypeConditions = append(exportTypeConditions, fmt.Sprintf("rs.export_type = '%v'", exportType))
		}
		conditions = append(conditions, "("+strings.Join(exportTypeConditions, " OR ")+")")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY rs.created_at DESC"

	err := tr.db.Select(&listMerchantExport, query)
	if err != nil {
		return listMerchantExport, err
	}

	return listMerchantExport, nil
}

func (tr *TransactionsReads) GetListMerchantExportRepo(params dto.GetListMerchantExportFilter) ([]entity.ReportStoragesEntity, error) {
	var listMerchantExport []entity.ReportStoragesEntity

	query := `
	SELECT
		rs.ID,
		rs.created_at,
		rs.merchant_id,
		m.merchant_name,
		m.currency,
		rs.export_type,
		rs.period,
		rs.status,
		rs.report_url,
		rs.created_by_user
	FROM
		report_storages rs
		JOIN merchants m ON m.merchant_id = rs.merchant_id
	`

	if params.Merchants == constant.InternalExport {
		query = `
		SELECT
			rs.ID,
			rs.created_at,
			rs.export_type,
			rs.period,
			rs.status,
			rs.report_url,
			rs.created_by_user
		FROM
			report_storages rs
		`
	}

	var conditions []string

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("rs.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("rs.created_at <= '%v'", params.MaxDate))
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(rs.merchant_id LIKE '%v')", searchStr))
	}

	if params.Merchants != "" {
		sliceMerchantName := helper.SplitString(params.Merchants)
		var merchantNameConditions []string
		for _, merchantName := range sliceMerchantName {
			merchantNameConditions = append(merchantNameConditions, fmt.Sprintf("rs.merchant_id = '%v'", merchantName))
		}
		conditions = append(conditions, "("+strings.Join(merchantNameConditions, " OR ")+")")
	}

	if params.ExportStatus != "" {
		sliceStatus := helper.SplitString(params.ExportStatus)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("rs.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.ExportType != "" {
		sliceExportType := helper.SplitString(params.ExportType)
		var exportTypeConditions []string
		for _, exportType := range sliceExportType {
			exportTypeConditions = append(exportTypeConditions, fmt.Sprintf("rs.export_type = '%v'", exportType))
		}
		conditions = append(conditions, "("+strings.Join(exportTypeConditions, " OR ")+")")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY rs.created_at DESC"

	err := tr.db.Select(&listMerchantExport, query)
	if err != nil {
		return listMerchantExport, err
	}

	return listMerchantExport, nil
}

func (tr *TransactionsReads) GetBankDataDetailRepo(bankName string) (entity.BankListDto, error) {
	var bankData entity.BankListDto

	query := `
	SELECT
		bl.ID,
		bl.bank_name,
		bl.bank_code,
		bl.created_at
	FROM
		bank_lists bl
	WHERE bl.bank_name = $1;
	`

	err := tr.db.Get(&bankData, query, bankName)
	if err != nil {
		return bankData, err
	}

	return bankData, nil
}

func (tr *TransactionsReads) GetBankDataDetailByBankCodeRepo(bankCode string) (entity.BankListDto, error) {
	var bankData entity.BankListDto

	query := `
	SELECT
		bl.ID,
		bl.bank_name,
		bl.bank_code,
		bl.created_at
	FROM
		bank_lists bl
	WHERE bl.bank_code = $1;
	`

	err := tr.db.Get(&bankData, query, bankCode)
	if err != nil && err != sql.ErrNoRows {
		return bankData, err
	}

	return bankData, nil
}

func (tr *TransactionsReads) GetTransactionDataByProviderChannelRepo(payload dto.GetProviderAnalyticsDtoReq) ([]entity.PaymentDetailMerchantProvider, error) {
	var transactionData []entity.PaymentDetailMerchantProvider

	query := `
	SELECT
		t.ID AS transaction_id,
		t.payment_id,
		t.merchant_reference_number,
		t.provider_reference_number,
		t.transaction_amount,
		t.bank_code,
		t.status,
		t.client_ip_address,
		t.merchant_callback_url,
		t.request_method,
		pm.name AS payment_method_name,
		pm.pay_type AS pay_type,
		t.created_at AS transaction_created_at,
		t.updated_at AS transaction_updated_at,
		m.merchant_id,
		m.merchant_name,
		mp.merchant_payment_method_id,
		mp.segment,
		mp.fee AS merchant_fee,
		mp.fee_type AS merchant_fee_type,
		mp.status AS merchant_status,
		mp.min_transaction AS merchant_min_transaction,
		mp.max_transaction AS merchant_max_transaction,
		mp.max_daily_transaction AS merchant_max_daily_transaction,
		mp.merchant_paychannel_code,
		mp.created_at AS merchant_created_at,
		mp.updated_at AS merchant_updated_at,
		p.provider_name,
		pp.provider_payment_method_id,
		pp.bank_code AS provider_bank_code,
		pp.paychannel_name,
		pp.fee AS provider_fee,
		pp.fee_type AS provider_fee_type,
		pp.status AS provider_status,
		pp.min_transaction AS provider_min_transaction,
		pp.max_transaction AS provider_max_transaction,
		pp.max_daily_transaction AS provider_max_daily_transaction,
		pp.interface_setting,
		pp.created_at AS provider_created_at,
		pp.updated_at AS provider_updated_at
	FROM
		transactions t
		LEFT JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		LEFT JOIN provider_paychannels pp ON t.provider_paychannel_id = pp.ID
		JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		JOIN providers p ON ppm.provider_id = p.ID
	WHERE
		t.merchant_paychannel_id IS NOT NULL
		AND t.provider_paychannel_id IS NOT NULL
		AND t.provider_paychannel_id = $1
	`

	var conditions []string

	if payload.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", payload.MinDate))
	}

	if payload.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", payload.MaxDate))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.created_at DESC"

	err := tr.db.Select(&transactionData, query, payload.ProviderChannelId)
	if err != nil {
		return transactionData, err
	}

	return transactionData, nil
}

func buildTransactionQuery(params dto.QueryParams) (string, error) {
	pageInt := converter.ToInt(params.Page)
	pageSizeInt := converter.ToInt(params.PageSize)

	query := `
	SELECT
		t.id,
		t.payment_id,
		t.merchant_reference_number,
		t.transaction_amount,
		t.bank_code,
		t.status AS transaction_status,
		t.created_at AS transaction_created_at,
		t.updated_at AS transaction_updated_at,
		
		m.merchant_id,
		m.merchant_name,
		
		pm.name AS payment_method_name,
		pm.pay_type AS payment_method_type,

		ppch.paychannel_name AS provider_paychannel_name,
		bl.bank_name AS bank_name

	FROM
		transactions t
		JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN provider_paychannels ppch ON t.provider_paychannel_id = ppch.ID
		JOIN provider_payment_methods ppm ON ppch.provider_payment_method_id = ppm.ID 
		JOIN provider_payment_method_bank_lists ppmbl ON ppm.ID = ppmbl.provider_payment_method_id 
		JOIN bank_lists bl ON ppmbl.bank_list_id = bl.id 
		JOIN providers pp ON ppm.provider_id = pp.ID
	WHERE
		bl.bank_code = t.bank_code
`

	var conditions []string

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", params.MaxDate))
	}

	if params.AmountMin != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_amount >= '%v'", params.AmountMin))
	}

	if params.AmountMax != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_amount <= '%v'", params.AmountMax))
	}

	if params.PaymentMethod != "" {
		slicePaymentMethod := helper.SplitString(params.PaymentMethod)
		var paymentMethodConditions []string
		for _, paymentMethod := range slicePaymentMethod {
			paymentMethodConditions = append(paymentMethodConditions, fmt.Sprintf("pm.name = '%v'", paymentMethod))
		}
		conditions = append(conditions, "("+strings.Join(paymentMethodConditions, " OR ")+")")
	}

	if params.PayType != "" {
		slicePayType := helper.SplitString(params.PayType)
		var payTypeConditions []string
		for _, payType := range slicePayType {
			payTypeConditions = append(payTypeConditions, fmt.Sprintf("pm.pay_type = '%v'", payType))
		}
		conditions = append(conditions, "("+strings.Join(payTypeConditions, " OR ")+")")
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("t.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.RequestMethod != "" {
		sliceRequestMethod := helper.SplitString(params.RequestMethod)
		var requestMethodConditions []string
		for _, requestMethod := range sliceRequestMethod {
			requestMethodConditions = append(requestMethodConditions, fmt.Sprintf("t.request_method = '%v'", requestMethod))
		}
		conditions = append(conditions, "("+strings.Join(requestMethodConditions, " OR ")+")")
	}

	if params.PayChannel != "" {
		slicePayChannel := helper.SplitString(params.PayChannel)
		var payChannelConditions []string
		for _, payChannel := range slicePayChannel {
			payChannelConditions = append(payChannelConditions, fmt.Sprintf("ppch.paychannel_name = '%v'", payChannel))
		}
		conditions = append(conditions, "("+strings.Join(payChannelConditions, " OR ")+")")
	}

	if params.MerchantName != "" {
		sliceMerchantName := helper.SplitString(params.MerchantName)
		var merchantNameConditions []string
		for _, merchantName := range sliceMerchantName {
			merchantNameConditions = append(merchantNameConditions, fmt.Sprintf("m.merchant_name = '%v'", merchantName))
		}
		conditions = append(conditions, "("+strings.Join(merchantNameConditions, " OR ")+")")
	}

	if params.ProviderName != "" {
		sliceProviderName := helper.SplitString(params.ProviderName)
		var providerNameConditions []string
		for _, providerName := range sliceProviderName {
			providerNameConditions = append(providerNameConditions, fmt.Sprintf("pp.provider_name = '%v'", providerName))
		}
		conditions = append(conditions, "("+strings.Join(providerNameConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(t.payment_id LIKE '%v' OR t.merchant_reference_number LIKE '%v' OR t.provider_reference_number LIKE '%v')", searchStr, searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.created_at DESC"

	if pageSizeInt > 0 {
		query += fmt.Sprintf(" LIMIT %d", pageSizeInt)
	}

	if pageInt > 0 {
		offset := (pageInt - 1) * pageSizeInt
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	return query, nil
}

func buildCountListInQuery(params dto.QueryParams) string {
	query := `
	SELECT
		COUNT(*)
	FROM
		transactions t
		JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN provider_paychannels ppch ON t.provider_paychannel_id = ppch.ID
		JOIN provider_payment_methods ppm ON ppch.provider_payment_method_id = ppm.ID 
		JOIN provider_payment_method_bank_lists ppmbl ON ppm.ID = ppmbl.provider_payment_method_id 
		JOIN bank_lists bl ON ppmbl.bank_list_id = bl.id 
		JOIN providers pp ON ppm.provider_id = pp.ID
	WHERE
		bl.bank_code = t.bank_code
		AND pm.ID NOT IN (4)
	`

	var conditions []string

	if params.MerchantId != "" {
		conditions = append(conditions, fmt.Sprintf("m.merchant_id = '%v'", params.MerchantId))
	}

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", params.MaxDate))
	}

	if params.PaymentMethod != "" {
		slicePaymentMethod := helper.SplitString(params.PaymentMethod)
		var paymentMethodConditions []string
		for _, paymentMethod := range slicePaymentMethod {
			paymentMethodConditions = append(paymentMethodConditions, fmt.Sprintf("pm.name = '%v'", paymentMethod))
		}
		conditions = append(conditions, "("+strings.Join(paymentMethodConditions, " OR ")+")")
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("t.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(t.payment_id LIKE '%v' OR t.merchant_reference_number LIKE '%v')", searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	return query
}

func buildCountQuery(params dto.QueryParams) string {
	var query string

	query = `
SELECT
    COUNT(*)
FROM
    transactions t
    JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
    JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
    JOIN merchants m ON mpm.merchant_id = m.ID
    JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
    JOIN provider_paychannels ppch ON t.provider_paychannel_id = ppch.ID
    JOIN provider_payment_methods ppm ON ppch.provider_payment_method_id = ppm.ID 
    JOIN provider_payment_method_bank_lists ppmbl ON ppm.ID = ppmbl.provider_payment_method_id 
    JOIN bank_lists bl ON ppmbl.bank_list_id = bl.id 
    JOIN providers pp ON ppm.provider_id = pp.ID
WHERE
    bl.bank_code = t.bank_code
`

	var conditions []string

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", params.MaxDate))
	}

	if params.AmountMin != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_amount >= '%v'", params.AmountMin))
	}

	if params.AmountMax != "" {
		conditions = append(conditions, fmt.Sprintf("t.transaction_amount <= '%v'", params.AmountMax))
	}

	if params.PaymentMethod != "" {
		slicePaymentMethod := helper.SplitString(params.PaymentMethod)
		var paymentMethodConditions []string
		for _, paymentMethod := range slicePaymentMethod {
			paymentMethodConditions = append(paymentMethodConditions, fmt.Sprintf("pm.name = '%v'", paymentMethod))
		}
		conditions = append(conditions, "("+strings.Join(paymentMethodConditions, " OR ")+")")
	}

	if params.PayType != "" {
		slicePayType := helper.SplitString(params.PayType)
		var payTypeConditions []string
		for _, payType := range slicePayType {
			payTypeConditions = append(payTypeConditions, fmt.Sprintf("pm.pay_type = '%v'", payType))
		}
		conditions = append(conditions, "("+strings.Join(payTypeConditions, " OR ")+")")
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("t.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.RequestMethod != "" {
		sliceRequestMethod := helper.SplitString(params.RequestMethod)
		var requestMethodConditions []string
		for _, requestMethod := range sliceRequestMethod {
			requestMethodConditions = append(requestMethodConditions, fmt.Sprintf("t.request_method = '%v'", requestMethod))
		}
		conditions = append(conditions, "("+strings.Join(requestMethodConditions, " OR ")+")")
	}

	if params.PayChannel != "" {
		slicePayChannel := helper.SplitString(params.PayChannel)
		var payChannelConditions []string
		for _, payChannel := range slicePayChannel {
			payChannelConditions = append(payChannelConditions, fmt.Sprintf("ppch.paychannel_name = '%v'", payChannel))
		}
		conditions = append(conditions, "("+strings.Join(payChannelConditions, " OR ")+")")
	}

	if params.MerchantName != "" {
		sliceMerchantName := helper.SplitString(params.MerchantName)
		var merchantNameConditions []string
		for _, merchantName := range sliceMerchantName {
			merchantNameConditions = append(merchantNameConditions, fmt.Sprintf("m.merchant_name = '%v'", merchantName))
		}
		conditions = append(conditions, "("+strings.Join(merchantNameConditions, " OR ")+")")
	}

	if params.ProviderName != "" {
		sliceProviderName := helper.SplitString(params.ProviderName)
		var providerNameConditions []string
		for _, providerName := range sliceProviderName {
			providerNameConditions = append(providerNameConditions, fmt.Sprintf("pp.provider_name = '%v'", providerName))
		}
		conditions = append(conditions, "("+strings.Join(providerNameConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(t.payment_id LIKE '%v' OR t.merchant_reference_number LIKE '%v' OR t.provider_reference_number LIKE '%v')", searchStr, searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	return query
}

func buildListMerchantCallbackQuery(params dto.QueryParamsMerchantCallback) (string, error) {
	pageInt := converter.ToInt(params.Page)
	pageSizeInt := converter.ToInt(params.PageSize)

	query := `
	WITH first_last_callbacks AS (
	    SELECT 
		mc.*,
		MIN(mc.created_at) OVER (PARTITION BY mc.payment_id) AS first_created_at,
		MAX(mc.created_at) OVER (PARTITION BY mc.payment_id) AS last_created_at
	    FROM 
		merchant_callbacks mc
	),
	latest_callbacks AS (
	    SELECT 
		flc.*
	    FROM 
		first_last_callbacks flc
	    WHERE 
		flc.created_at = flc.last_created_at
	)

	SELECT 
	    lc.ID,
	    lc.payment_id,
	    lc.callback_status,
	    lc.payment_status_in_callback,
	    lc.callback_result,
	    lc.created_at AS latest_created_at,
	    lc.first_created_at,
	    lc.last_created_at,
	    t.merchant_reference_number,
	    m.merchant_name,
		m.merchant_id,
	    pm.pay_type
	FROM 
	    latest_callbacks lc
	    JOIN transactions t ON lc.payment_id = t.payment_id
	    JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
	    JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
	    JOIN merchants m ON mpm.merchant_id = m.ID
	    JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
	`

	var conditions []string

	if params.MerchantName != "" {
		merchantNames := helper.SplitString(params.MerchantName)
		quotedMerchantNames := make([]string, len(merchantNames))
		for i, name := range merchantNames {
			quotedMerchantNames[i] = fmt.Sprintf("'%v'", name)
		}
		conditions = append(conditions, fmt.Sprintf("m.merchant_name IN (%v)", strings.Join(quotedMerchantNames, ", ")))
	}

	if params.PayType != "" {
		payTypes := helper.SplitString(params.PayType)
		quotedPayTypes := make([]string, len(payTypes))
		for i, payType := range payTypes {
			quotedPayTypes[i] = fmt.Sprintf("'%v'", payType)
		}
		conditions = append(conditions, fmt.Sprintf("pm.pay_type IN (%v)", strings.Join(quotedPayTypes, ", ")))
	}

	if params.CallbackStatus != "" {
		callbackStatuses := helper.SplitString(params.CallbackStatus)
		quotedCallbackStatuses := make([]string, len(callbackStatuses))
		for i, status := range callbackStatuses {
			quotedCallbackStatuses[i] = fmt.Sprintf("'%v'", status)
		}
		conditions = append(conditions, fmt.Sprintf("lc.callback_status IN (%v)", strings.Join(quotedCallbackStatuses, ", ")))
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(t.payment_id LIKE '%v' OR t.merchant_reference_number LIKE '%v' OR t.provider_reference_number LIKE '%v')", searchStr, searchStr, searchStr))
	}

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("lc.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("lc.created_at <= '%v'", params.MaxDate))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
	GROUP BY
	    lc.ID,
	    lc.payment_id,
	    lc.callback_status,
	    lc.payment_status_in_callback,
	    lc.callback_result,
	    lc.created_at,
	    lc.first_created_at,
	    lc.last_created_at,
	    t.merchant_reference_number,
	    m.merchant_name,
		m.merchant_id,
	    pm.pay_type
	ORDER BY lc.created_at DESC
	`

	if pageSizeInt > 0 {
		query += fmt.Sprintf(" LIMIT %d", pageSizeInt)
	}

	if pageInt > 0 {
		offset := (pageInt - 1) * pageSizeInt
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	return query, nil
}

func buildCountQueryListMerchantCallback(params dto.QueryParamsMerchantCallback) string {
	var query string

	query = `
	WITH first_last_callbacks AS (
		SELECT 
			mc.*,
			MIN(mc.created_at) OVER (PARTITION BY mc.payment_id) AS first_created_at,
			MAX(mc.created_at) OVER (PARTITION BY mc.payment_id) AS last_created_at
		FROM 
			merchant_callbacks mc
	),
	latest_callbacks AS (
		SELECT 
			flc.*
		FROM 
			first_last_callbacks flc
		WHERE 
			flc.created_at = flc.last_created_at
	)
	SELECT
		COUNT(*)
	FROM
		latest_callbacks lc
		JOIN transactions t ON lc.payment_id = t.payment_id
		JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
	`

	var conditions []string

	if params.MerchantName != "" {
		merchantNames := helper.SplitString(params.MerchantName)
		quotedMerchantNames := make([]string, len(merchantNames))
		for i, name := range merchantNames {
			quotedMerchantNames[i] = fmt.Sprintf("'%v'", name)
		}
		conditions = append(conditions, fmt.Sprintf("m.merchant_name IN (%v)", strings.Join(quotedMerchantNames, ", ")))
	}

	if params.PayType != "" {
		payTypes := helper.SplitString(params.PayType)
		quotedPayTypes := make([]string, len(payTypes))
		for i, payType := range payTypes {
			quotedPayTypes[i] = fmt.Sprintf("'%v'", payType)
		}
		conditions = append(conditions, fmt.Sprintf("pm.pay_type IN (%v)", strings.Join(quotedPayTypes, ", ")))
	}

	if params.CallbackStatus != "" {
		callbackStatuses := helper.SplitString(params.CallbackStatus)
		quotedCallbackStatuses := make([]string, len(callbackStatuses))
		for i, status := range callbackStatuses {
			quotedCallbackStatuses[i] = fmt.Sprintf("'%v'", status)
		}
		conditions = append(conditions, fmt.Sprintf("lc.callback_status IN (%v)", strings.Join(quotedCallbackStatuses, ", ")))
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(t.payment_id LIKE '%v' OR t.merchant_reference_number LIKE '%v' OR t.provider_reference_number LIKE '%v')", searchStr, searchStr, searchStr))
	}

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("lc.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("lc.created_at <= '%v'", params.MaxDate))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query
}

func buildListInMerchantQuery(params dto.QueryParams) (string, error) {
	pageInt := converter.ToInt(params.Page)
	pageSizeInt := converter.ToInt(params.PageSize)

	query := `
	SELECT
		t.id,
		t.payment_id,
		m.merchant_id,
		pm.name AS payment_method_name,
		t.transaction_amount,
		t.created_at AS transaction_created_at,
		t.updated_at AS transaction_updated_at,
		t.status AS transaction_status
	FROM
		transactions t
		JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN provider_paychannels ppch ON t.provider_paychannel_id = ppch.ID
		JOIN provider_payment_methods ppm ON ppch.provider_payment_method_id = ppm.ID 
		JOIN provider_payment_method_bank_lists ppmbl ON ppm.ID = ppmbl.provider_payment_method_id 
		JOIN bank_lists bl ON ppmbl.bank_list_id = bl.id 
		JOIN providers pp ON ppm.provider_id = pp.ID
	WHERE
		bl.bank_code = t.bank_code
		AND pm.ID NOT IN (4)
	`

	var conditions []string

	if params.MerchantId != "" {
		conditions = append(conditions, fmt.Sprintf("m.merchant_id = '%v'", params.MerchantId))
	}

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", params.MaxDate))
	}

	if params.PaymentMethod != "" {
		slicePaymentMethod := helper.SplitString(params.PaymentMethod)
		var paymentMethodConditions []string
		for _, paymentMethod := range slicePaymentMethod {
			paymentMethodConditions = append(paymentMethodConditions, fmt.Sprintf("pm.name = '%v'", paymentMethod))
		}
		conditions = append(conditions, "("+strings.Join(paymentMethodConditions, " OR ")+")")
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("t.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(t.payment_id LIKE '%v' OR t.merchant_reference_number LIKE '%v')", searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.created_at DESC"

	if pageSizeInt > 0 {
		query += fmt.Sprintf(" LIMIT %d", pageSizeInt)
	}

	if pageInt > 0 {
		offset := (pageInt - 1) * pageSizeInt
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	return query, nil
}

func buildListOutMerchantQuery(params dto.QueryParams) string {
	pageInt := converter.ToInt(params.Page)
	pageSizeInt := converter.ToInt(params.PageSize)

	query := `
	SELECT
		t.id,
		t.payment_id,
		m.merchant_id,
		pm.name AS payment_method_name,
		t.transaction_amount,
		t.created_at AS transaction_created_at,
		t.updated_at AS transaction_updated_at,
		t.status AS transaction_status
	FROM
		transactions t
		JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN provider_paychannels ppch ON t.provider_paychannel_id = ppch.ID
		JOIN provider_payment_methods ppm ON ppch.provider_payment_method_id = ppm.ID 
		JOIN provider_payment_method_bank_lists ppmbl ON ppm.ID = ppmbl.provider_payment_method_id 
		JOIN bank_lists bl ON ppmbl.bank_list_id = bl.id 
		JOIN providers pp ON ppm.provider_id = pp.ID
	WHERE
		bl.bank_code = t.bank_code
		AND pm.ID NOT IN (1,2,3)
	`

	var conditions []string

	if params.MerchantId != "" {
		conditions = append(conditions, fmt.Sprintf("m.merchant_id = '%v'", params.MerchantId))
	}

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", params.MaxDate))
	}

	if params.PaymentMethod != "" {
		slicePaymentMethod := helper.SplitString(params.PaymentMethod)
		var paymentMethodConditions []string
		for _, paymentMethod := range slicePaymentMethod {
			paymentMethodConditions = append(paymentMethodConditions, fmt.Sprintf("pm.name = '%v'", paymentMethod))
		}
		conditions = append(conditions, "("+strings.Join(paymentMethodConditions, " OR ")+")")
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("t.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(t.payment_id LIKE '%v' OR t.merchant_reference_number LIKE '%v')", searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.created_at DESC"

	if pageSizeInt > 0 {
		query += fmt.Sprintf(" LIMIT %d", pageSizeInt)
	}

	if pageInt > 0 {
		offset := (pageInt - 1) * pageSizeInt
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	return query
}

func buildCountListOutMerchantQuery(params dto.QueryParams) string {

	query := `
	SELECT
		COUNT(*)
	FROM
		transactions t
		JOIN merchant_paychannels mp ON t.merchant_paychannel_id = mp.ID
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN merchants m ON mpm.merchant_id = m.ID
		JOIN payment_methods pm ON mpm.payment_method_id = pm.ID
		JOIN provider_paychannels ppch ON t.provider_paychannel_id = ppch.ID
		JOIN provider_payment_methods ppm ON ppch.provider_payment_method_id = ppm.ID 
		JOIN provider_payment_method_bank_lists ppmbl ON ppm.ID = ppmbl.provider_payment_method_id 
		JOIN bank_lists bl ON ppmbl.bank_list_id = bl.id 
		JOIN providers pp ON ppm.provider_id = pp.ID
	WHERE
		bl.bank_code = t.bank_code
		AND pm.ID NOT IN (1,2,3)
	`

	var conditions []string

	if params.MerchantId != "" {
		conditions = append(conditions, fmt.Sprintf("m.merchant_id = '%v'", params.MerchantId))
	}

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= '%v'", params.MaxDate))
	}

	if params.PaymentMethod != "" {
		slicePaymentMethod := helper.SplitString(params.PaymentMethod)
		var paymentMethodConditions []string
		for _, paymentMethod := range slicePaymentMethod {
			paymentMethodConditions = append(paymentMethodConditions, fmt.Sprintf("pm.name = '%v'", paymentMethod))
		}
		conditions = append(conditions, "("+strings.Join(paymentMethodConditions, " OR ")+")")
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("t.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(t.payment_id LIKE '%v' OR t.merchant_reference_number LIKE '%v')", searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	return query
}
