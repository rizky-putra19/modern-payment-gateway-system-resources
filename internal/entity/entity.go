package entity

import "time"

type Permission struct {
	PermissionID   int    `db:"permission_id" json:"permissionId"`
	PermissionDesc string `db:"permission_desc" json:"permissionDesc"`
}

type User struct {
	UserID         int          `db:"user_id" json:"userId"`
	Email          string       `db:"email" json:"email,omitempty"`
	Username       string       `db:"username" json:"username,omitempty"`
	UserType       string       `db:"user_type" json:"userType"`
	Password       string       `db:"password" json:"-"`
	Pin            string       `db:"pin" json:"-"`
	MerchantID     *string      `db:"merchant_id" json:"merchantId,omitempty"`
	MerchantName   *string      `db:"merchant_name" json:"merchantName,omitempty"`
	MerchantSecret *string      `db:"merchant_secret" json:"-"`
	Currency       *string      `db:"currency" json:"currency,omitempty"`
	MerchantStatus *string      `db:"merchant_status" json:"merchantStatus,omitempty"`
	RoleId         int          `db:"role_id" json:"-"`
	RoleName       string       `db:"role_name" json:"roleName"`
	Permissions    []Permission `json:"permissions"`
	UserStatus     string       `db:"user_status" json:"userStatus"`
	UserCreatedAt  time.Time    `db:"user_created_at" json:"createdAt"`
	UserUpdatedAt  time.Time    `db:"user_updated_at" json:"updatedAt"`
}

// Transaction represents a transaction record.
type Transaction struct {
	Id                      int     `db:"id" json:"id"`
	PaymentID               string  `db:"payment_id" json:"transactionId"`
	MerchantReferenceNumber string  `db:"merchant_reference_number" json:"merchantReferenceNumber"`
	TransactionAmount       float64 `db:"transaction_amount" json:"transactionAmount"`
	TransactionStatus       string  `db:"transaction_status" json:"transactionStatus"`
	MerchantId              string  `db:"merchant_id" json:"merchantId"`
	MerchantName            string  `db:"merchant_name" json:"merchantName"`
	PaymentMethodName       string  `db:"payment_method_name" json:"paymentMethodName"`
	PaymentMethodType       string  `db:"payment_method_type" json:"paymentMethodType"`
	ProviderPaychannelName  string  `db:"provider_paychannel_name" json:"providerPaychannelName"`
	BankName                *string `db:"bank_name" json:"bankName"`
	BankCode                *string `db:"bank_code" json:"bankCode"`
	TransactionCreatedAt    string  `db:"transaction_created_at" json:"createdAt"`
	TransactionUpdatedAt    string  `db:"transaction_updated_at" json:"updatedAt"`
}

type AccountData struct {
	Id              int     `db:"id" json:"accountInformationId"`
	PaymentId       string  `db:"payment_id" json:"transactionId"`
	AccountName     *string `db:"account_name" json:"accountName"`
	AccountNumber   *string `db:"account_number" json:"accountNumber"`
	ReferenceNumber *string `db:"reference_number" json:"referenceNumber"`
	BankName        *string `db:"bank_name" json:"bankName"`
	BankCode        *string `db:"bank_code" json:"bankCode"`
	Remark          *string `db:"remark" json:"remark"`
	PhoneNumber     *string `db:"phone_number" json:"phoneNumber"`
	Email           *string `db:"email" json:"email"`
	AccountType     *string `db:"account_type" json:"accountType"`
	CreatedAt       string  `db:"created_at" json:"createdAt"`
}

