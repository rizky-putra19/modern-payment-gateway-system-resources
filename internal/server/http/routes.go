package http

import (
	"github.com/hypay-id/backend-dashboard-hypay/internal/server/http/controller"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, ctrl *controller.Controller) {
	e.GET("/", ctrl.ReturnOK)

	// dashboard login
	u := e.Group("/v1")
	u.POST("/login", ctrl.Authentication)

	// operations dashboard
	ops := e.Group("/operation-dashboard/v1")
	// PATCH method
	ops.PATCH("/update-status", ctrl.AuthMiddleware(ctrl.UpdateStatusTransaction))
	ops.PATCH("/merchant-update-status", ctrl.AuthMiddleware(ctrl.UpdateMerchantStatusCtrl))
	ops.PATCH("/merchant-paychannel-update-status", ctrl.AuthMiddleware(ctrl.UpdateStatusMerchantPaychannel))
	ops.PATCH("/update-fee-limit", ctrl.AuthMiddleware(ctrl.UpdateLimitOrFeeCtrl))
	ops.PATCH("/update-fee-limit-interface-pchannel", ctrl.AuthMiddleware(ctrl.UpdateLimitFeeInterfacePchannelCtrl))

	// GET method
	ops.GET("/transaction-list", ctrl.AuthMiddleware(ctrl.GetListTransaction))
	ops.GET("/payment-details/merchant-provider", ctrl.AuthMiddleware(ctrl.GetPaymentDetailProviderMerchant))
	ops.GET("/payment-details/account-information", ctrl.AuthMiddleware(ctrl.GetPaymentDetailAccountInformation))
	ops.GET("/payment-details/fee", ctrl.AuthMiddleware(ctrl.GetPaymentDetailFee))
	ops.GET("/payment-details/provider-confirmation-detail", ctrl.AuthMiddleware(ctrl.GetPaymentDetailProviderConfirmDetail))
	ops.GET("/payment-details/status-change-logs", ctrl.AuthMiddleware(ctrl.GetPaymentDetailStatusChangeLogs))
	ops.GET("/payment-details/transactions", ctrl.AuthMiddleware(ctrl.GetPaymentDetailTransactions))
	ops.GET("/payment-details/latest-merchant-callback", ctrl.AuthMiddleware(ctrl.GetPaymentDetailLatestCallback))
	ops.GET("/list-merchant-callbacks", ctrl.AuthMiddleware(ctrl.GetListMerchantCallback))
	ops.GET("/list-merchant-callbacks/callback-attempt", ctrl.AuthMiddleware(ctrl.GetCallbackAttempts))
	ops.GET("/list-manual-payment", ctrl.AuthMiddleware(ctrl.GetListManualPayment))
	ops.GET("/manual-payment-detail", ctrl.AuthMiddleware(ctrl.GetManualPaymentDetailCtrl))
	ops.GET("/list-filter", ctrl.AuthMiddleware(ctrl.GetListFilter))
	ops.GET("/list-merchant", ctrl.AuthMiddleware(ctrl.GetListMerchantWithFilterCtrl))
	ops.GET("/merchant-balance", ctrl.AuthMiddleware(ctrl.GetMerchantBalanceCtrl))
	ops.GET("/merchant-analytics", ctrl.AuthMiddleware(ctrl.GetMerchantAnalyticsCtrl))
	ops.GET("/merchant-paychannel", ctrl.AuthMiddleware(ctrl.GetListMerchantPaychannleCtrl))
	ops.GET("/merchant-paychannel-analytics", ctrl.AuthMiddleware(ctrl.GetMerchantPaychannelAnalyticsCtrl))
	ops.GET("/list-capital-transaction", ctrl.AuthMiddleware(ctrl.GetListCapitalFlowTransactionCtrl))
	ops.GET("/paychannel-tier", ctrl.AuthMiddleware(ctrl.GetPaychannelTierCtrl))
	ops.GET("/list-merchant-account", ctrl.AuthMiddleware(ctrl.GetListMerchantAccountCtrl))
	ops.GET("/routed-paychannel", ctrl.AuthMiddleware(ctrl.GetRoutedPaychannelCtrl))
	ops.GET("/paychannel-payment-operators", ctrl.AuthMiddleware(ctrl.GetPaymentOperatorsCtrl))
	ops.GET("/aggregated-paychannel", ctrl.AuthMiddleware(ctrl.GetAggregatedPaychannelCtrl))
	ops.GET("/list-payment-method", ctrl.AuthMiddleware(ctrl.GetListChannelCreateMerchantPaychannelCtrl))
	ops.GET("/active-available-paychannel", ctrl.AuthMiddleware(ctrl.GetActiveAvailableChannel))
	ops.GET("/get-list-providers", ctrl.AuthMiddleware(ctrl.GetListProvidersCtrl))
	ops.GET("/provider-analytics", ctrl.AuthMiddleware(ctrl.GetProviderAnalyticsCtrl))
	ops.GET("/get-list-provider-paychannel", ctrl.AuthMiddleware(ctrl.GetListProviderPaychannelCtrl))
	ops.GET("/get-list-merchant-export", ctrl.AuthMiddleware(ctrl.GetListMerchantExportCtrl))
	ops.GET("/get-list-internal-export", ctrl.AuthMiddleware(ctrl.GetListInternalExport))
	ops.GET("/get-list-filter-export", ctrl.AuthMiddleware(ctrl.GetListFilterExportCtrl))
	ops.GET("/get-merchant-information", ctrl.AuthMiddleware(ctrl.GetMerchantInformationCtrl))
	ops.GET("/get-roles", ctrl.AuthMiddleware(ctrl.GetRolesCtrl))
	ops.GET("/get-users-merchant", ctrl.AuthMiddleware(ctrl.GetListUserMerchants))
	ops.GET("/get-list-provider-channel-all", ctrl.AuthMiddleware(ctrl.GetListProviderChannelAllCtrl))
	ops.GET("/provider-channel-analytics", ctrl.AuthMiddleware(ctrl.GetProviderChannelAnalyticsCtrl))
	ops.GET("/get-list-pchannel-operators", ctrl.AuthMiddleware(ctrl.GetProviderChannelOperatorsCtrl))
	ops.GET("/get-routed-provider-channel", ctrl.AuthMiddleware(ctrl.GetListRoutedProviderChannelCtrl))

	// POST Method
	ops.POST("/top-up", ctrl.AuthMiddleware(ctrl.TopUpMerchantCtrl))
	ops.POST("/hold-balance", ctrl.AuthMiddleware(ctrl.HoldBalanceCtrl))
	ops.POST("/add-settlement", ctrl.AuthMiddleware(ctrl.AddSettlementCtrl))
	ops.POST("/balance-transfer", ctrl.AuthMiddleware(ctrl.BalanceTransferCtrl))
	ops.POST("/send-callback", ctrl.AuthMiddleware(ctrl.SendCallbackCtrl))
	ops.POST("/payout-settlement", ctrl.AuthMiddleware(ctrl.SendPayoutSettlementCtrl))
	ops.POST("/reverse-manual-payment", ctrl.AuthMiddleware(ctrl.ReverseManualPaymentCtrl))
	ops.POST("/create-merchant", ctrl.AuthMiddleware(ctrl.CreateMerchantCtrl))
	ops.POST("/add-segment", ctrl.AuthMiddleware(ctrl.AddSegmentCtrl))
	ops.POST("/add-channel", ctrl.AuthMiddleware(ctrl.AddChannelCtrl))
	ops.POST("/routing-paychannel", ctrl.AuthMiddleware(ctrl.AddRoutingPaychannelCtrl))
	ops.POST("/merchant-export", ctrl.AuthMiddleware(ctrl.CreateMerchantExportCtrl))
	ops.POST("/internal-export", ctrl.AuthMiddleware(ctrl.CreateInternalExportCtrl))
	ops.POST("/display-api-key", ctrl.AuthMiddleware(ctrl.DisplaySecretKeyCtrl))
	ops.POST("/generate-api-key-merchant", ctrl.AuthMiddleware(ctrl.GenerateSecretKeyCtrl))
	ops.POST("/invite-user-merchant", ctrl.AuthMiddleware(ctrl.InviteUserMerchantCtrl))
	ops.POST("/add-operator-channel", ctrl.AuthMiddleware(ctrl.AddOperatorProviderChannelCtrl))

	mrn := e.Group("/merchant-dashboard/v1")

	// get method
	mrn.GET("/home-analytics", ctrl.AuthMiddleware(ctrl.HomeAnalyticsCtrl))
	mrn.GET("/transaction-in-list", ctrl.AuthMiddleware(ctrl.GetTransactionInListCtrl))
	mrn.GET("/transaction-detail", ctrl.AuthMiddleware(ctrl.GetDetailTransactionCtrl))
	mrn.GET("/transaction-out-list", ctrl.AuthMiddleware(ctrl.GetTransactionOutListCtrl))
	mrn.GET("/get-other-transactions", ctrl.AuthMiddleware(ctrl.GetOtherTransactionListCtrl))
	mrn.GET("/get-detail-other-transaction", ctrl.AuthMiddleware(ctrl.GetDetailOtherTransactionsCtrl))
	mrn.GET("/get-list-callback-merchants", ctrl.AuthMiddleware(ctrl.GetListMerchantCallbackCtrl))
	mrn.GET("/get-merchant-account-balance", ctrl.AuthMiddleware(ctrl.GetMerchantAccountBalanceCtrl))
	mrn.GET("/get-bank-list-disbursement", ctrl.AuthMiddleware(ctrl.GetBankListForDisbursementCtrl))
	mrn.GET("/get-report-list", ctrl.AuthMiddleware(ctrl.GetReportListCtrl))
	mrn.GET("/get-list-filter-report", ctrl.AuthMiddleware(ctrl.GetListFilterMerchantReportCtrl))

	// post method
	mrn.POST("/resend-callback", ctrl.AuthMiddleware(ctrl.ResendCallbackMerchantCtrl))
	mrn.POST("/disbursement", ctrl.AuthMiddleware(ctrl.DisbursementCtrl))
	mrn.POST("/count-disbursement", ctrl.AuthMiddleware(ctrl.CountDisbursementTotalAmountCtrl))
	mrn.POST("/provider-jack/disbursement", ctrl.JackDisbursementCallbackCtrl)
	mrn.POST("/create-report", ctrl.AuthMiddleware(ctrl.CreateMerchantReportCtrl))
}
