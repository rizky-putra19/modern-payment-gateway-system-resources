package dto

import "github.com/hypay-id/backend-dashboard-hypay/internal/entity"

type GetProviderAnalyticsDtoReq struct {
	MinDate    string `json:"minDate"`
	MaxDate    string `json:"maxDate"`
	ProviderId int    `json:"providerId"`
}

type AnalyticsProviderRespDto struct {
	TransactionIn      AnalyticsDataRespDto              `json:"transactionIn"`
	TransactionOut     AnalyticsDataRespDto              `json:"transactionOut"`
	ProviderInterfaces []entity.ProviderInterfacesEntity `json:"providerInterfaces"`
}
