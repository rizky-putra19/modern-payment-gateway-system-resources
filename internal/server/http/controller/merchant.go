package controller

import (
	"net/http"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/labstack/echo/v4"
)

func (ctrl *Controller) GetCallbackAttempts(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	callbackAttempt, err := ctrl.merchantService.GetCallbackAttemptsByPaymentId(paymentId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, callbackAttempt)
}

func (ctrl *Controller) GetPaymentDetailLatestCallback(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	callbackData, err := ctrl.merchantService.GetLatestMerchantCallback(paymentId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, callbackData)
	}

	return c.JSON(http.StatusOK, callbackData)
}

func (ctrl *Controller) GetListMerchantCallback(c echo.Context) error {
	userType := c.Get("userType").(string)
	var queryParamsMerchantCallback dto.QueryParamsMerchantCallback

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	queryParamsMerchantCallback.PayType = c.QueryParam("payType")
	queryParamsMerchantCallback.MerchantName = c.QueryParam("merchantName")
	queryParamsMerchantCallback.CallbackStatus = c.QueryParam("callbackStatus")
	queryParamsMerchantCallback.Page = c.QueryParam("page")
	queryParamsMerchantCallback.PageSize = c.QueryParam("pageSize")
	queryParamsMerchantCallback.Search = c.QueryParam("search")
	queryParamsMerchantCallback.MaxDate = c.QueryParam("maxDate")
	queryParamsMerchantCallback.MinDate = c.QueryParam("minDate")

	callbackList, err := ctrl.merchantService.GetListMerchantCallback(queryParamsMerchantCallback)
	if err != nil {
		return c.JSON(http.StatusOK, callbackList)
	}

	return c.JSON(http.StatusOK, callbackList)
}

func (ctrl *Controller) GetListManualPayment(c echo.Context) error {
	userType := c.Get("userType").(string)
	var queryParamsManualPayment dto.QueryParamsManualPayment

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	queryParamsManualPayment.MaxDate = c.QueryParam("maxDate")
	queryParamsManualPayment.MinDate = c.QueryParam("minDate")
	queryParamsManualPayment.MerchantName = c.QueryParam("merchantName")
	queryParamsManualPayment.Status = c.QueryParam("status")
	queryParamsManualPayment.Page = c.QueryParam("page")
	queryParamsManualPayment.PageSize = c.QueryParam("pageSize")
	queryParamsManualPayment.ReasonName = c.QueryParam("reasonName")
	queryParamsManualPayment.Search = c.QueryParam("search")
	queryParamsManualPayment.AmountMax = c.QueryParam("amountMax")
	queryParamsManualPayment.AmountMin = c.QueryParam("amountMin")

	listManualPayment, err := ctrl.merchantService.GetListManualPaymentWithFilter(queryParamsManualPayment)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listManualPayment)
	}

	return c.JSON(http.StatusOK, listManualPayment)
}

func (ctrl *Controller) GetListMerchantWithFilterCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	var queryParams dto.QueryParams

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	queryParams.Status = c.QueryParam("status")
	queryParams.MerchantName = c.QueryParam("merchantName")
	queryParams.Search = c.QueryParam("search")

	listMerchant, err := ctrl.merchantService.GetListMerchantWithFilterSvc(queryParams)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listMerchant)
	}

	return c.JSON(http.StatusOK, listMerchant)
}

func (ctrl *Controller) GetManualPaymentDetailCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	// proceed on service handle
	manualDetailResp, err := ctrl.merchantService.GetDetailManualPaymentSvc(paymentId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, manualDetailResp)
	}

	return c.JSON(http.StatusOK, manualDetailResp)
}

func (ctrl *Controller) TopUpMerchantCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.AdjustBalanceReqPayload

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin && roleName != constant.RoleNameFinance {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin and finance can do top up",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantId == "" || payload.Notes == "" || payload.Amount == 0 || payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id, notes, amount, and pin is mandatory",
		})
	}

	payload.Username = username
	topResp, err := ctrl.merchantService.TopUpMerchantSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, topResp)
	}

	return c.JSON(http.StatusOK, topResp)
}

