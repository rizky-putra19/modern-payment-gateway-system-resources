package internal

import (
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
)

type MerchantCallbackItf interface {
	SendCallbackAdptr(url string, transactionEntity entity.PaymentDetailMerchantProvider, transactionStatusLogLatest entity.TransactionStatusLogs, merchantSecret string) (interface{}, error)
}

type JackProviderItf interface {
	InquiryAccount(payload dto.MerchantDisbursement, credentials []entity.ProviderCredentialsEntity, bankCode string) (dto.InquiryAccountResponse, error)
	GetBalance(username string, credentials []entity.ProviderCredentialsEntity) (int, error)
	ConfirmDisbursement(payload dto.ConfirmTransactionPayload, credentials []entity.ProviderCredentialsEntity) (dto.CreateDisbursementRequestResponse, error)
	CreateDisbursement(payload dto.MerchantDisbursement, credentials []entity.ProviderCredentialsEntity, bankCode string, paymentId string) (dto.CreateDisbursementRequestResponse, error)
}