type PaymentDetailMerchantProvider struct {
	TransactionID          int       `db:"transaction_id"`
	PaymentID              string    `db:"payment_id"`
	MerchantRefNumber      string    `db:"merchant_reference_number"`
	ProviderRefNumber      string    `db:"provider_reference_number"`
	TransactionAmount      float64   `db:"transaction_amount"`
	BankCode               string    `db:"bank_code"`
	Status                 string    `db:"status"`
	ClientIPAddress        string    `db:"client_ip_address"`
	MerchantCallbackURL    string    `db:"merchant_callback_url"`
	RequestMethod          string    `db:"request_method"`
	PaymentMethodName      string    `db:"payment_method_name"`
	PayType                string    `db:"pay_type"`
	TransactionCreatedAt   time.Time `db:"transaction_created_at"`
	TransactionUpdatedAt   time.Time `db:"transaction_updated_at"`
	MerchantPaychannelID   int       `db:"merchant_payment_method_id"`
	Segment                *string   `db:"segment"`
	MerchantId             string    `db:"merchant_id"`
	MerchantName           string    `db:"merchant_name"`
	MerchantFee            float64   `db:"merchant_fee"`
	MerchantFeeType        string    `db:"merchant_fee_type"`
	MerchantStatus         string    `db:"merchant_status"`
	MerchantMinTrans       float64   `db:"merchant_min_transaction"`
	MerchantMaxTrans       float64   `db:"merchant_max_transaction"`
	MerchantMaxDailyTrans  float64   `db:"merchant_max_daily_transaction"`
	MerchantPaychannelCode string    `db:"merchant_paychannel_code"`
	MerchantCreatedAt      time.Time `db:"merchant_created_at"`
	MerchantUpdatedAt      time.Time `db:"merchant_updated_at"`
	ProviderName           string    `db:"provider_name"`
	ProviderPaychannelID   int       `db:"provider_payment_method_id"`
	ProviderBankCode       int       `db:"provider_bank_code"`
	PaychannelName         string    `db:"paychannel_name"`
	ProviderFee            float64   `db:"provider_fee"`
	ProviderFeeType        string    `db:"provider_fee_type"`
	ProviderStatus         string    `db:"provider_status"`
	ProviderMinTrans       float64   `db:"provider_min_transaction"`
	ProviderMaxTrans       float64   `db:"provider_max_transaction"`
	ProviderMaxDailyTrans  float64   `db:"provider_max_daily_transaction"`
	InterfaceSetting       *string   `db:"interface_setting"`
	ProviderCreatedAt      time.Time `db:"provider_created_at"`
	ProviderUpdatedAt      time.Time `db:"provider_updated_at"`
}

type ProviderConfirmDetail struct {
	Id                 int       `db:"id" json:"id"`
	PaymentId          string    `db:"payment_id" json:"paymentId"`
	Type               string    `db:"type" json:"type"`
	ConfirmationResult string    `db:"confirmation_result" json:"confirmationResult"`
	Notes              *string   `db:"notes" json:"notes"`
	ReceivedAt         time.Time `db:"created_at" json:"receivedAt"`
	ProcessedAt        time.Time `db:"updated_at" json:"processedAt"`
}

type TransactionStatusLogs struct {
	Id        int       `db:"id" json:"id"`
	PaymentId string    `db:"payment_id" json:"paymentId"`
	StatusLog string    `db:"status_log" json:"statusLog"`
	ChangeBy  string    `db:"change_by" json:"changeBy"`
	Notes     *string   `db:"notes" json:"notes"`
	RealNotes *string   `db:"real_notes" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"changeAt"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type CapitalFlows struct {
	Id                int       `db:"id" json:"id"`
	PaymentId         string    `db:"payment_id" json:"paymentId"`
	MerchantAccountId int       `db:"merchant_account_id" json:"merchantAccountId"`
	MerchantId        string    `db:"merchant_id" json:"merchantCode"`
	MerchantName      string    `db:"merchant_name" json:"merchantName"`
	TempBalance       float64   `db:"temp_balance" json:"tempBalance"`
	ReasonId          int       `db:"reason_id" json:"reasonId"`
	ReasonName        string    `db:"reason_name" json:"reasonName"`
	ReasonDescription string    `db:"reason_description" json:"reasonDescription"`
	Status            string    `db:"status" json:"status"`
	Notes             *string   `db:"notes" json:"notes"`
	Amount            float64   `db:"amount" json:"amount"`
	CapitalType       string    `db:"capital_type" json:"capitalType"`
	CreatedBy         string    `db:"created_by" json:"createdBy"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
}