func (ctrl *Controller) HoldBalanceCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.AdjustBalanceReqPayload

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin && roleName != constant.RoleNameFinance {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin and finance can do hold balance",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantId == "" || payload.Notes == "" || payload.Amount == 0 || payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id, notes, amount, and pin is mandatory",
		})
	}

	payload.Username = username
	holdBalanceResp, err := ctrl.merchantService.HoldBalanceSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, holdBalanceResp)
	}

	return c.JSON(http.StatusOK, holdBalanceResp)
}

func (ctrl *Controller) AddSettlementCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.AdjustBalanceReqPayload

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin && roleName != constant.RoleNameFinance {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin and finance can do top up",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantId == "" || payload.Notes == "" || payload.Amount == 0 || payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id, notes, amount, and pin is mandatory",
		})
	}

	payload.Username = username
	settlementResp, err := ctrl.merchantService.SettlementBalanceSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, settlementResp)
	}

	return c.JSON(http.StatusOK, settlementResp)
}

func (ctrl *Controller) BalanceTransferCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.BalanceTrfReqPayload

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin && roleName != constant.RoleNameFinance {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin and finance can do top up",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.AccountFrom.MerchantId == "" ||
		payload.AccountTo.MerchantId == "" ||
		payload.Amount == 0 ||
		payload.Pin == "" ||
		payload.Notes == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "all payload is mandatory",
		})
	}

	payload.Username = username
	balanceTrfResp, err := ctrl.merchantService.BalanceTransferSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, balanceTrfResp)
	}

	return c.JSON(http.StatusOK, balanceTrfResp)
}

func (ctrl *Controller) SendCallbackCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.SendCallbackReqPayload

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can do send callback",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.PaymentId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "transaction id is mandatory",
		})
	}

	sendCallbackResp, err := ctrl.merchantService.SendCallbackSvc(payload.PaymentId, username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, sendCallbackResp)
	}

	return c.JSON(http.StatusOK, sendCallbackResp)
}

func (ctrl *Controller) SendPayoutSettlementCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.AdjustBalanceReqPayload

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin && roleName != constant.RoleNameFinance {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin and finance can do out settlement",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantId == "" || payload.Notes == "" || payload.Amount == 0 || payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id, notes, amount, and pin is mandatory",
		})
	}

	payload.Username = username
	outSettlementResp, err := ctrl.merchantService.PayoutSettlementSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, outSettlementResp)
	}

	return c.JSON(http.StatusOK, outSettlementResp)
}

func (ctrl *Controller) ReverseManualPaymentCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)
	roleName := c.Get("roleName").(string)

	var payload dto.UpdateStatusTransaction

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin && roleName != constant.RoleNameFinance {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin and finance can do top up",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.PaymentId == "" || payload.Notes == "" || payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "payment id, status, notes, and pin is mandatory",
		})
	}

	reverseResp, err := ctrl.merchantService.ReverseManualPaymentSvc(payload, username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, reverseResp)
	}

	return c.JSON(http.StatusOK, reverseResp)
}

func (ctrl *Controller) GetMerchantBalanceCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	merchantId := c.QueryParam("merchantId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	merchantAccountResp, err := ctrl.merchantService.GetMerchantBalanceSvc(merchantId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, merchantAccountResp)
	}

	return c.JSON(http.StatusOK, merchantAccountResp)
}

func (ctrl *Controller) CreateMerchantCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.CreateMerchantDtoReq

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can create merchant",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantName == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant name is mandatory",
		})
	}

	createMerchantResp, err := ctrl.merchantService.CreateMerchantSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, createMerchantResp)
	}

	return c.JSON(http.StatusOK, createMerchantResp)
}

func (ctrl *Controller) GetMerchantAnalyticsCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	minDate := c.QueryParam("minDate")
	maxdate := c.QueryParam("maxDate")
	merchandId := c.QueryParam("merchantId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if merchandId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id is mandatory",
		})
	}

	payload := dto.GetMerchantAnalyticsDtoReq{
		MinDate:    minDate,
		MaxDate:    maxdate,
		MerchantId: merchandId,
	}

	merchantAnalyticsRes, err := ctrl.merchantService.GetMerchantAnalyticsSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, merchantAnalyticsRes)
	}

	return c.JSON(http.StatusOK, merchantAnalyticsRes)
}

func (ctrl *Controller) UpdateMerchantStatusCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)

	var payload dto.AccountData

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can update status",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id is mandatory",
		})
	}

	merchantUpdateStatusRes, err := ctrl.merchantService.MerchantUpdateStatusSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, merchantUpdateStatusRes)
	}

	return c.JSON(http.StatusOK, merchantUpdateStatusRes)
}

