package internal

import (
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
)

type TransactionsReadsRepositoryItf interface {
	GetTransactionList(params dto.QueryParams) ([]entity.Transaction, dto.PaginatedResponse, error)
	GetPaymentDetailProviderMerchant(paymentId string) (entity.PaymentDetailMerchantProvider, error)
	GetPaymentDetailAccountInformation(paymentId string) ([]entity.AccountData, error)
	GetPaymentDetailPrvConfirmDetail(paymentId string) ([]entity.ProviderConfirmDetail, error)
	GetStatusChangeLogData(paymentId string) ([]entity.TransactionStatusLogs, error)
	GetCapitalFlowsWitPaymentId(paymentId string) ([]entity.CapitalFlows, error)
	GetTransactionAnalyticsRepo(params dto.GetMerchantAnalyticsDtoReq) ([]entity.PaymentDetailMerchantProvider, error)
	GetTransactionCapitalFlowRepo(params dto.QueryParams) ([]entity.TransactionCapitalFlows, dto.PaginatedResponse, error)
	GetTransactionListByMerchantPaychannelRepo(payload dto.GetMerchantAnalyticsDtoReq) ([]entity.PaymentDetailMerchantProvider, error)
	GetTransactionDataForProviderAnalyticsRepo(payload dto.GetProviderAnalyticsDtoReq) ([]entity.PaymentDetailMerchantProvider, error)
	GetListMerchantExportRepo(params dto.GetListMerchantExportFilter) ([]entity.ReportStoragesEntity, error)
	GetListMerchantReportRepo(params dto.GetListMerchantExportFilter, merchantId string) ([]entity.ReportStoragesEntity, error)
	GetTransactionInListRepo(params dto.QueryParams) ([]entity.MerchantTransactionList, dto.PaginatedResponse, error)
	GetAccountInformationByPaymentIdAccountType(paymentId string, accountType string) (entity.AccountData, error)
	GetTransactionOutListRepo(params dto.QueryParams) ([]entity.MerchantTransactionList, dto.PaginatedResponse, error)
	GetBankDataDetailRepo(bankName string) (entity.BankListDto, error)
	GetTransactionDataByProviderChannelRepo(payload dto.GetProviderAnalyticsDtoReq) ([]entity.PaymentDetailMerchantProvider, error)
	GetBankDataDetailByBankCodeRepo(bankCode string) (entity.BankListDto, error)
}

type TransactionsWritesRepositoryItf interface {
	UpdateStatus(status string, paymentId string) error
	CreateTransactionStatusLog(paymentId string, statusLog string, changeBy string, notes string, realNotes string) (int, error)
	UpdateReportStoragesByFileName(publicUrl string, fileName string, status string) error
	CreateListReportStoragesRepo(payload dto.CreateReportStorageDto) (int, error)
	CreateMerchantExportCapitalFlowRepo(payload dto.CreateMerchantExportReqDto) ([]entity.MerchantExportCapitalFlowEntity, error)
	CreateTransactionsRepo(payload dto.CreateTransactionsDto) (int, error)
	CreateAccountInformationRepo(payload dto.CreateAccountInformationDto) (int, error)
}

type MerchantReadsRepositoryItf interface {
	GetListMerchantCallbackWithFilter(params dto.QueryParamsMerchantCallback) ([]entity.ListMerchantCallback, dto.PaginatedResponse, error)
	GetListMerchantCallback(paymentId string) ([]entity.MerchantCallback, error)
	GetMerchantAccountByMerchantId(merchantId string) (entity.MerchantAccount, error)
	GetMerchantDataByMerchantId(merchantId string) (entity.Merchants, error)
	GetListManualPayment(params dto.QueryParamsManualPayment) ([]entity.ManualPayment, dto.PaginatedResponse, error)
	GetListFilter() (dto.FilterResponseDto, error)
	GetDetailManualPayment(paymentId string) ([]entity.ManualPayment, error)
	GetListMerchantWithFilterRepo(params dto.QueryParams) ([]entity.Merchants, error)
	GetMerchantPaychannelByMerchantId(merchantId string) ([]entity.MerchantPaychannel, error)
	CheckReverseStatusRepo(paymentId string) ([]entity.ManualPayment, error)
	GetListMerchantAccountRepo(params dto.QueryParams) ([]entity.ListMerchantAccountDto, error)
	GetMerchantPaychannelDetailById(id int) (entity.MerchantPaychannel, error)
	GetListRoutedPaychannelByIdMerchantPaychannelRepo(id int) ([]entity.RoutedPaychanneDto, error)
	GetBankListProviderPaymentMethodRepo(routedChannelName string) ([]entity.BankListDto, error)
	GetBankListFromProviderPaychannelRepo(routedChannelName string) ([]string, error)
	GetMerchantPaychannelByPaymentMethodId(id int) ([]entity.MerchantPaychannel, error)
	GetAggregatedPaychannelByIdRepo(id int) (entity.AggregatedPaychannelEntity, error)
	GetMerchantPaymentMethodByIdMerchantRepo(id int) ([]entity.PaymentMethods, []entity.PaymentMethods, error)
	GetActiveAndAvailableChannelRepo(merchantPaychannelId int, paymentMethodName string) ([]entity.ProviderPaychannelEntity, []entity.ProviderPaychannelEntity, error)
	GetSecretKeyByMerchantIdRepo(merchantId string) (entity.MerchantSecret, error)
	GetBankListForDisbursementRepo(routedChannelName string) ([]entity.BankListDto, error)
}

