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

type ProviderCredentialsEntity struct {
	Id               int       `db:"id"`
	ProviderId       string    `db:"provider_id"`
	Key              string    `db:"key"`
	Value            string    `db:"value"`
	InterfaceSetting string    `db:"interface_setting"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