func (ctrl *Controller) GetListMerchantPaychannleCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	merchandId := c.QueryParam("merchantId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if merchandId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id is mandatory",
		})
	}

	merchantPaychannelList, err := ctrl.merchantService.GetMerchantPaychannelSvc(merchandId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, merchantPaychannelList)
	}

	return c.JSON(http.StatusOK, merchantPaychannelList)
}

func (ctrl *Controller) GetListCapitalFlowTransactionCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	var params dto.QueryParams

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	params.MerchantId = c.QueryParam("merchantId")
	params.PayType = c.QueryParam("direction")
	params.PaymentMethod = c.QueryParam("transactionType")
	params.Status = c.QueryParam("status")
	params.MinDate = c.QueryParam("minDate")
	params.MaxDate = c.QueryParam("maxDate")
	params.Search = c.QueryParam("search")
	params.Page = c.QueryParam("page")
	params.PageSize = c.QueryParam("pageSize")

	if params.MerchantId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id is mandatory",
		})
	}

	listTransactionCapitalFlow, err := ctrl.merchantService.GetListCapitalFlowTransactionSvc(params)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listTransactionCapitalFlow)
	}

	return c.JSON(http.StatusOK, listTransactionCapitalFlow)
}

func (ctrl *Controller) GetPaychannelTierCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	merchantPaychannelId := c.QueryParam("merchantPaychannelId")
	intMerchantPaychannelId := converter.ToInt(merchantPaychannelId)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can add tier paychannel",
		})
	}

	tierList, err := ctrl.merchantService.GetListTierPaychannelSvc(intMerchantPaychannelId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tierList)
	}

	return c.JSON(http.StatusOK, tierList)
}

func (ctrl *Controller) GetListMerchantAccountCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	var queryParams dto.QueryParams

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	queryParams.Status = c.QueryParam("status")
	queryParams.MerchantName = c.QueryParam("merchantName")
	queryParams.Search = c.QueryParam("search")

	listAccount, err := ctrl.merchantService.GetListMerchantAccountSvc(queryParams)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listAccount)
	}

	return c.JSON(http.StatusOK, listAccount)
}

func (ctrl *Controller) GetMerchantPaychannelAnalyticsCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	minDate := c.QueryParam("minDate")
	maxdate := c.QueryParam("maxDate")
	merchantPaychannelId := c.QueryParam("merchantPaychannelId")
	intMerchantPaychannelId := converter.ToInt(merchantPaychannelId)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if merchantPaychannelId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id is mandatory",
		})
	}

	payload := dto.GetMerchantAnalyticsDtoReq{
		MinDate:              minDate,
		MaxDate:              maxdate,
		MerchantPaychannelId: intMerchantPaychannelId,
	}

	merchantPaychannelAnalyticsData, err := ctrl.merchantService.GetMerchantPaychannelAnalyticsSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, merchantPaychannelAnalyticsData)
	}

	return c.JSON(http.StatusOK, merchantPaychannelAnalyticsData)
}

func (ctrl *Controller) GetRoutedPaychannelCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	merchantPaychannelId := c.QueryParam("merchantPaychannelId")
	intMerchantPaychannelId := converter.ToInt(merchantPaychannelId)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	routedPaychannelList, err := ctrl.merchantService.GetRoutedPaychannelSvc(intMerchantPaychannelId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, routedPaychannelList)
	}

	return c.JSON(http.StatusOK, routedPaychannelList)
}

func (ctrl *Controller) GetPaymentOperatorsCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	routedChannelId := c.QueryParam("routedChannelId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	paymentOperatorList, err := ctrl.merchantService.GetPaymentOperatorsSvc(routedChannelId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, paymentOperatorList)
	}

	return c.JSON(http.StatusOK, paymentOperatorList)
}

func (ctrl *Controller) UpdateLimitOrFeeCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.AdjustLimitOrFeePayload

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can adjust fee and limit",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantPaychannelId == 0 {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Merchant paychannel id is mandatory",
		})
	}

	if payload.FeeType != nil {
		if !helper.StringInSlice(*payload.FeeType, constant.FeeType) {
			return c.JSON(http.StatusBadRequest, dto.ResponseDto{
				ResponseCode:    http.StatusBadRequest,
				ResponseMessage: "Fee type only FIXED_FEE and PERCENTAGE",
			})
		}
	}

	payload.Username = username
	adjustResp, err := ctrl.merchantService.UpdateLimitOrFeeMerchantPaychannelSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, adjustResp)
	}

	return c.JSON(http.StatusOK, adjustResp)
}

