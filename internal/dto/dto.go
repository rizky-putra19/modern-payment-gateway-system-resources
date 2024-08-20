package dto

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	UserType string `json:"userType"`
	RoleName string `json:"roleName"`
	jwt.StandardClaims
}

type ResponseDto struct {
	ResponseCode    int         `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Data            interface{} `json:"data,omitempty"`
	Pagination      interface{} `json:"pagination,omitempty"`
}

type QueryParams struct {
	MinDate       string `json:"minDate"`
	MaxDate       string `json:"maxDate"`
	AmountMax     string `json:"amountMax"`
	AmountMin     string `json:"amountMin"`
	PaymentMethod string `json:"paymentMethod"`
	PayType       string `json:"payType"`
	Status        string `json:"status"`
	RequestMethod string `json:"requestMethod"`
	PayChannel    string `json:"payChannel"`
	MerchantName  string `json:"merchantName"`
	MerchantId    string `json:"merchantId"`
	ProviderName  string `json:"providerName"`
	Search        string `json:"search"`
	Page          string `json:"page"`
	PageSize      string `json:"pageSize"`
	Username      string `json:"username"`
	ReasonName    string `json:"reasonName"`
}

type QueryParamsMerchantCallback struct {
	PayType        string `json:"payType"`
	CallbackStatus string `json:"callbackStatus"`
	MerchantName   string `json:"merchantName"`
	Search         string `json:"search"`
	Page           string `json:"page"`
	PageSize       string `json:"pageSize"`
	MinDate        string `json:"minDate"`
	MaxDate        string `json:"maxDate"`
	Username       string `json:"username"`
}

type QueryParamsManualPayment struct {
	AmountMin    string `json:"amountMin"`
	AmountMax    string `json:"amountMax"`
	MinDate      string `json:"minDate"`
	MaxDate      string `json:"maxDate"`
	ReasonName   string `json:"reasonName"`
	Status       string `json:"status"`
	MerchantName string `json:"merchantName"`
	Search       string `json:"search"`
	PageSize     string `json:"pageSize"`
	Page         string `json:"page"`
	Username     string `json:"username"`
}

type PaginatedResponse struct {
	CurrentPage int  `json:"currentPage"`
	PageSize    int  `json:"pageSize"`
	TotalItems  int  `json:"totalItems"`
	TotalPages  int  `json:"totalPages"`
	HasNextPage bool `json:"hasNextPage"`
	HasPrevPage bool `json:"hasPrevPage"`
}

type FeeResponse struct {
	MerchantFee      FeeData `json:"merchantFee"`
	ProviderFee      FeeData `json:"providerFee"`
	CalculatedProfit float64 `json:"calculatedProfit"`
	FeeType          string  `json:"feeType"`
}

type FeeData struct {
	ConfiguredFee float64 `json:"configuredFee"`
	CalculatedFee float64 `json:"calculatedFee"`
	ChargedFee    float64 `json:"chargedFee"`
}

type FilterResponseDto struct {
	PaymentMethod interface{} `json:"paymentMethod"`
	ProviderName  interface{} `json:"providerName"`
	MerchantName  interface{} `json:"merchantName"`
	PayChannel    interface{} `json:"payChannel"`
	Reason        interface{} `json:"reason"`
	FeeType       interface{} `json:"feeType"`
	PayType       interface{} `json:"payType"`
}

type PayType struct {
	Id      int    `json:"id"`
	PayType string `json:"payType"`
}

type FeeType struct {
	Id      int    `json:"id"`
	FeeType string `json:"feeType"`
}

// TransactionStatusLog represents a record in the transaction_status_logs table
type TransactionStatusLog struct {
	PaymentID string
	Status    string
	CreatedAt time.Time
}

type TiersDto struct {
	Id    int    `json:"id"`
	Tiers string `json:"tiers"`
}

type GoogleCredentials struct {
	Type           string `json:"type"`
	ProjectId      string `json:"project_id"`
	PrivateKeyId   string `json:"private_key_id"`
	PrivateKey     string `json:"private_key"`
	ClientEmail    string `json:"client_email"`
	ClientId       string `json:"client_id"`
	AuthUri        string `json:"auth_uri"`
	TokenUri       string `json:"token_uri"`
	AuthProvider   string `json:"auth_provider_x509_cert_url"`
	ClientCertUrl  string `json:"client_x509_cert_url"`
	UniverseDomain string `json:"universe_domain"`
}
