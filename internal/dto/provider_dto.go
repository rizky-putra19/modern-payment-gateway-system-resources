package dto

import "github.com/hypay-id/backend-dashboard-hypay/internal/entity"

type GetProviderAnalyticsDtoReq struct {
	MinDate           string `json:"minDate"`
	MaxDate           string `json:"maxDate"`
	ProviderId        int    `json:"providerId"`
	ProviderChannelId int
}

type AnalyticsProviderRespDto struct {
	TransactionIn      AnalyticsDataRespDto              `json:"transactionIn"`
	TransactionOut     AnalyticsDataRespDto              `json:"transactionOut"`
	ProviderInterfaces []entity.ProviderInterfacesEntity `json:"providerInterfaces"`
}

type InquiryAccountResponse struct {
	Status int         `json:"status"`
	Data   InquiryData `json:"data"`
}

type InquiryData struct {
	ID           interface{} `json:"id"`
	AccountNo    string      `json:"account_no"`
	BankName     string      `json:"bank_name"`
	AccountName  string      `json:"account_name"`
	ErrorMessage string      `json:"error_message"`
	Errors       string      `json:"errors"`
}

type ConfirmTransactionPayload struct {
	Username   string
	ProviderID string
}

type CreateDisbursementRequest struct {
	ReferenceID string          `json:"reference_id"`
	CallbackURL string          `json:"callback_url"`
	PayerID     string          `json:"payer_id"`
	Mode        string          `json:"mode"`
	Sender      SenderData      `json:"sender"`
	Source      SourceData      `json:"source"`
	Destination DestinationData `json:"destination"`
	Beneficiary BeneficiaryData `json:"beneficiary"`
	Notes       string          `json:"notes"`
}

type CreateDisbursementRequestResponse struct {
	Status int                                   `json:"status"`
	Data   CreateDisbursementRequestResponseData `json:"data"`
}

type CreateDisbursementRequestResponseData struct {
	ID              int             `json:"id"`
	ReferenceID     string          `json:"reference_id"`
	CallbackURL     string          `json:"callback_url"`
	PayerID         int             `json:"payer_id"`
	Mode            string          `json:"mode"`
	Sender          SenderData      `json:"sender"`
	Source          SourceData      `json:"source"`
	Destination     DestinationData `json:"destination"`
	Beneficiary     BeneficiaryData `json:"beneficiary"`
	Notes           string          `json:"notes"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	UserID          int             `json:"user_id"`
	State           string          `json:"state"`
	Amount          int             `json:"amount"`
	PaidAt          string          `json:"paid_at"`
	Rate            string          `json:"rate"`
	Fee             string          `json:"fee"`
	PartnerID       int             `json:"partner_id"`
	SentAmount      string          `json:"sent_amount"`
	ErrorCode       interface{}     `json:"error_code"`
	ErrorMessage    string          `json:"error_message"`
	TransactionType string          `json:"transaction_type"`
}

type SenderData struct {
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	CountryIsoCode string `json:"country_iso_code"`
}

type SourceData struct {
	Currency       string `json:"currency"`
	CountryIsoCode string `json:"country_iso_code"`
}

type DestinationData struct {
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	CountryIsoCode string `json:"country_iso_code"`
}

type BeneficiaryData struct {
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	CountryIsoCode string `json:"country_iso_code"`
	Account        string `json:"account"`
	Bank           string `json:"bank"`
}

type JackBankCode struct {
	Id       string
	BankName string
}

type JackCredentialsDto struct {
	InquiryUrl      string
	GetBalanceUrl   string
	ApiKey          string
	DisbursementUrl string
}

type JackGetBalanceResponse struct {
	Status  int              `json:"status"`
	Data    JackBalancesData `json:"data"`
	Message string           `json:"message"`
}

type JackBalancesData struct {
	Total    int                     `json:"total"`
	Balances []JackBalanceDetailData `json:"balances"`
}

type JackBalanceDetailData struct {
	ID        int    `json:"id"`
	Currency  string `json:"currency"`
	Balance   int    `json:"balance"`
	IsActive  bool   `json:"is_active"`
	UserId    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	PartnerId int    `json:"partner_id"`
}

type ProviderChannelAnalyticsResDto struct {
	AnalyticsData         AnalyticsDataRespDto               `json:"analyticsData"`
	ProviderChannelDetail entity.ProviderChannelDetailEntity `json:"providerChannelDetailData"`
}

type AdjustLimitOrFeeProviderPayload struct {
	ProviderChannelId int      `json:"providerChannelId"`
	MinAmount         *float64 `json:"minAmount,omitempty"`
	MaxAmount         *float64 `json:"maxAmount,omitempty"`
	MaxDailyLimit     *float64 `json:"maxDailyAmount,omitempty"`
	Fee               *float64 `json:"fee,omitempty"`
	FeeType           *string  `json:"feeType,omitempty"`
	InterfaceSetting  *string  `json:"interfaceSetting,omitempty"`
}

type AddOperatorProviderChannelPayload struct {
	ProviderChannelId int    `json:"providerChannelId"`
	BankCode          string `json:"bankCode"`
	CheckListFlagging bool   `json:"checkListFlagging"`
}

type UpdateStatusProviderPaychannelDto struct {
	ProviderChannelId int `json:"providerChannelId"`
}

type CreateProviderChannelDto struct {
	ProviderInterfaceId int                                 `json:"providerInterfaceId"`
	PaychannelName      string                              `json:"paychannelName"`
	InterfaceSetting    *string                             `json:"interfaceSetting"`
	MinAmount           *float64                            `json:"minAmount"`
	MaxAmount           *float64                            `json:"maxAmount"`
	DailyLimit          *float64                            `json:"dailyLimit"`
	Fee                 *float64                            `json:"fee"`
	FeeType             *string                             `json:"feeType"`
	BankOperator        []AddOperatorProviderChannelPayload `json:"paymentOperator"`
}