func (ctrl *Controller) GetAggregatedPaychannelCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	merchantPaychannelId := c.QueryParam("merchantPaychannelId")
	intMerchantPaychannelId := converter.ToInt(merchantPaychannelId)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can adjust fee and limit",
		})
	}

	aggregatedData, err := ctrl.merchantService.GetAggregatedPaychannelSvc(intMerchantPaychannelId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, aggregatedData)
	}

	return c.JSON(http.StatusOK, aggregatedData)
}

func (ctrl *Controller) AddSegmentCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.AddSegmentDtoReq

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can adjust fee and limit",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantPaychannelId == 0 && payload.TierName == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant paychannel id and tier name is mandatory",
		})
	}

	addSegmentResp, err := ctrl.merchantService.AddSegmentMerchantPaychannelSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, addSegmentResp)
	}

	return c.JSON(http.StatusOK, addSegmentResp)
}

func (ctrl *Controller) UpdateStatusMerchantPaychannel(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.AccountData

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can adjust fee and limit",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantPaychannelId == 0 {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant paychannel id is mandatory",
		})
	}

	updateStatusRes, err := ctrl.merchantService.ActivateOrDeactivateMerchantPaychannelSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, updateStatusRes)
	}

	return c.JSON(http.StatusOK, updateStatusRes)
}

func (ctrl *Controller) AddChannelCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	idMerchant := c.QueryParam("idMerchant")
	intIdMerchant := converter.ToInt(idMerchant)
	var payload []dto.PaymentMethodData

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can create merchant",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	createChanneResp, err := ctrl.merchantService.AddChannelSvc(intIdMerchant, payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, createChanneResp)
	}

	return c.JSON(http.StatusOK, createChanneResp)
}

func (ctrl *Controller) GetListChannelCreateMerchantPaychannelCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	idMerchant := c.QueryParam("idMerchant")
	intIdMerchant := converter.ToInt(idMerchant)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can create merchant",
		})
	}

	listPaymentMethods, err := ctrl.merchantService.GetListMerchantPaymentMethodsSvc(intIdMerchant)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listPaymentMethods)
	}

	return c.JSON(http.StatusOK, listPaymentMethods)
}

func (ctrl *Controller) AddRoutingPaychannelCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	merchantPaychannelId := c.QueryParam("merchantPaychannelId")
	intMerchantPaychannelId := converter.ToInt(merchantPaychannelId)
	var payload dto.AddPaychannelRouting

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin can create merchant",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	AddRoutingResp, err := ctrl.merchantService.AddRoutingPaychannelSvc(intMerchantPaychannelId, payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, AddRoutingResp)
	}

	return c.JSON(http.StatusOK, AddRoutingResp)
}

func (ctrl *Controller) GetActiveAvailableChannel(c echo.Context) error {
	userType := c.Get("userType").(string)
	merchantPaychannelId := c.QueryParam("merchantPaychannelId")
	intMerchantPaychannelId := converter.ToInt(merchantPaychannelId)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	ActiveAndAvailablePaychannelList, err := ctrl.merchantService.GetActiveAvailablePaychannelSvc(intMerchantPaychannelId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ActiveAndAvailablePaychannelList)
	}

	return c.JSON(http.StatusOK, ActiveAndAvailablePaychannelList)
}

func (ctrl *Controller) HomeAnalyticsCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchants can access this endpoint",
		})
	}

	payload := dto.HomeAnalyticsDto{
		MaxDate:  c.QueryParam("maxDate"),
		MinDate:  c.QueryParam("minDate"),
		Username: username,
	}

	homeAnalyticsRes, err := ctrl.merchantService.HomeAnalyticsSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, homeAnalyticsRes)
	}

	return c.JSON(http.StatusOK, homeAnalyticsRes)
}

