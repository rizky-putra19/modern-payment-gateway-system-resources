package psql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/jmoiron/sqlx"
)

type ProviderReads struct {
	db *sqlx.DB
}

func NewProviderReads(db *sqlx.DB) *ProviderReads {
	return &ProviderReads{
		db: db,
	}
}

func (pr *ProviderReads) CountProviderChannelByPaymentMethodRepo(paymentMethodName string) (int, error) {
	var totalItems int

	query := `
	SELECT
		COUNT(*)
	FROM
		provider_paychannels pp
		JOIN provider_payment_methods ppm ON ppm.ID = pp.provider_payment_method_id
		JOIN payment_methods pm ON pm.ID = ppm.payment_method_id
	WHERE
		pm.name = $1;
	`

	err := pr.db.Get(&totalItems, query, paymentMethodName)
	if err != nil {
		return 0, err
	}

	return totalItems, nil
}

func (pr *ProviderReads) CountActiveProviderChannelRepo(merchantPaychannelId int) (int, error) {
	var totalItems int

	query := `
	SELECT
		COUNT(*)
	FROM
		paychannel_routings pr
	WHERE
		pr.merchant_paychannel_id = $1;
	`

	err := pr.db.Get(&totalItems, query, merchantPaychannelId)
	if err != nil {
		return 0, err
	}

	return totalItems, nil
}

