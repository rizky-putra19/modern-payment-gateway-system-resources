package dto

import (
	"time"
)

type GetPaymentDetailsRequest struct {
	PaymentId string
	Username  string
}

type PaymentDetailProviderMerchantResponse struct {
	MerchantData DetailData `json:"merchantData"`
	ProviderData DetailData `json:"providerData"`
}

type DetailData struct {
	Name             string  `json:"name"`
	PaymentId        string  `json:"orderId"`
	Channel          string  `json:"channel"`
	Segment          string  `json:"segment"`
	Routing          string  `json:"routing"`
	Fee              float64 `json:"fee"`
	FeeType          string  `json:"feeType"`
	ClientIpAddress  string  `json:"clientIpAddress"`
	InterfaceSetting string  `json:"interfaceSetting"`
	RequestMethod    string  `json:"requestMethod"`
	RequestedAmount  float64 `json:"requestedAmount"`
}

type TransactionMerchantData struct {
	Name                  string  `json:"name"`
	MerchantTransactionId string  `json:"merchantTransactionId"`
	TransactionFee        float64 `json:"transactionFee"`
	TransactionFeeType    string  `json:"transactionFeeType"`
	TransactionAmount     float64 `json:"transactionAmount"`
	TransactionNetAmount  float64 `json:"transactionNetAmount"`
	TransactionMethod     string  `json:"transactionMethod"`
	AccountName           *string `json:"accountName"`
	AccountNumber         *string `json:"accountNumber"`
	BankName              *string `json:"bankName"`
	IpAddress             string  `json:"ipAddress"`
}

type TransactionsCapitalFlow struct {
	TransactionType string    `json:"transactionType"`
	MerchantAccount string    `json:"merchantAccount"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	CapitalType     string    `json:"capitalType"`
	CreatedAt       time.Time `json:"createdAt"`
}

type UpdateStatusTransaction struct {
	PaymentId string `json:"paymentId"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
	Pin       string `json:"pin"`
}

type CreateMerchantExportReqDto struct {
	ExportType string `json:"exportType"`
	MinDate    string `json:"minDate"`
	MaxDate    string `json:"maxDate"`
	MerchantId string `json:"merchantId"`
	UserType   string
}

type GetListMerchantExportFilter struct {
	ExportType   string `json:"exportType"`
	MinDate      string `json:"minDate"`
	MaxDate      string `json:"maxDate"`
	ExportStatus string `json:"exportStatus"`
	Search       string `json:"search"`
	Merchants    string `json:"merchants"`
}

type CreateReportStorageDto struct {
	MerchantId    string
	Period        string
	ExportType    string
	Status        string
	ReportUrl     string
	CreatedByUser string
	FileName      string
}

type ListFilterExportDto struct {
	ExportType   []ExportData       `json:"exportType"`
	ExportStatus []ExportStatusData `json:"exportStatus"`
}

type ExportData struct {
	Id             int    `json:"id"`
	ExportTypeName string `json:"exportTypeName"`
}

type ExportStatusData struct {
	Id               int    `json:"id"`
	ExportStatusName string `json:"exportStatusName"`
}

type InternalExportDto struct {
	Id         int       `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	Currency   string    `json:"currency"`
	ExportType string    `json:"exportType"`
	Period     string    `json:"period"`
	Status     string    `json:"status"`
	ReportUrl  string    `json:"reportUrl"`
}

type FailureInformationDto struct {
	FailedAt string `json:"failedAt"`
	Message  string `json:"message"`
}

type TransactionMerchantDetailDto struct {
	TransactionData    TransactionMerchantData `json:"transactionData"`
	MerchantCallback   CallbackMerchantResp    `json:"merchantCallback"`
	FailureInformation FailureInformationDto   `json:"failureInformation"`
}

type CountDisbursementTotalAmountDto struct {
	Amount   int `json:"amount"`
	Username string
}

type CountDisbursementRespDto struct {
	FeeAmount   float64 `json:"feeAmount"`
	TotalAmount float64 `json:"totalAmount"`
}