func (ctrl *Controller) GetTransactionInListCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)
	var params dto.QueryParams

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchants can access this endpoint",
		})
	}

	params.MinDate = c.QueryParam("minDate")
	params.MaxDate = c.QueryParam("maxDate")
	params.PaymentMethod = c.QueryParam("transactionMethod")
	params.Status = c.QueryParam("status")
	params.Search = c.QueryParam("search")
	params.Page = c.QueryParam("page")
	params.PageSize = c.QueryParam("pageSize")

	params.Username = username
	transactionInList, err := ctrl.transactionService.GetTransactionInListSvc(params)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, transactionInList)
	}

	return c.JSON(http.StatusOK, transactionInList)
}

func (ctrl *Controller) GetTransactionOutListCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)
	var params dto.QueryParams

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchants can access this endpoint",
		})
	}

	params.MinDate = c.QueryParam("minDate")
	params.MaxDate = c.QueryParam("maxDate")
	params.PaymentMethod = c.QueryParam("transactionMethod")
	params.Status = c.QueryParam("status")
	params.Search = c.QueryParam("search")
	params.Page = c.QueryParam("page")
	params.PageSize = c.QueryParam("pageSize")

	params.Username = username
	transactionOutList, err := ctrl.transactionService.GetTransactionOutListSvc(params)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, transactionOutList)
	}

	return c.JSON(http.StatusOK, transactionOutList)
}

func (ctrl *Controller) ResendCallbackMerchantCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)

	var payload dto.SendCallbackReqPayload

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only user merchants can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant admin can do send callback",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.PaymentId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "transaction id is mandatory",
		})
	}

	sendCallbackResp, err := ctrl.merchantService.SendCallbackSvc(payload.PaymentId, username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, sendCallbackResp)
	}

	return c.JSON(http.StatusOK, sendCallbackResp)
}

func (ctrl *Controller) GetMerchantInformationCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	merchantId := c.QueryParam("merchantId")
	if merchantId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id is mandatory",
		})
	}

	merchantInformation, err := ctrl.merchantService.GetMerchantInformationSvc(merchantId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, merchantInformation)
	}

	return c.JSON(http.StatusOK, merchantInformation)
}

func (ctrl *Controller) DisplaySecretKeyCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	username := c.Get("username").(string)
	var payload dto.BalanceTrfReqPayload
	merchantId := c.QueryParam("merchantId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin display secret key",
		})
	}

	if merchantId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchantId on query params is mandatory",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "pin is mandatory",
		})
	}

	displayKey, err := ctrl.merchantService.DisplaySecretKeySvc(payload.Pin, username, merchantId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, displayKey)
	}

	return c.JSON(http.StatusOK, displayKey)
}

func (ctrl *Controller) GenerateSecretKeyCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	username := c.Get("username").(string)
	var payload dto.BalanceTrfReqPayload
	merchantId := c.QueryParam("merchantId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin display secret key",
		})
	}

	if merchantId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchantId on query params is mandatory",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "pin is mandatory",
		})
	}

	generateKey, err := ctrl.merchantService.GenerateSecretKeySvc(payload.Pin, username, merchantId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, generateKey)
	}

	return c.JSON(http.StatusOK, generateKey)
}

func (ctrl *Controller) GetRolesCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	listRoles, err := ctrl.userService.GetListRolesSvc()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listRoles)
	}

	return c.JSON(http.StatusOK, listRoles)
}

func (ctrl *Controller) GetOtherTransactionListCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)
	var params dto.QueryParamsManualPayment

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only user merchants can access this endpoint",
		})
	}

	params.MinDate = c.QueryParam("minDate")
	params.MaxDate = c.QueryParam("maxDate")
	params.Status = c.QueryParam("status")
	params.Page = c.QueryParam("page")
	params.PageSize = c.QueryParam("pageSize")
	params.Search = c.QueryParam("search")
	params.ReasonName = c.QueryParam("reasonName")
	params.AmountMax = c.QueryParam("amountMax")
	params.AmountMin = c.QueryParam("amountMin")
	params.Username = username

	listOtherTransactions, err := ctrl.merchantService.GetListOtherTransactionsSvc(params)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listOtherTransactions)
	}

	return c.JSON(http.StatusOK, listOtherTransactions)
}

func (ctrl *Controller) GetDetailOtherTransactionsCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only user merchants can access this endpoint",
		})
	}

	// proceed on service handle
	manualDetailResp, err := ctrl.merchantService.GetDetailManualPaymentSvc(paymentId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, manualDetailResp)
	}

	return c.JSON(http.StatusOK, manualDetailResp)
}