type TransactionCapitalFlows struct {
	Id              int       `db:"id" json:"id"`
	PaymentId       string    `db:"payment_id" json:"transactionId"`
	TransactionType string    `db:"reason_name" json:"transactionType"`
	Amount          float64   `db:"amount" json:"amount"`
	CreatedAt       time.Time `db:"created_at" json:"processedAt"`
	MerchantId      string    `db:"merchant_id" json:"correspondingAccount"`
	Status          string    `db:"status" json:"status"`
	TempBalance     float64   `db:"temp_balance" json:"balance"`
	CapitalType     string    `db:"capital_type" json:"capitalType"`
}

type MerchantCallback struct {
	Id                      int       `db:"id" json:"id"`
	PaymentId               string    `db:"payment_id" json:"paymentId"`
	CallbackStatus          string    `db:"callback_status" json:"callbackStatus"`
	PaymentStatusInCallback string    `db:"payment_status_in_callback" json:"paymentStatusInCallback"`
	CallbackResult          string    `db:"callback_result" json:"callbackResult"`
	TriggeredBy             string    `db:"triggered_by" json:"triggeredBy"`
	CallbackRequest         string    `db:"callback_request" json:"-"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	StartedAt               string    `json:"startedAt,omitempty"`
	RetriedAt               string    `json:"retriedAt,omitempty"`
	CallbackAt              string    `json:"callbackAt,omitempty"`
	CallbackRequestResp     string    `json:"callbackRequest,omitempty"`
}

type ListMerchantCallback struct {
	ID                      int         `db:"id" json:"id"`
	PaymentID               string      `db:"payment_id" json:"paymentId"`
	CallbackStatus          string      `db:"callback_status" json:"callbackStatus"`
	PaymentStatusInCallback string      `db:"payment_status_in_callback" json:"paymentStatusInCallback"`
	CallbackResult          string      `db:"callback_result" json:"callbackResult"`
	LatestCreatedAt         time.Time   `db:"latest_created_at" json:"-"`
	FirstCreatedAt          time.Time   `db:"first_created_at" json:"startedAt"`
	LastCreatedAt           time.Time   `db:"last_created_at" json:"retriedAt"`
	MerchantReferenceNumber string      `db:"merchant_reference_number" json:"merchantReferenceNumber"`
	MerchantId              string      `db:"merchant_id" json:"-"`
	MerchantName            string      `db:"merchant_name" json:"merchantName"`
	MerchantDetailData      interface{} `json:"merchantDetailData"`
	PaymentType             string      `db:"pay_type" json:"payType"`
}

type MerchantAccount struct {
	Id                    int       `db:"id" json:"id"`
	MerchantId            string    `db:"merchant_id" json:"merchantId"`
	SettledBalance        float64   `db:"settle_balance" json:"settleBalance"`
	NotSettledBalance     float64   `db:"not_settle_balance" json:"notSettleBalance"`
	HoldBalance           float64   `db:"hold_balance" json:"holdBalance"`
	BalanceCapitalFlow    float64   `db:"balance_capital_flow" json:"balanceCapitalFlow"`
	PendingTransactionOut float64   `db:"pending_transaction_out" json:"pendingTransactionOut"`
	CreatedAt             time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time `db:"updated_at" json:"updatedAt"`
}

type ManualPayment struct {
	Id          int       `db:"id" json:"id"`
	PaymentId   string    `db:"payment_id" json:"paymentId"`
	MerchantId  string    `db:"merchant_id" json:"merchantAccount"`
	ReasonId    int       `db:"reason_id" json:"-"`
	ReasonName  string    `db:"reason_name" json:"reasonName"`
	Amount      float64   `db:"amount" json:"amount"`
	Status      string    `db:"status" json:"status"`
	Notes       *string   `db:"notes" json:"-"`
	CapitalType string    `db:"capital_type" json:"capitalType"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type PaymentMethods struct {
	Id        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	PayType   string    `db:"pay_type" json:"payType"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type Providers struct {
	Id           int       `db:"id" json:"id"`
	ProviderId   string    `db:"provider_id" json:"providerId"`
	ProviderName string    `db:"provider_name" json:"providerName"`
	Currency     string    `db:"currency" json:"currency"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

type Merchants struct {
	Id             int       `db:"id" json:"id"`
	MerchantId     string    `db:"merchant_id" json:"merchantId"`
	MerchantName   string    `db:"merchant_name" json:"merchantName"`
	MerchantSecret string    `db:"merchant_secret" json:"-"`
	Currency       string    `db:"currency" json:"currency"`
	Status         string    `db:"status" json:"status"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}

type MerchantSecret struct {
	SecretKey string `db:"merchant_secret" json:"secretKey"`
}

type PayChannels struct {
	Id             int       `db:"id" json:"id"`
	PayChannelName string    `db:"paychannel_name" json:"payChannelName"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}

type Reasons struct {
	Id         int       `db:"id" json:"id"`
	ReasonName string    `db:"reason_name" json:"reasonName"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
}

type MerchantPaychannel struct {
	Id                      int       `db:"id" json:"id"`
	MerchantPaymentMethodId int       `db:"merchant_payment_method_id" json:"-"`
	PayChannelCode          string    `db:"merchant_paychannel_code" json:"merchantPaychannelCode"`
	PaymentMethodChannel    string    `db:"name" json:"paymentMethodChannel"`
	PayTypeChannel          string    `db:"pay_type" json:"payTypeChannel"`
	Fee                     float64   `db:"fee" json:"channelFee"`
	FeeType                 string    `db:"fee_type" json:"feeType"`
	Status                  string    `db:"status" json:"status"`
	Segment                 string    `db:"segment" json:"segment"`
	MinTransaction          float64   `db:"min_transaction" json:"minTransaction"`
	MaxTransaction          float64   `db:"max_transaction" json:"maxTransaction"`
	MaxDailyTransaction     float64   `db:"max_daily_transaction" json:"maxDailyTransaction"`
	CreatedAt               time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt               time.Time `db:"updated_at" json:"updatedAt"`
	ActiveAvailableChannel  string    `json:"activeAvailableChannel"`
}

type ListMerchantAccountDto struct {
	Id                    int       `db:"id" json:"id"`
	MerchantId            string    `db:"merchant_id" json:"merchantId"`
	MerchantName          string    `db:"merchant_name" json:"merchantName"`
	SettleBalane          float64   `db:"settle_balance" json:"settleBalance"`
	NotSettleBalance      float64   `db:"not_settle_balance" json:"notSettleBalance"`
	HoldBalance           float64   `db:"hold_balance" json:"holdBalance"`
	PendingTransactionOut float64   `db:"pending_transaction_out" json:"pendingTransactionOut"`
	BalanceCapitalFlow    float64   `db:"balance_capital_flow" json:"balanceCapitalFlow"`
	CreatedAt             time.Time `db:"created_at" json:"createdAt"`
}

type RoutedPaychanneDto struct {
	Id                     int     `db:"id" json:"id"`
	ProviderPaychannelName string  `db:"paychannel_name" json:"providerPaychannelName"`
	ProviderName           string  `db:"provider_name" json:"providerName"`
	PaymentMethodType      string  `db:"pay_type" json:"paymentType"`
	PaymentMethodName      string  `db:"name" json:"paymentMethod"`
	FeeResp                string  `json:"fees"`
	FeesDb                 float64 `db:"fee" json:"-"`
	FeeType                string  `db:"fee_type" json:"-"`
	MinTransaction         float64 `db:"min_transaction" json:"minTransaction"`
	MaxTransaction         float64 `db:"max_transaction" json:"maxTransaction"`
	MaxDailyTransaction    float64 `db:"max_daily_transaction" json:"maxDailyTransaction"`
	Status                 string  `db:"status" json:"status"`
}

type BankListDto struct {
	Id                int       `db:"id" json:"id"`
	BankName          string    `db:"bank_name" json:"bankName"`
	BankCode          string    `db:"bank_code" json:"bankCode"`
	CheckListFlagging bool      `json:"checkListFlagging"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
}
