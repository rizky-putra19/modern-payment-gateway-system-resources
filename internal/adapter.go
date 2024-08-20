package internal

import "github.com/hypay-id/backend-dashboard-hypay/internal/entity"

type MerchantCallbackItf interface {
	SendCallbackAdptr(url string, transactionEntity entity.PaymentDetailMerchantProvider, transactionStatusLogLatest entity.TransactionStatusLogs, merchantSecret string) (interface{}, error)
}