func (ctrl *Controller) GetListMerchantCallbackCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)
	var queryParamsMerchantCallback dto.QueryParamsMerchantCallback

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only user merchants can access this endpoint",
		})
	}

	queryParamsMerchantCallback.PayType = c.QueryParam("payType")
	queryParamsMerchantCallback.CallbackStatus = c.QueryParam("callbackStatus")
	queryParamsMerchantCallback.Page = c.QueryParam("page")
	queryParamsMerchantCallback.PageSize = c.QueryParam("pageSize")
	queryParamsMerchantCallback.Search = c.QueryParam("search")
	queryParamsMerchantCallback.MaxDate = c.QueryParam("maxDate")
	queryParamsMerchantCallback.MinDate = c.QueryParam("minDate")
	queryParamsMerchantCallback.Username = username

	callbackList, err := ctrl.merchantService.GetListCallbackMerchantSvc(queryParamsMerchantCallback)
	if err != nil {
		return c.JSON(http.StatusOK, callbackList)
	}

	return c.JSON(http.StatusOK, callbackList)
}

func (ctrl *Controller) GetMerchantAccountBalanceCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only user merchants can access this endpoint",
		})
	}

	accountBalanceRes, err := ctrl.merchantService.GetMerchantAccountBalanceSvc(username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, accountBalanceRes)
	}

	return c.JSON(http.StatusOK, accountBalanceRes)
}

func (ctrl *Controller) DisbursementCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.MerchantDisbursement

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin && roleName != constant.RoleNameFinance {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin and finance can do disbursement",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.BankAccountName == "" || payload.Amount == 0 || payload.Pin == "" || payload.BankAccountNumber == "" || payload.BankName == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "bank account name, amount, pin, bank account number, and bank name is mandatory",
		})
	}

	payload.Username = username
	disbursementResp, err := ctrl.transactionService.MerchantDisbursementSvc(payload)
	if err != nil {
		if err.Error() == "wrong pin" || err.Error() == "insufficient" {
			return c.JSON(http.StatusBadRequest, disbursementResp)
		}
		return c.JSON(http.StatusUnprocessableEntity, disbursementResp)
	}

	return c.JSON(http.StatusOK, disbursementResp)
}

func (ctrl *Controller) GetBankListForDisbursementCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	bankListDisbursementResp, err := ctrl.transactionService.GetBankListDisbursementSvc(username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, bankListDisbursementResp)
	}

	return c.JSON(http.StatusOK, bankListDisbursementResp)
}

func (ctrl *Controller) GetReportListCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	exportStatus := c.QueryParam("exportStatus")
	exportType := c.QueryParam("exportType")
	minDate := c.QueryParam("minDate")
	maxDate := c.QueryParam("maxDate")
	search := c.QueryParam("search")

	params := dto.GetListMerchantExportFilter{
		ExportStatus: exportStatus,
		Search:       search,
		MinDate:      minDate,
		MaxDate:      maxDate,
		ExportType:   exportType,
	}

	reportList, err := ctrl.transactionService.GetReportListMerchantSvc(params, username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, reportList)
	}

	return c.JSON(http.StatusOK, reportList)
}

func (ctrl *Controller) GetInformationMerchantCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	merchantInformation, err := ctrl.merchantService.GetInformationMerchantSvc(username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, merchantInformation)
	}

	return c.JSON(http.StatusOK, merchantInformation)
}

func (ctrl *Controller) DisplayMerchantKeyCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	username := c.Get("username").(string)
	var payload dto.BalanceTrfReqPayload

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin display secret key",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "pin is mandatory",
		})
	}

	displayKey, err := ctrl.merchantService.DisplayMerchantKeySvc(username, payload.Pin)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, displayKey)
	}

	return c.JSON(http.StatusOK, displayKey)
}

func (ctrl *Controller) GenerateMerchantKeyCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	username := c.Get("username").(string)
	var payload dto.BalanceTrfReqPayload

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	if roleName != constant.RoleNameAdmin {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only admin display secret key",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	generateRes, err := ctrl.merchantService.GenerateMerchantKeySvc(payload.Pin, username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, generateRes)
	}

	return c.JSON(http.StatusOK, generateRes)
}
