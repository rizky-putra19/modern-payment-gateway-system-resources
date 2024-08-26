package entity

import "time"

type ProviderPaychannelEntity struct {
	Id             int     `db:"id" json:"id"`
	PaychannelName string  `db:"paychannel_name" json:"code"`
	Provider       string  `db:"provider_name" json:"provider"`
	Currency       string  `db:"currency" json:"currency"`
	MinAmount      float64 `db:"min_transaction" json:"minAmount"`
	MaxAmount      float64 `db:"max_transaction" json:"maxAmount"`
}

type InterfacePaychannelEntity struct {
	Id             int     `db:"id" json:"id"`
	PaychannelName string  `db:"paychannel_name" json:"code"`
	Provider       string  `db:"provider_name" json:"provider"`
	Currency       string  `db:"currency" json:"currency"`
	MinAmount      float64 `db:"min_transaction" json:"minAmount"`
	MaxAmount      float64 `db:"max_transaction" json:"maxAmount"`
	MaxDailyLimit  float64 `db:"max_daily_transaction" json:"maxDailyLimit"`
	Status         string  `db:"status" json:"status"`
}

type ProviderListEntity struct {
	Id            int       `db:"id" json:"id"`
	ProviderId    string    `db:"provider_id" json:"providerCode"`
	ProviderName  string    `db:"provider_name" json:"providerName"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	Currency      string    `db:"currency" json:"currencies"`
	PaymentMethod string    `db:"payment_methods" json:"paymentMethods"`
	Interfaces    string    `json:"interface"`
}

type ProviderInterfacesEntity struct {
	Id              int    `db:"id" json:"id"`
	InterfaceCode   string `json:"interfaceCode"`
	Currency        string `db:"currency" json:"currency"`
	Provider        string `db:"provider_name" json:"provider"`
	PaymentType     string `db:"pay_type" json:"paymentType"`
	PaymentMethod   string `db:"name" json:"paymentMethod"`
	PaymentOperator string `db:"payment_operators" json:"paymentOperators"`
}

type ProviderPaychannelAllEntity struct {
	Id            int     `db:"id" json:"id"`
	Paychannel    string  `db:"paychannel_name" json:"paychannel_name"`
	Currency      string  `db:"currency" json:"currency"`
	Provider      string  `db:"provider_name" json:"provider"`
	PaymentType   string  `db:"pay_type" json:"paymentType"`
	PaymentMethod string  `db:"name" json:"paymentMethod"`
	MinAmount     float64 `db:"min_transaction" json:"minAmount"`
	MaxAmount     float64 `db:"max_transaction" json:"maxAmount"`
	MaxDailyLimit float64 `db:"max_daily_transaction" json:"maxDailyLimit"`
	Status        string  `db:"status" json:"status"`
}

type ProviderChannelDetailEntity struct {
	Id               int     `db:"id" json:"id"`
	Paychannel       string  `db:"paychannel_name" json:"paychannel_name"`
	Currency         string  `db:"currency" json:"currency"`
	Provider         string  `db:"provider_name" json:"provider"`
	PaymentType      string  `db:"pay_type" json:"paymentType"`
	PaymentMethod    string  `db:"name" json:"paymentMethod"`
	Fee              float64 `db:"fee" json:"fee"`
	FeeType          string  `db:"fee_type" json:"feeType"`
	MinAmount        float64 `db:"min_transaction" json:"minAmount"`
	MaxAmount        float64 `db:"max_transaction" json:"maxAmount"`
	MaxDailyLimit    float64 `db:"max_daily_transaction" json:"maxDailyLimit"`
	InterfaceSetting string  `db:"interface_setting" json:"interfaceSetting"`
	Status           string  `db:"status" json:"status"`
}

type ProviderCredentialsEntity struct {
	Id               int       `db:"id"`
	ProviderId       string    `db:"provider_id"`
	Key              string    `db:"key"`
	Value            string    `db:"value"`
	InterfaceSetting string    `db:"interface_setting"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type ProviderRoutedChannelEntity struct {
	Id                      int       `db:"id" json:"id"`
	PaychannelName          string    `db:"merchant_paychannel_code" json:"paychannelName"`
	MerchantName            string    `db:"merchant_name" json:"merchantName"`
	Fee                     float64   `db:"fee" json:"fee"`
	FeeType                 string    `db:"fee_type" json:"feeType"`
	Status                  string    `db:"status" json:"status"`
	ActiveAvailableChannels string    `json:"activeAvailableChannels"`
	MinTransaction          float64   `db:"min_transaction" json:"minTransaction"`
	MaxTransaction          float64   `db:"max_transaction" json:"maxTransaction"`
	MaxDailyTransaction     string    `db:"max_daily_transaction" json:"maxDailyTransaction"`
	CreatedAt               time.Time `db:"created_at" json:"createdAt"`
	PaymentMethodName       string    `db:"name" json:"-"`
}

type ProviderPaychannelBankListEntity struct {
	Id                   int       `db:"id" json:"id"`
	ProviderPaychannelId int       `db:"provider_paychannel_id" json:"providerPaychannelId"`
	BankListId           int       `db:"bank_list_id" json:"bankListId"`
	CreatedAt            time.Time `db:"created_at" json:"createdAt"`
}
