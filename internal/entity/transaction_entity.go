package entity

import (
	"time"
)

type MerchantExportCapitalFlowEntity struct {
	Id               int       `db:"id" json:"id"`
	PaymentId        string    `db:"payment_id" json:"paymentId"`
	MerchantId       *string   `db:"merchant_id" json:"merchantId"`
	MerchantName     *string   `db:"merchant_name" json:"merchantName"`
	Amount           *float64  `db:"amount" json:"amount"`
	ReasonName       *string   `db:"reason_name" json:"reasonName"`
	PaymentMethod    *string   `db:"payment_method" json:"paymentMethod"`
	Fee              *float64  `db:"fee" json:"fee"`
	FeeType          *string   `db:"fee_type" json:"feeType"`
	MerchantBalance  *float64  `db:"temp_balance" json:"merchantBalance"`
	Provider         *string   `db:"provider_name" json:"provider"`
	PaychannelRouted *string   `db:"paychannel_name" json:"paychannelName"`
	ProviderFee      *float64  `db:"provider_fee" json:"providerFee"`
	ProviderFeeType  *string   `db:"provider_fee_type" json:"providerFeeType"`
	Status           *string   `db:"status" json:"status"`
	Notes            *string   `db:"notes" json:"notes"`
	CapitalType      *string   `db:"capital_type" json:"capitalType"`
	ReverseFrom      *string   `db:"reverse_from" json:"reverseFrom"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}

type ReportStoragesEntity struct {
	Id            int       `db:"id" json:"id"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	MerchantId    string    `db:"merchant_id" json:"merchantId"`
	MerchantName  string    `db:"merchant_name" json:"merchantName"`
	Currency      string    `db:"currency" json:"currency"`
	ExportType    string    `db:"export_type" json:"exportType"`
	Period        string    `db:"period" json:"period"`
	Status        string    `db:"status" json:"status"`
	ReportUrl     string    `db:"report_url" json:"reportUrl"`
	CreatedByUser string    `db:"created_by_user" json:"-"`
}

type MerchantTransactionList struct {
	Id                   int     `db:"id" json:"id"`
	PaymentId            string  `db:"payment_id" json:"transactionId"`
	MerchanId            string  `db:"merchant_id" json:"merchantId"`
	TransactionAmount    float64 `db:"transaction_amount" json:"transactionAmount"`
	TransactionStatus    string  `db:"transaction_status" json:"transactionStatus"`
	PaymentMethodName    string  `db:"payment_method_name" json:"paymentMethodName"`
	TransactionCreatedAt string  `db:"transaction_created_at" json:"createdAt"`
	TransactionUpdatedAt string  `db:"transaction_updated_at" json:"updatedAt"`
}
