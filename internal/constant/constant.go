package constant

const CreateBySystem = "SYSTEM"
const BucketName = "hypaystagingreportstorage"
const BusinessHypayEmail = "business@hypay.id"
const MainType = "MAIN"
const GeneralErrMsg = "failed to do disbursement, please ask admin for more info"
const BankAccountNameSimilarityMatchInPercent = 0.3

const (
	IpAddressHypay   = "0.0.0.0"
	CallbackUrlHypay = "https://www.hypay.id"
)

const (
	ReasonNameSettlement     = "Settlement"
	ReasonNameInTransaction  = "In"
	ReasonNameOutTransaction = "Out"
	ReasonNameFeeTransaction = "Fee"
)

const (
	UserOperation = "USER_OPERATION"
	UserMerchant  = "USER_MERCHANT"
)

const (
	FlaggingTransactionList   = "Transaction List"
	FlaggingMerchantAnalytics = "Merchant Analytics"
)

const (
	AccountTypeCreditor = "CREDITOR"
	AccountTypeDebitor  = "DEBITOR"
)

const (
	FeeTypeFixedFee   = "FIXED_FEE"
	FeeTypePercentage = "PERCENTAGE"
)

const (
	RoleNameAdmin           = "admin"
	RoleNameCustomerSupport = "customer support"
	RoleNameFinance         = "finance"
)

const (
	StatusSuccess    = "SUCCESS"
	StatusFailed     = "FAILED"
	StatusProcessing = "PROCESSING"
	StatusReversed   = "REVERSED"
)

const (
	ReportStatusFinished = "FINISHED"
	ReportStatusError    = "ERROR"
	ReportStatusPending  = "PENDING"
	ReportStatusNoData   = "NO DATA"
)

const (
	StatusLogSuccess            = "CONFIRMED BY PROVIDER"
	StatusLogFailed             = "FAILED / REFUSED BY PROVIDER"
	StatusLogAcceptedByPlatform = "ACCEPTED BY PLATFORM"
	StatusLogAcceptedByProvider = "ACCEPTED BY PROVIDER"
)

const (
	CapitalTypeCredit            = "CREDIT"
	CapitalTypeDebit             = "DEBIT"
	CapitalTypeNotDebitNotCredit = "UNCHANGED"
)

const (
	PayTypePayin  = "in"
	PayTypePayout = "out"
)

const (
	ReasonIdPayin           = 5
	ReasonIdPayout          = 6
	ReasonIdFee             = 7
	ReasonIdTopUp           = 1
	ReasonIdHoldBalance     = 2
	ReasonIdSettlement      = 3
	ReasonIdBalanceTransfer = 8
	ReasonIdOutSettlement   = 4
)

const (
	SettleBalance     = "settledBalance"
	NotSettledBalance = "notSettledBalance"
)

const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
)

var PayType = []string{
	"in",
	"out",
}

var FeeType = []string{
	"FIXED_FEE",
	"PERCENTAGE",
}

var TransformPaymentMethodNameIntoCode = map[string]string{
	"Virtual account": "VA",
	"Qris":            "QRIS",
	"Ewallet":         "EWALLET",
	"Disbursement":    "DISBURSEMENT",
}

const (
	QrisPaymentMethod           = "Qris"
	VirtualAccountPaymentMethod = "Virtual account"
	EwalletPaymentMethod        = "Ewallet"
	DisbursementPaymentMethod   = "Disbursement"
)

var Alphabet = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

var TransformExportType = map[string]string{
	"Balance transaction flow": "BTF",
	"Transaction in":           "TI",
	"Transaction out":          "TO",
}

const (
	ExportTypeCapitalFlow = "Balance transaction flow"
	ExportTypeIn          = "Transaction in"
	ExportTypeOut         = "Transaction out"
)

const GetSecret = "service-account-credentials"
const InternalExport = "[INTERNAL-EXPORT]"

const (
	ProviderJack = "ID-JACK"
)

var ProviderListName = []string{
	"ID-JACK",
}
