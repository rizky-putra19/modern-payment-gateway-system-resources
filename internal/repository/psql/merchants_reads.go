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

type MerchantReads struct {
	db *sqlx.DB
}

func NewMerchantReads(db *sqlx.DB) *MerchantReads {
	return &MerchantReads{
		db: db,
	}
}

func (mr *MerchantReads) GetListManualPayment(params dto.QueryParamsManualPayment) ([]entity.ManualPayment, dto.PaginatedResponse, error) {
	var manualPaymentData []entity.ManualPayment
	var pagination dto.PaginatedResponse
	pageSizeInt := converter.ToInt(params.PageSize)
	pageInt := converter.ToInt(params.Page)

	// Set default values for pagination if not provided
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 {
		pageSizeInt = 50 // Default page size
	}

	query := `
	SELECT
		mcf.ID,
		mcf.payment_id,
		ma.merchant_id,
		mcf.reason_id,
		r.reason_name,
		mcf.amount,
		mcf.status,
		mcf.notes,
		mcf.capital_type,
		mcf.created_at
	FROM
		merchant_capital_flows mcf
		JOIN merchant_accounts ma ON mcf.merchant_account_id = ma.ID
		JOIN reason_lists r ON mcf.reason_id = r.ID
		JOIN merchants m ON ma.merchant_id = m.merchant_id
	WHERE
		mcf.reason_id NOT IN (5,6,7)
	`

	var conditions []string

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.created_at <= '%v'", params.MaxDate))
	}

	if params.AmountMax != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.amount <= '%v'", params.AmountMax))
	}

	if params.AmountMin != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.amount >= '%v'", params.AmountMin))
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("mcf.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.ReasonName != "" {
		sliceReasonName := helper.SplitString(params.ReasonName)
		var reasonNameConditions []string
		for _, reasonName := range sliceReasonName {
			reasonNameConditions = append(reasonNameConditions, fmt.Sprintf("r.reason_name = '%v'", reasonName))
		}
		conditions = append(conditions, "("+strings.Join(reasonNameConditions, " OR ")+")")
	}

	if params.MerchantName != "" {
		sliceMerchantName := helper.SplitString(params.MerchantName)
		var merchantNameConditions []string
		for _, merchantName := range sliceMerchantName {
			merchantNameConditions = append(merchantNameConditions, fmt.Sprintf("m.merchant_name = '%v'", merchantName))
		}
		conditions = append(conditions, "("+strings.Join(merchantNameConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(mcf.payment_id LIKE '%v')", searchStr))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY mcf.created_at DESC"

	if pageSizeInt > 0 {
		query += fmt.Sprintf(" LIMIT %d", pageSizeInt)
	}

	if pageInt > 0 {
		offset := (pageInt - 1) * pageSizeInt
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	err := mr.db.Select(&manualPaymentData, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, pagination, err
	}

	countQueryListManualPayment := buildCountQueryManualPayment(params)
	var totalItems int
	err = mr.db.Get(&totalItems, countQueryListManualPayment)
	if err != nil && err != sql.ErrNoRows {
		return manualPaymentData, pagination, err
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

	return manualPaymentData, pagination, nil
}

func (mr *MerchantReads) CheckReverseStatusRepo(paymentId string) ([]entity.ManualPayment, error) {
	var manualPayment []entity.ManualPayment

	query := `
	SELECT
		mcf.ID,
		mcf.payment_id,
		ma.merchant_id,
		mcf.reason_id,
		r.reason_name,
		mcf.amount,
		mcf.status,
		mcf.notes,
		mcf.capital_type,
		mcf.created_at
	FROM
		merchant_capital_flows mcf
		JOIN merchant_accounts ma ON mcf.merchant_account_id = ma.ID
		JOIN reason_lists r ON mcf.reason_id = r.ID
		JOIN merchants m ON ma.merchant_id = m.merchant_id
	WHERE
		mcf.reverse_from = $1;
	`

	err := mr.db.Select(&manualPayment, query, paymentId)
	if err != nil {
		return manualPayment, err
	}

	return manualPayment, nil
}

func (mr *MerchantReads) GetListMerchantCallbackWithFilter(params dto.QueryParamsMerchantCallback) ([]entity.ListMerchantCallback, dto.PaginatedResponse, error) {
	var listMerchantCallback []entity.ListMerchantCallback
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

	query, err := buildListMerchantCallbackQuery(params)
	if err != nil {
		return nil, pagination, err
	}

	err = mr.db.Select(&listMerchantCallback, query)
	if err != nil {
		return nil, pagination, err
	}

	countQueryListCallback := buildCountQueryListMerchantCallback(params)
	var totalItems int
	err = mr.db.Get(&totalItems, countQueryListCallback)
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

	return listMerchantCallback, pagination, nil
}

func (mr *MerchantReads) GetListMerchantCallback(paymentId string) ([]entity.MerchantCallback, error) {
	var callbackData []entity.MerchantCallback

	query := `
	SELECT 
    mc.ID,
    mc.payment_id,
    mc.callback_status,
    mc.payment_status_in_callback,
    mc.callback_result,
    t.merchant_callback_url AS callback_request,
    mc.triggered_by,
    mc.created_at
	FROM 
		merchant_callbacks mc
		LEFT JOIN transactions t ON mc.payment_id = t.payment_id 
	WHERE mc.payment_id = $1
	ORDER BY mc.created_at DESC;
	`

	err := mr.db.Select(&callbackData, query, paymentId)
	if err != nil {
		return callbackData, err
	}

	return callbackData, nil
}

func (mr *MerchantReads) GetMerchantAccountByMerchantId(merchantId string) (entity.MerchantAccount, error) {
	var merchantAccountData entity.MerchantAccount

	query := `
	SELECT *
	FROM merchant_accounts
	WHERE merchant_id = $1;
	`

	err := mr.db.Get(&merchantAccountData, query, merchantId)
	if err != nil {
		return merchantAccountData, err
	}

	return merchantAccountData, nil
}

func (mr *MerchantReads) GetListFilter() (dto.FilterResponseDto, error) {
	var filterResp dto.FilterResponseDto
	var paymentMethods []entity.PaymentMethods
	var providers []entity.Providers
	var merchants []entity.Merchants
	var payChannels []entity.PayChannels
	var reasons []entity.Reasons
	feeType := make([]dto.FeeType, len((constant.FeeType)))
	payType := make([]dto.PayType, len(constant.PayType))

	queryPaymentMethod := `
	SELECT *
	FROM payment_methods pm
	ORDER BY pm.created_at DESC;
	`

	err := mr.db.Select(&paymentMethods, queryPaymentMethod)
	if err != nil {
		return filterResp, err
	}

	queryProviderName := `
	SELECT *
	FROM providers p
	WHERE p.status = 'ACTIVE'
	ORDER BY p.created_at DESC;
	`

	err = mr.db.Select(&providers, queryProviderName)
	if err != nil {
		return filterResp, err
	}

	queryMerchantName := `
	SELECT *
	FROM merchants m
	WHERE m.status = 'ACTIVE'
	ORDER BY m.created_at DESC;
	`

	err = mr.db.Select(&merchants, queryMerchantName)
	if err != nil {
		return filterResp, err
	}

	queryPayChannel := `
	SELECT 
		pc.ID,
		pc.paychannel_name,
		pc.created_at,
		pc.updated_at
	FROM provider_paychannels pc
	WHERE
		pc.status = 'ACTIVE'
	ORDER BY pc.created_at DESC;
	`

	err = mr.db.Select(&payChannels, queryPayChannel)
	if err != nil {
		return filterResp, err
	}

	queryReason := `
	SELECT
		r.ID,
		r.reason_name,
		r.created_at,
		r.updated_at
	FROM reason_lists r
	ORDER BY r.created_at DESC;
	`
	err = mr.db.Select(&reasons, queryReason)
	if err != nil {
		return filterResp, err
	}

	for i := range constant.FeeType {
		id := i + 1
		feeType[i] = dto.FeeType{
			Id:      id,
			FeeType: constant.FeeType[i],
		}
	}

	for i := range constant.PayType {
		id := i + 1
		payType[i] = dto.PayType{
			Id:      id,
			PayType: constant.PayType[i],
		}
	}

	filterResp = dto.FilterResponseDto{
		PaymentMethod: paymentMethods,
		ProviderName:  providers,
		MerchantName:  merchants,
		PayChannel:    payChannels,
		Reason:        reasons,
		FeeType:       feeType,
		PayType:       payType,
	}

	return filterResp, nil
}

func (mr *MerchantReads) GetMerchantDataByMerchantId(merchantId string) (entity.Merchants, error) {
	var merchantData entity.Merchants

	query := `
	SELECT *
	FROM merchants m
	WHERE m.merchant_id = $1;
	`

	err := mr.db.Get(&merchantData, query, merchantId)
	if err != nil {
		return merchantData, err
	}

	return merchantData, nil
}

func (mr *MerchantReads) GetSecretKeyByMerchantIdRepo(merchantId string) (entity.MerchantSecret, error) {
	var secretKey entity.MerchantSecret

	query := `
	SELECT 
		m.merchant_secret
	FROM merchants m
	WHERE m.merchant_id = $1;
	`

	err := mr.db.Get(&secretKey, query, merchantId)
	if err != nil {
		return secretKey, err
	}

	return secretKey, nil
}

func (mr *MerchantReads) GetDetailManualPayment(paymentId string) ([]entity.ManualPayment, error) {
	var manualPayment []entity.ManualPayment

	query := `
	SELECT
		mcf.ID,
		mcf.payment_id,
		ma.merchant_id,
		mcf.reason_id,
		r.reason_name,
		mcf.amount,
		mcf.status,
		mcf.notes,
		mcf.capital_type,
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

	err := mr.db.Select(&manualPayment, query, paymentId)
	if err != nil {
		return manualPayment, err
	}

	return manualPayment, nil
}

func (mr *MerchantReads) GetListMerchantWithFilterRepo(params dto.QueryParams) ([]entity.Merchants, error) {
	var merchantList []entity.Merchants

	query := `
	SELECT *
	FROM merchants m
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

	if params.Status != "" {
		statuses := helper.SplitString(params.Status)
		quotedStatuses := make([]string, len(statuses))
		for i, status := range statuses {
			quotedStatuses[i] = fmt.Sprintf("'%v'", status)
		}
		conditions = append(conditions, fmt.Sprintf("m.status IN (%v)", strings.Join(quotedStatuses, ", ")))
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(m.merchant_id LIKE '%v' OR m.merchant_name LIKE '%v')", searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY m.created_at DESC"

	err := mr.db.Select(&merchantList, query)
	if err != nil {
		return merchantList, err
	}

	return merchantList, nil
}

func (mr *MerchantReads) GetMerchantPaychannelByMerchantId(merchantId string) ([]entity.MerchantPaychannel, error) {
	var merchantPaychannelList []entity.MerchantPaychannel

	query := `
	SELECT
		mpc.ID,
		mpc.merchant_payment_method_id,
		mpc.merchant_paychannel_code,
		pm.name,
		pm.pay_type,
		mpc.fee,
		mpc.fee_type,
		mpc.status,
		mpc.segment,
		mpc.min_transaction,
		mpc.max_transaction,
		mpc.max_daily_transaction,
		mpc.created_at,
		mpc.updated_at
	FROM
		merchant_paychannels mpc
		JOIN merchant_payment_methods mpm ON mpm.ID = mpc.merchant_payment_method_id
		JOIN merchants m ON m.ID = mpm.merchant_id
		JOIN payment_methods pm ON pm.ID = mpm.payment_method_id
	WHERE
		m.merchant_id = $1
	ORDER BY mpc.created_at;
	`

	err := mr.db.Select(&merchantPaychannelList, query, merchantId)
	if err != nil {
		return merchantPaychannelList, err
	}

	return merchantPaychannelList, nil
}

func (mr *MerchantReads) GetMerchantPaychannelDetailById(id int) (entity.MerchantPaychannel, error) {
	var merchantPaychannelData entity.MerchantPaychannel

	query := `
	SELECT
		mpc.ID,
		mpc.merchant_payment_method_id,
		mpc.merchant_paychannel_code,
		pm.name,
		pm.pay_type,
		mpc.fee,
		mpc.fee_type,
		mpc.status,
		mpc.segment,
		mpc.min_transaction,
		mpc.max_transaction,
		mpc.max_daily_transaction,
		mpc.created_at,
		mpc.updated_at
	FROM
		merchant_paychannels mpc
		JOIN merchant_payment_methods mpm ON mpm.ID = mpc.merchant_payment_method_id
		JOIN merchants m ON m.ID = mpm.merchant_id
		JOIN payment_methods pm ON pm.ID = mpm.payment_method_id
	WHERE
		mpc.ID = $1;
	`

	err := mr.db.Get(&merchantPaychannelData, query, id)
	if err != nil {
		return merchantPaychannelData, err
	}

	return merchantPaychannelData, nil
}

func (mr *MerchantReads) GetListMerchantAccountRepo(params dto.QueryParams) ([]entity.ListMerchantAccountDto, error) {
	var listAccount []entity.ListMerchantAccountDto

	query := `
	SELECT
		ma.ID,
		ma.merchant_id,
		m.merchant_name,
		ma.settle_balance,
		ma.not_settle_balance,
		ma.hold_balance,
		ma.pending_transaction_out,
		ma.balance_capital_flow,
		ma.created_at
	FROM
		merchant_accounts ma
		JOIN merchants m ON m.merchant_id = ma.merchant_id
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

	if params.Status != "" {
		statuses := helper.SplitString(params.Status)
		quotedStatuses := make([]string, len(statuses))
		for i, status := range statuses {
			quotedStatuses[i] = fmt.Sprintf("'%v'", status)
		}
		conditions = append(conditions, fmt.Sprintf("m.status IN (%v)", strings.Join(quotedStatuses, ", ")))
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(m.merchant_id LIKE '%v' OR m.merchant_name LIKE '%v')", searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY m.created_at DESC"

	err := mr.db.Select(&listAccount, query)
	if err != nil {
		return listAccount, err
	}

	return listAccount, nil
}

func (mr *MerchantReads) GetListRoutedPaychannelByIdMerchantPaychannelRepo(id int) ([]entity.RoutedPaychanneDto, error) {
	var routedPaychannelList []entity.RoutedPaychanneDto

	query := `
	SELECT
		pr.ID,
		pp.paychannel_name,
		p.provider_name,
		pm.pay_type,
		pm.name,
		pp.fee,
		pp.fee_type,
		pp.min_transaction,
		pp.max_transaction,
		pp.max_daily_transaction,
		pp.status
	FROM
		paychannel_routings pr
		JOIN provider_paychannels pp ON pr.provider_paychannel_id = pp.ID
		JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		JOIN payment_methods pm ON ppm.payment_method_id = pm.ID
		JOIN providers p ON ppm.provider_id = p.ID
	WHERE
		pr.merchant_paychannel_id = $1
	ORDER BY pr.created_at DESC;
	`

	err := mr.db.Select(&routedPaychannelList, query, id)
	if err != nil {
		return routedPaychannelList, err
	}

	return routedPaychannelList, nil
}

func (mr *MerchantReads) GetBankListProviderPaymentMethodRepo(routedChannelName string) ([]entity.BankListDto, error) {
	var bankList []entity.BankListDto

	query := `
	WITH unique_bank_list AS (
		SELECT DISTINCT ON (bl.bank_name)
			bl.ID,
			bl.bank_name,
			bl.bank_code,
			bl.created_at,
			pp.paychannel_name
		FROM
			provider_paychannels pp
			JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
			JOIN provider_payment_method_bank_lists ppmbl ON ppm.ID = ppmbl.provider_payment_method_id
			JOIN bank_lists bl ON ppmbl.bank_list_id = bl.ID
	`

	var conditions []string

	if routedChannelName != "" {
		routedChannelIds := helper.SplitString(routedChannelName)
		quotedRoutedChannelIds := make([]string, len(routedChannelIds))
		for i, name := range routedChannelIds {
			quotedRoutedChannelIds[i] = fmt.Sprintf("'%v'", name)
		}
		conditions = append(conditions, fmt.Sprintf("pp.paychannel_name IN (%v)", strings.Join(quotedRoutedChannelIds, ", ")))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
	)
	SELECT 
		ID,
		bank_name,
		bank_code,
		created_at
	FROM 
		unique_bank_list
	`

	err := mr.db.Select(&bankList, query)
	if err != nil {
		return bankList, err
	}

	return bankList, nil
}

func (mr *MerchantReads) GetBankListFromProviderPaychannelRepo(routedChannelName string) ([]string, error) {
	var bankName []string

	query := `
	WITH unique_bank_list AS (
		SELECT DISTINCT ON (bl.bank_name)
			bl.bank_name
		FROM
			provider_paychannels pp
			JOIN provider_paychannel_bank_lists ppbl ON pp.ID = ppbl.provider_paychannel_id
			JOIN bank_lists bl ON ppbl.bank_list_id = bl.ID
	`

	var conditions []string

	if routedChannelName != "" {
		routedChannelIds := helper.SplitString(routedChannelName)
		quotedRoutedChannelIds := make([]string, len(routedChannelIds))
		for i, name := range routedChannelIds {
			quotedRoutedChannelIds[i] = fmt.Sprintf("'%v'", name)
		}
		conditions = append(conditions, fmt.Sprintf("pp.paychannel_name IN (%v)", strings.Join(quotedRoutedChannelIds, ", ")))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
	)
	SELECT 
		bank_name
	FROM 
		unique_bank_list
	`

	err := mr.db.Select(&bankName, query)
	if err != nil {
		return bankName, err
	}

	return bankName, nil
}

func (mr *MerchantReads) GetBankListForDisbursementRepo(routedChannelName string) ([]entity.BankListDto, error) {
	var bankList []entity.BankListDto

	query := `
	WITH unique_bank_list AS (
		SELECT DISTINCT ON (bl.bank_name)
			bl.ID,
			bl.bank_name,
			bl.bank_code,
			bl.created_at,
			pp.paychannel_name
		FROM
			provider_paychannels pp
			JOIN provider_paychannel_bank_lists ppbl ON pp.ID = ppbl.provider_paychannel_id
			JOIN bank_lists bl ON ppbl.bank_list_id = bl.ID
	`

	var conditions []string

	if routedChannelName != "" {
		routedChannelIds := helper.SplitString(routedChannelName)
		quotedRoutedChannelIds := make([]string, len(routedChannelIds))
		for i, name := range routedChannelIds {
			quotedRoutedChannelIds[i] = fmt.Sprintf("'%v'", name)
		}
		conditions = append(conditions, fmt.Sprintf("pp.paychannel_name IN (%v)", strings.Join(quotedRoutedChannelIds, ", ")))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
	)
	SELECT 
		ID,
		bank_name,
		bank_code,
		created_at
	FROM 
		unique_bank_list
	`

	err := mr.db.Select(&bankList, query)
	if err != nil {
		return bankList, err
	}

	return bankList, nil
}

func (mr *MerchantReads) GetMerchantPaychannelByPaymentMethodId(id int) ([]entity.MerchantPaychannel, error) {
	var merchantPaychannelList []entity.MerchantPaychannel

	query := `
	SELECT
		mpc.ID,
		mpc.merchant_payment_method_id,
		mpc.merchant_paychannel_code,
		pm.name,
		pm.pay_type,
		mpc.fee,
		mpc.fee_type,
		mpc.status,
		mpc.segment,
		mpc.min_transaction,
		mpc.max_transaction,
		mpc.max_daily_transaction,
		mpc.created_at,
		mpc.updated_at
	FROM
		merchant_paychannels mpc
		JOIN merchant_payment_methods mpm ON mpm.ID = mpc.merchant_payment_method_id
		JOIN merchants m ON m.ID = mpm.merchant_id
		JOIN payment_methods pm ON pm.ID = mpm.payment_method_id
	WHERE
		mpc.merchant_payment_method_id = $1
	ORDER BY mpc.created_at;
	`

	err := mr.db.Select(&merchantPaychannelList, query, id)
	if err != nil {
		return merchantPaychannelList, err
	}

	return merchantPaychannelList, nil
}

func (mr *MerchantReads) GetAggregatedPaychannelByIdRepo(id int) (entity.AggregatedPaychannelEntity, error) {
	var aggregatedData entity.AggregatedPaychannelEntity

	query := `
	SELECT
		mpc.ID,
		mpc.merchant_paychannel_code,
		pm.name,
		m.merchant_name
	FROM
		merchant_paychannels mpc
		JOIN merchant_payment_methods mpm ON mpm.ID = mpc.merchant_payment_method_id
		JOIN merchants m ON m.ID = mpm.merchant_id
		JOIN payment_methods pm ON pm.ID = mpm.payment_method_id
	WHERE
		mpc.ID = $1;
	`

	err := mr.db.Get(&aggregatedData, query, id)
	if err != nil {
		return aggregatedData, err
	}

	return aggregatedData, nil
}

func (mr *MerchantReads) GetMerchantPaymentMethodByIdMerchantRepo(id int) ([]entity.PaymentMethods, []entity.PaymentMethods, error) {
	var listMerchantPaymentMethods []entity.PaymentMethods
	var listPaymentMethods []entity.PaymentMethods

	queryPaymentMethod := `
	SELECT *
	FROM payment_methods pm
	ORDER BY pm.created_at DESC;
	`

	err := mr.db.Select(&listPaymentMethods, queryPaymentMethod)
	if err != nil {
		return listMerchantPaymentMethods, listPaymentMethods, err
	}

	queryListMerchantPaymentMethod := `
	SELECT
		pm.ID,
		pm.name,
		pm.pay_type,
		pm.created_at,
		pm.updated_at
	FROM
		merchant_payment_methods mpm
		JOIN payment_methods pm ON pm.ID = mpm.payment_method_id
		JOIN merchants m ON m.ID = mpm.merchant_id
	WHERE m.ID = $1
	ORDER BY mpm.created_at DESC;
	`

	err = mr.db.Select(&listMerchantPaymentMethods, queryListMerchantPaymentMethod, id)
	if err != nil {
		return listMerchantPaymentMethods, listPaymentMethods, err
	}

	return listMerchantPaymentMethods, listPaymentMethods, err
}

func (mr *MerchantReads) GetActiveAndAvailableChannelRepo(merchantPaychannelId int, paymentMethodName string) ([]entity.ProviderPaychannelEntity, []entity.ProviderPaychannelEntity, error) {
	var activePaychannel []entity.ProviderPaychannelEntity
	var availablePaychannel []entity.ProviderPaychannelEntity

	queryListActivePaychannel := `
	SELECT
		pp.ID,
		pp.paychannel_name,
		p.provider_name,
		p.currency,
		pp.min_transaction,
		pp.max_transaction
	FROM
		paychannel_routings pr
		JOIN provider_paychannels pp ON pp.ID = pr.provider_paychannel_id
		JOIN merchant_paychannels mp ON mp.ID = pr.merchant_paychannel_id
		JOIN merchant_payment_methods mpm ON mp.merchant_payment_method_id = mpm.ID
		JOIN provider_payment_methods ppm ON pp.provider_payment_method_id = ppm.ID
		JOIN providers p ON ppm.provider_id = p.ID
		JOIN payment_methods pm ON pm.ID = mpm.payment_method_id
	WHERE
		mp.ID = $1
	ORDER BY pr.created_at DESC;
	`

	err := mr.db.Select(&activePaychannel, queryListActivePaychannel, merchantPaychannelId)
	if err != nil {
		return activePaychannel, availablePaychannel, err
	}

	queryListAvailabelChannel := `
	SELECT
		pp.ID,
		pp.paychannel_name,
		p.provider_name,
		p.currency,
		pp.min_transaction,
		pp.max_transaction
	FROM
		provider_paychannels pp
		JOIN provider_payment_methods ppm ON ppm.ID = pp.provider_payment_method_id
		JOIN providers p ON ppm.provider_id = p.ID
		JOIN payment_methods pm ON pm.ID = ppm.payment_method_id
	WHERE
		pm.name = $1
	ORDER BY pp.created_at DESC;
	`

	err = mr.db.Select(&availablePaychannel, queryListAvailabelChannel, paymentMethodName)
	if err != nil {
		return activePaychannel, availablePaychannel, err
	}

	return activePaychannel, availablePaychannel, nil
}

func buildCountQueryManualPayment(params dto.QueryParamsManualPayment) string {

	query := `
	SELECT
		COUNT(*)
	FROM
		merchant_capital_flows mcf
		JOIN merchant_accounts ma ON mcf.merchant_account_id = ma.ID
		JOIN reason_lists r ON mcf.reason_id = r.ID
		JOIN merchants m ON ma.merchant_id = m.merchant_id
	WHERE
		mcf.reason_id NOT IN (5,6,7)
	`

	var conditions []string

	if params.MinDate != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.created_at >= '%v'", params.MinDate))
	}

	if params.MaxDate != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.created_at <= '%v'", params.MaxDate))
	}

	if params.AmountMax != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.amount <= '%v'", params.AmountMax))
	}

	if params.AmountMin != "" {
		conditions = append(conditions, fmt.Sprintf("mcf.amount >= '%v'", params.AmountMin))
	}

	if params.Status != "" {
		sliceStatus := helper.SplitString(params.Status)
		var statusConditions []string
		for _, status := range sliceStatus {
			statusConditions = append(statusConditions, fmt.Sprintf("mcf.status = '%v'", status))
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if params.ReasonName != "" {
		sliceReasonName := helper.SplitString(params.ReasonName)
		var reasonNameConditions []string
		for _, reasonName := range sliceReasonName {
			reasonNameConditions = append(reasonNameConditions, fmt.Sprintf("r.reason_name = '%v'", reasonName))
		}
		conditions = append(conditions, "("+strings.Join(reasonNameConditions, " OR ")+")")
	}

	if params.MerchantName != "" {
		sliceMerchantName := helper.SplitString(params.MerchantName)
		var merchantNameConditions []string
		for _, merchantName := range sliceMerchantName {
			merchantNameConditions = append(merchantNameConditions, fmt.Sprintf("m.merchant_name = '%v'", merchantName))
		}
		conditions = append(conditions, "("+strings.Join(merchantNameConditions, " OR ")+")")
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(mcf.payment_id LIKE '%v')", searchStr))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	return query
}
