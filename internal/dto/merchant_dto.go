package dto

import (
	"time"

	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
)

type CreateMerchantCapitalFlowPayload struct {
	PaymentId         string
	MerchantAccountId int
	TempBalance       float64
	ReasonId          int
	Status            string
	CreateBy          string
	Amount            float64
	Notes             string
	CapitalType       string
	ReverseFrom       string
}

type AdjustBalanceReqPayload struct {
	Amount     int    `json:"amount"`
	Notes      string `json:"notes"`
	MerchantId string `json:"merchantId"`
	Pin        string `json:"pin"`
	Username   string
}

type AdjustLimitOrFeePayload struct {
	MerchantPaychannelId int      `json:"merchantPaychannelId"`
	Username             string   `json:"username"`
	MinAmount            *float64 `json:"minAmount,omitempty"`
	MaxAmount            *float64 `json:"maxAmount,omitempty"`
	MaxDailyLimit        *float64 `json:"maxDailyAmount,omitempty"`
	Fee                  *float64 `json:"fee,omitempty"`
	FeeType              *string  `json:"feeType,omitempty"`
}

type SendCallbackReqPayload struct {
	PaymentId string `json:"transactionId"`
}

type BalanceTrfReqPayload struct {
	AccountFrom AccountData `json:"accountFrom"`
	AccountTo   AccountData `json:"accountTo"`
	Pin         string      `json:"pin"`
	Amount      int         `json:"amount"`
	Notes       string      `json:"notes"`
	Username    string
}

type AccountData struct {
	MerchantId           string `json:"merchantId"`
	MerchantPaychannelId int    `json:"merchantPaychannelId"`
}

type MerchantCallbackDto struct {
	TransactionId         string    `json:"transactionId"`
	MerchantTransactionId string    `json:"merchantTransactionId"`
	Status                string    `json:"status"`
	FailedReason          string    `json:"failedReason,omitempty"`
	Amount                float64   `json:"amount"`
	TransactionType       string    `json:"transactionType"`
	TransactionCreatedAt  time.Time `json:"transactionCreatedAt"`
	TransactionUpdatedAt  time.Time `json:"transactionUpdatedAt"`
}

type ManualPaymentDetailDto struct {
	DebitedAccount  string `json:"debitedAccount"`
	CreditedAccount string `json:"creditedAccount"`
	Reason          string `json:"reason"`
	Notes           string `json:"notes"`
}

type CreateMerchantDtoReq struct {
	MerchantName  string              `json:"merchantName"`
	PaymentMethod []PaymentMethodData `json:"paymentMethod"`
}

type GetMerchantAnalyticsDtoReq struct {
	MinDate              string `json:"minDate"`
	MaxDate              string `json:"maxDate"`
	MerchantId           string `json:"merchantId"`
	MerchantPaychannelId int    `json:"merchantPaychannelId"`
}

type PaymentMethodData struct {
	PaymentMethodId         int     `json:"paymentMethodId"`
	PaymentMethodName       string  `json:"paymentMethodName"`
	MinAmountPerTransaction float64 `json:"minAmount"`
	MaxAmountPerTransaction float64 `json:"maxAmount"`
	DailyLimit              float64 `json:"dailyLimit"`
	Fee                     float64 `json:"fee"`
	FeeType                 string  `json:"feeType"`
}

type AnalyticsMerchantRespDto struct {
	TransactionIn  AnalyticsDataRespDto `json:"transactionIn"`
	TransactionOut AnalyticsDataRespDto `json:"transactionOut"`
}

type HomeAnalyticsRespDto struct {
	TransactionIn       HomeAnalyticsDataRespDto `json:"transactionIn"`
	TransactionOut      HomeAnalyticsDataRespDto `json:"transactionOut"`
	TotalSuccessPayment SuccessPayment           `json:"totalSuccessPayment"`
}

type SuccessPayment struct {
	Qris           HomeAnalyticsDataRespDto `json:"qris"`
	Ewallet        HomeAnalyticsDataRespDto `json:"eWallet"`
	VirtualAccount HomeAnalyticsDataRespDto `json:"virtualAccount"`
	Disbursement   HomeAnalyticsDataRespDto `json:"disbursement"`
}

type HomeAnalyticsDataRespDto struct {
	TotalNumber int     `json:"totalNumber"`
	TotalAmount float64 `json:"totalAmount"`
}

type AnalyticsDataRespDto struct {
	TotalVolume        float64 `json:"totalVolume,omitempty"`
	SuccessRate        string  `json:"successRate,omitempty"`
	CompletionRate     string  `json:"completionRate,omitempty"`
	TransactionTotal   int     `json:"transactionTotal,omitempty"`
	SuccessTransaction int     `json:"successTransaction,omitempty"`
	FailedTransaction  int     `json:"failedTransaction,omitempty"`
}

type MerchantDataDtoRes struct {
	Id           int    `json:"id"`
	MerchantId   string `json:"merchantId"`
	MerchantName string `json:"merchantName"`
}

type MerchantPaychannelAnalyticsRspDto struct {
	AnalyticsData            AnalyticsDataRespDto      `json:"analyticsData"`
	MerchantPaychannelDetail entity.MerchantPaychannel `json:"merchantPaychannelDetailData"`
}

type AddSegmentDtoReq struct {
	MerchantPaychannelId int    `json:"merchantPaychannelId"`
	TierName             string `json:"tierName"`
}

type AddPaychannelRouting struct {
	ProviderPaychannelId []int `json:"providerPaychannelRouting"`
}

type ActiveAvailableChannelRespDto struct {
	ActivePaychannel    []entity.ProviderPaychannelEntity `json:"activePaychannel"`
	AvailablePaychannel []entity.ProviderPaychannelEntity `json:"availablePaychannel"`
}

type HomeAnalyticsDto struct {
	MinDate  string `json:"minDate"`
	MaxDate  string `json:"maxDate"`
	Username string `json:"username"`
}

type CallbackMerchantResp struct {
	StatusCallback              string `json:"statusCallback"`
	TransactionStatusInCallback string `json:"transactionStatusInCallback"`
	BeginsAt                    string `json:"beginAt"`
	LatestAt                    string `json:"latestAt"`
	MerchantResponse            string `json:"merchantResponse"`
}

type MerchantDisbursement struct {
	Amount            int    `json:"amount"`
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
	Note              string `json:"note"`
	Pin               string `json:"pin"`
	Username          string
}
