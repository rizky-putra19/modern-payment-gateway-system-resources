package controller

import (
	"net/http"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/labstack/echo/v4"
)

func (ctrl *Controller) GetListTransaction(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	var params dto.QueryParams

	params.MinDate = c.QueryParam("minDate")
	params.MaxDate = c.QueryParam("maxDate")
	params.AmountMax = c.QueryParam("amountMax")
	params.AmountMin = c.QueryParam("amountMin")
	params.PaymentMethod = c.QueryParam("paymentMethod")
	params.PayType = c.QueryParam("payType")
	params.Status = c.QueryParam("status")
	params.RequestMethod = c.QueryParam("requestMethod")
	params.PayChannel = c.QueryParam("payChannel")
	params.MerchantName = c.QueryParam("merchantName")
	params.ProviderName = c.QueryParam("providerName")
	params.Search = c.QueryParam("search")
	params.Page = c.QueryParam("page")
	params.PageSize = c.QueryParam("pageSize")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	params.Username = username
	transactionList, err := ctrl.transactionService.GetTransactionList(params)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "unprocessable entity",
		})
	}

	return c.JSON(http.StatusOK, transactionList)
}

func (ctrl *Controller) GetPaymentDetailProviderMerchant(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	req := dto.GetPaymentDetailsRequest{
		PaymentId: paymentId,
		Username:  username,
	}

	detailData, err := ctrl.transactionService.GetPaymentDetailMerchantProvider(req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, detailData)
}

func (ctrl *Controller) GetPaymentDetailAccountInformation(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	req := dto.GetPaymentDetailsRequest{
		PaymentId: paymentId,
		Username:  username,
	}

	accountDetailData, err := ctrl.transactionService.GetPaymentDetailAccountInformation(req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, accountDetailData)
	}

	return c.JSON(http.StatusOK, accountDetailData)
}

func (ctrl *Controller) GetPaymentDetailFee(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	feeData, err := ctrl.transactionService.GetPaymentDetailFee(paymentId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, feeData)
	}

	return c.JSON(http.StatusOK, feeData)
}

func (ctrl *Controller) GetPaymentDetailProviderConfirmDetail(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	detailData, err := ctrl.transactionService.GetPaymentDetailConfirmData(paymentId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, detailData)
	}

	return c.JSON(http.StatusOK, detailData)
}

func (ctrl *Controller) GetPaymentDetailStatusChangeLogs(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	detailData, err := ctrl.transactionService.GetStatusChangeLog(paymentId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, detailData)
	}

	return c.JSON(http.StatusOK, detailData)
}

func (ctrl *Controller) GetPaymentDetailTransactions(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentId := c.QueryParam("paymentId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	transactions, err := ctrl.transactionService.GetPaymentDetailCapitalFlow(paymentId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, transactions)
	}

	return c.JSON(http.StatusOK, transactions)
}

func (ctrl *Controller) UpdateStatusTransaction(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.UpdateStatusTransaction

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

	if payload.PaymentId == "" || payload.Status == "" || payload.Notes == "" || payload.Pin == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "payment id, status, notes, and pin is mandatory",
		})
	}

	status, err := ctrl.transactionService.UpdateStatusTransaction(payload.PaymentId, payload.Status, username, payload.Notes, payload.Pin)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, status)
	}

	return c.JSON(http.StatusOK, status)
}

func (ctrl *Controller) GetListFilter(c echo.Context) error {
	listFilter, err := ctrl.transactionService.GetListFilterSvc()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listFilter)
	}

	return c.JSON(http.StatusOK, listFilter)
}

func (ctrl *Controller) CreateMerchantExportCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.CreateMerchantExportReqDto

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
			ResponseMessage: "only admin and finance can do export",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.ExportType == "" ||
		payload.MaxDate == "" ||
		payload.MinDate == "" ||
		payload.MerchantId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "all payload is mandatory",
		})
	}

	payload.UserType = userType
	exportRes, err := ctrl.transactionService.CreateMerchantExportSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, exportRes)
	}

	return c.JSON(http.StatusOK, exportRes)
}

func (ctrl *Controller) GetListMerchantExportCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	exportStatus := c.QueryParam("exportStatus")
	merchants := c.QueryParam("merchant")
	exportType := c.QueryParam("exportType")
	minDate := c.QueryParam("minDate")
	maxDate := c.QueryParam("maxDate")
	search := c.QueryParam("search")

	params := dto.GetListMerchantExportFilter{
		ExportStatus: exportStatus,
		Search:       search,
		Merchants:    merchants,
		MinDate:      minDate,
		MaxDate:      maxDate,
		ExportType:   exportType,
	}

	exportList, err := ctrl.transactionService.GetListMerchantExportSvc(params)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, exportList)
	}

	return c.JSON(http.StatusOK, exportList)
}

func (ctrl *Controller) GetListFilterExportCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	listFilter, err := ctrl.transactionService.GetListFilterExportSvc()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listFilter)
	}

	return c.JSON(http.StatusOK, listFilter)
}

func (ctrl *Controller) GetListFilterMerchantReportCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	listFilter, err := ctrl.transactionService.GetListFilterExportSvc()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listFilter)
	}

	return c.JSON(http.StatusOK, listFilter)
}

func (ctrl *Controller) CreateInternalExportCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.CreateMerchantExportReqDto

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
			ResponseMessage: "only admin and finance can do export",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.ExportType == "" ||
		payload.MaxDate == "" ||
		payload.MinDate == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "all payload is mandatory",
		})
	}

	payload.UserType = userType
	exportRes, err := ctrl.transactionService.CreateInternalExportSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, exportRes)
	}

	return c.JSON(http.StatusOK, exportRes)
}

func (ctrl *Controller) GetListInternalExport(c echo.Context) error {
	userType := c.Get("userType").(string)

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
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

	exportList, err := ctrl.transactionService.GetListInternalExportSvc(params)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, exportList)
	}

	return c.JSON(http.StatusOK, exportList)
}

func (ctrl *Controller) GetDetailTransactionCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchants can access this endpoint",
		})
	}

	transactionDetail, err := ctrl.transactionService.GetTransactionDetailSvc(c.QueryParam("transactionId"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, transactionDetail)
	}

	return c.JSON(http.StatusOK, transactionDetail)
}

func (ctrl *Controller) CountDisbursementTotalAmountCtrl(c echo.Context) error {
	username := c.Get("username").(string)
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.CountDisbursementTotalAmountDto

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
			ResponseMessage: "only admin and finance can do export",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.Amount == 0 {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "amount is mandatory",
		})
	}

	payload.Username = username
	totalAmountResp, err := ctrl.transactionService.CountDisbursementTotalAmountSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, totalAmountResp)
	}

	return c.JSON(http.StatusOK, totalAmountResp)
}

func (ctrl *Controller) CreateMerchantReportCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	username := c.Get("username").(string)
	var payload dto.CreateReportMerchantReqDto

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
			ResponseMessage: "only admin and finance can create report",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.ExportType == "" ||
		payload.MaxDate == "" ||
		payload.MinDate == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "all payload is mandatory",
		})
	}

	payload.UserType = userType
	payload.Username = username
	reportRes, err := ctrl.transactionService.CreateReportMerchantSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, reportRes)
	}

	return c.JSON(http.StatusOK, reportRes)
}