func (pr *ProviderReads) GetListProvidersWithFilterRepo(paymentMethod string, search string) ([]entity.ProviderListEntity, error) {
	var listProvider []entity.ProviderListEntity

	query := `
	SELECT
		p.ID,
		p.provider_id,
		p.provider_name,
		p.created_at,
		p.currency,
		STRING_AGG(pm.name, ', ') AS payment_methods
	FROM
		providers p
		JOIN provider_payment_methods ppm ON ppm.provider_id = p.ID
		JOIN payment_methods pm ON pm.ID = ppm.payment_method_id
	`

	var conditions []string

	if paymentMethod != "" {
		slicePaymentMethod := helper.SplitString(paymentMethod)
		var paymentMethodConditions []string
		for _, paymentMethod := range slicePaymentMethod {
			paymentMethodConditions = append(paymentMethodConditions, fmt.Sprintf("pm.name = '%v'", paymentMethod))
		}
		conditions = append(conditions, "("+strings.Join(paymentMethodConditions, " OR ")+")")
	}

	if search != "" {
		searchStr := fmt.Sprintf("%%%v%%", search)
		conditions = append(conditions, fmt.Sprintf("(p.provider_name LIKE '%v' OR p.provider_id LIKE '%v')", searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
	GROUP BY p.ID, p.provider_id, p.provider_name, p.created_at
	ORDER BY p.created_at DESC
	`

	err := pr.db.Select(&listProvider, query)
	if err != nil {
		return listProvider, err
	}

	return listProvider, nil
}

func (pr *ProviderReads) CountInterfacesProviderPaychannelByIdProvider(id int) (int, int, error) {
	var activePaychannel int
	var providerPaymentMethods int

	query := `
	SELECT 
		COUNT(*)
	FROM
		provider_paychannels pp
		JOIN provider_payment_methods ppm ON ppm.ID = pp.provider_payment_method_id
		JOIN providers p ON p.ID = ppm.provider_id
	WHERE
		p.ID = $1;
	`

	err := pr.db.Get(&activePaychannel, query, id)
	if err != nil {
		return 0, 0, err
	}

	queryProviderPaymentMethods := `
	SELECT
		COUNT(*)
	FROM
		provider_payment_methods ppm
		JOIN providers p ON ppm.provider_id = p.ID
	WHERE
		p.ID = $1;
	`
	err = pr.db.Get(&providerPaymentMethods, queryProviderPaymentMethods, id)
	if err != nil {
		return 0, 0, err
	}

	return activePaychannel, providerPaymentMethods, nil
}

func (pr *ProviderReads) GetProviderInterfacesRepoById(id int) ([]entity.ProviderInterfacesEntity, error) {
	var listData []entity.ProviderInterfacesEntity

	query := `
	SELECT
		ppm.ID,
		p.currency,
		p.provider_name,
		pm.pay_type,
		pm.name,
		COUNT(ppbl.id) AS payment_operators
	FROM
		providers p
		JOIN provider_payment_methods ppm ON p.id = ppm.provider_id
		JOIN payment_methods pm ON ppm.payment_method_id = pm.id
		LEFT JOIN provider_payment_method_bank_lists ppbl ON ppm.id = ppbl.provider_payment_method_id
	WHERE
		p.id = $1
	GROUP BY
		ppm.ID,
		p.currency,
		p.provider_name,
		pm.pay_type,
		pm.name
	`

	err := pr.db.Select(&listData, query, id)
	if err != nil {
		return listData, err
	}

	if len(listData) > 0 {
		for i := range listData {
			listData[i].InterfaceCode = strings.ToUpper(listData[i].Provider) + "-" + constant.TransformPaymentMethodNameIntoCode[listData[i].PaymentMethod]
		}
	}

	return listData, nil
}

func (pr *ProviderReads) GetListProviderPaychannelById(id int) ([]entity.InterfacePaychannelEntity, error) {
	var listPaychannel []entity.InterfacePaychannelEntity

	query := `
	SELECT
		pp.ID,
		pp.paychannel_name,
		p.provider_name,
		p.currency,
		pp.min_transaction,
		pp.max_transaction,
		pp.max_daily_transaction,
		pp.status
	FROM
		provider_paychannels pp
		JOIN provider_payment_methods ppm ON ppm.ID = pp.provider_payment_method_id
		JOIN providers p ON ppm.provider_id = p.ID
		JOIN payment_methods pm ON pm.ID = ppm.payment_method_id
	WHERE
		ppm.ID = $1
	ORDER BY pp.created_at DESC;
	`

	err := pr.db.Select(&listPaychannel, query, id)
	if err != nil {
		return listPaychannel, err
	}

	return listPaychannel, nil
}

func (pr *ProviderReads) GetListProviderChannelAllRepo(params dto.QueryParams) ([]entity.ProviderPaychannelAllEntity, error) {
	var listProviderChannel []entity.ProviderPaychannelAllEntity

	query := `
	SELECT
		pp.ID,
		pp.paychannel_name,
		p.provider_name,
		p.currency,
		pm.name,
		pm.pay_type,
		pp.min_transaction,
		pp.max_transaction,
		pp.max_daily_transaction,
		pp.status
	FROM
		provider_paychannels pp
		JOIN provider_payment_methods ppm ON ppm.ID = pp.provider_payment_method_id
		JOIN providers p ON ppm.provider_id = p.ID
		JOIN payment_methods pm ON pm.ID = ppm.payment_method_id
	`

	var conditions []string

	if params.ProviderName != "" {
		providerNames := helper.SplitString(params.ProviderName)
		quotedProviderNames := make([]string, len(providerNames))
		for i, name := range providerNames {
			quotedProviderNames[i] = fmt.Sprintf("'%v'", name)
		}
		conditions = append(conditions, fmt.Sprintf("p.provider_name IN (%v)", strings.Join(quotedProviderNames, ", ")))
	}

	if params.Status != "" {
		statuses := helper.SplitString(params.Status)
		quotedStatuses := make([]string, len(statuses))
		for i, status := range statuses {
			quotedStatuses[i] = fmt.Sprintf("'%v'", status)
		}
		conditions = append(conditions, fmt.Sprintf("pp.status IN (%v)", strings.Join(quotedStatuses, ", ")))
	}

	if params.PayType != "" {
		payTypes := helper.SplitString(params.PayType)
		quotedPayType := make([]string, len(payTypes))
		for i, payType := range payTypes {
			quotedPayType[i] = fmt.Sprintf("'%v'", payType)
		}
		conditions = append(conditions, fmt.Sprintf("pm.pay_type IN (%v)", strings.Join(quotedPayType, ", ")))
	}

	if params.PaymentMethod != "" {
		paymentMethods := helper.SplitString(params.PaymentMethod)
		quotedPaymentMethod := make([]string, len(paymentMethods))
		for i, paymentMethod := range paymentMethods {
			quotedPaymentMethod[i] = fmt.Sprintf("'%v'", paymentMethod)
		}
		conditions = append(conditions, fmt.Sprintf("pm.name IN (%v)", strings.Join(quotedPaymentMethod, ", ")))
	}

	if params.Search != "" {
		searchStr := fmt.Sprintf("%%%v%%", params.Search)
		conditions = append(conditions, fmt.Sprintf("(p.provider_name LIKE '%v' OR pp.paychannel_name LIKE '%v')", searchStr, searchStr))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY pp.created_at DESC"

	err := pr.db.Select(&listProviderChannel, query)
	if err != nil {
		return listProviderChannel, err
	}

	return listProviderChannel, nil
}

func (pr *ProviderReads) GetAllCredentialsRepo(providerId string, interfaceSetting string) ([]entity.ProviderCredentialsEntity, error) {
	var listCredentials []entity.ProviderCredentialsEntity

	query := `
		SELECT *
		FROM provider_credentials pc
		WHERE pc.provider_id = $1
		AND pc.interface_setting = $2
		ORDER BY pc.created_at DESC;
	`

	err := pr.db.Select(&listCredentials, query, providerId, interfaceSetting)
	if err != nil && err != sql.ErrNoRows {
		return listCredentials, err
	}

	return listCredentials, nil
}