type MerchantWritesRepositoryItf interface {
	UpdateMerchantCapitalAndNotSettleBalance(notSettleBalance float64, balanceCapitalFlow float64, merchantId string) error
	UpdateMerchantCapitalAndSettleBalance(settleBalance float64, balanceCapitalFlow float64, merchantId string) error
	CreateMerchantCapitalFlow(payload dto.CreateMerchantCapitalFlowPayload) (int, error)
	UpdateMerchantHoldBalanceAndSettleBalance(settleBalance float64, holdBalance float64, merchantId string) error
	UpdateMerchantHoldBalanceAndNotSettleBalance(notSettleBalance float64, holdBalance float64, merchantId string) error
	UpdateMerchantSettlement(settleBalance float64, notSettleBalance float64, merchantId string) error
	UpdateMerchantCapitalPendingOut(pendingAmount float64, balanceCapitalFlow float64, merchantId string) error
	CreateMerchantCallback(paymentId string, callbackStatus string, paymentStatusInCallback string, callbackResult string, triggerBy string) (int, error)
	CreateMerchantRepo(merchantName string, merchantId string, merchantSecret string) (int, error)
	CreateMerchantPaymentMethodRepo(merchantId int, paymentMethodId int) (int, error)
	CreateMerchantPaychannelRepo(merchantPaymentMethodId int, segment string, fee float64, feeType string, minAmount float64, maxAmount float64, dailyLimit float64, merchantPaychannelCode string) (int, error)
	CreateMerchantAccountsRepo(merchantId string) (int, error)
	UpdateMerchantStatusRepo(merchantId string, status string) error
	UpdateMerchantPaychannelByIdRepo(payload dto.AdjustLimitOrFeePayload) error
	UpdateStatusMerchantPaychannelById(id int, status string) error
	DeleteRoutingPaychannelByMerchantPaychannelId(id int) error
	AddRoutingPaychannelRepo(merchantPaychannelId int, providerPaychannelId int) (int, error)
	UpdateMerchantSecretKeyRepo(secretKey string, merchantId string) error
	UpdateMerchantBalanceSettleAndPendingOutBalanceRepo(settleBalance float64, pendingOutBalance float64, merchantId string) error
}

type UserReadsRepositoryItf interface {
	GetUserByUsername(username string) (entity.User, error)
	GetPermissionByRoleId(id int) ([]entity.Permission, error)
	GetRolesRepo() ([]entity.RolesEntity, error)
	GetListUserByMerchantIdRepo(merchantId string) ([]entity.ListUsersEntity, error)
}

type UserWritesRepositoryItf interface {
	CreateUsersMerchantRepo(payload dto.InviteMerchantUserDto, credentials dto.EmailDataHtmlDto) (int, error)
}

type ProviderReadsRepositoryItf interface {
	CountProviderChannelByPaymentMethodRepo(paymentMethodName string) (int, error)
	CountActiveProviderChannelRepo(merchantPaychannelId int) (int, error)
	CountInterfacesProviderPaychannelByIdProvider(id int) (int, int, error)
	GetListProvidersWithFilterRepo(paymentMethod string, search string) ([]entity.ProviderListEntity, error)
	GetProviderInterfacesRepoById(id int) ([]entity.ProviderInterfacesEntity, error)
	GetListProviderPaychannelById(id int) ([]entity.InterfacePaychannelEntity, error)
	GetListProviderChannelAllRepo(params dto.QueryParams) ([]entity.ProviderPaychannelAllEntity, error)
	GetAllCredentialsRepo(providerId string, interfaceSetting string) ([]entity.ProviderCredentialsEntity, error)
	GetDetailProviderChannelById(id int) (entity.ProviderChannelDetailEntity, error)
	GetBankListProviderMethodRepo(providerChannelId int) ([]entity.BankListDto, error)
	GetBankListProviderChannelRepo(providerChannelId int) ([]entity.BankListDto, error)
	GetListRoutedProviderChannelRepo(providerChannelId int) ([]entity.ProviderRoutedChannelEntity, error)
	GetProviderBankListChannelRepo(providerChannelId int, bankListId int) (entity.ProviderPaychannelBankListEntity, error)
}

type ProviderWritesRepositoryItf interface {
	CreateProviderConfirmationDetail(source string, paymentId string, status string) (int, error)
	UpdateProviderPaychannelByIdRepo(payload dto.AdjustLimitOrFeeProviderPayload) error
	AddOperatorProviderChannelRepo(providerChannelId int, bankListId int) (int, error)
	DeleteOperatorProviderChannelRepo(providerChannelId int, bankListId int) error
	UpdateStatusProviderPaychannelRepo(id int, status string) error
}
