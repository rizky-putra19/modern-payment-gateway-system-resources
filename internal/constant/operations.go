package constant

var ExportStatus = []string{
	"NO DATA",
	"PENDING",
	"FINISHED",
	"ERROR",
}

var ExportType = []string{
	"Balance transaction flow",
	"Transaction in",
	"Transaction out",
}

const (
	JackInquiryKeyUrlCred    = "INQUIRY_ACCOUNT_URL"
	JackGetBalanceKeyUrlCred = "GET_BALANCE_URL"
	JackDisbursementUrlCred  = "DISBURSEMENT_TRANSACTIONS_URL"
	JackApiKeyCred           = "API_KEY"
)

const (
	JackStatusOk             = 200
	JackStatusInvalid        = 422
	JackStateStatusConfirm   = "confirmed"
	JackStateStatusCompleted = "completed"
	JackStateStatusDeclined  = "declined"
	JackStateStatusCanceled  = "canceled"
)

const (
	JackDisbursementFirstSenderName  = "Hypay"
	JackDisbursementSecondSenderName = "Indonesia"
	JackDisbursementCountryIsoName   = "IDN"
	JackDisbursementCurrency         = "IDR"
	JackDisbursementNotes            = "Send By Hypay"
	JackDisbursementMode             = "DESTINATION"
)

const (
	SourceCallback = "CALLBACK"
	SourceQuery    = "QUERY"
)
