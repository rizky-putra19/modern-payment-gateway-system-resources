package controller

import (
	"net/http"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"github.com/labstack/echo/v4"
)

func (ctrl *Controller) GetListProvidersCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	paymentMethod := c.QueryParam("paymentMethod")
	search := c.QueryParam("search")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	providerList, err := ctrl.providerService.GetListProvidersSvc(paymentMethod, search)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, providerList)
	}

	return c.JSON(http.StatusOK, providerList)
}

func (ctrl *Controller) GetProviderAnalyticsCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	minDate := c.QueryParam("minDate")
	maxdate := c.QueryParam("maxDate")
	providerId := c.QueryParam("providerId")

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if providerId == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "provider id is mandatory",
		})
	}

	payload := dto.GetProviderAnalyticsDtoReq{
		MinDate:    minDate,
		MaxDate:    maxdate,
		ProviderId: converter.ToInt(providerId),
	}

	providerAnalyticsRes, err := ctrl.providerService.GetProviderAnalyticsSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, providerAnalyticsRes)
	}

	return c.JSON(http.StatusOK, providerAnalyticsRes)
}

func (ctrl *Controller) GetListProviderPaychannelCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	providerInterfacesId := converter.ToInt(c.QueryParam("providerInterfacesId"))

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	if providerInterfacesId == 0 {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "provider interfaces id is mandatory",
		})
	}

	listPaychannel, err := ctrl.providerService.GetListProviderPaychannelSvc(providerInterfacesId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listPaychannel)
	}

	return c.JSON(http.StatusOK, listPaychannel)
}

func (ctrl *Controller) GetListProviderChannelAllCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	var params dto.QueryParams

	// blocked merchant user for further access
	if userType != constant.UserOperation {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only operations can access this endpoint",
		})
	}

	params.ProviderName = c.QueryParam("providers")
	params.Status = c.QueryParam("status")
	params.PayType = c.QueryParam("paymentType")
	params.PaymentMethod = c.QueryParam("paymentMethod")
	params.Search = c.QueryParam("search")

	listPaychannelResp, err := ctrl.providerService.GetListProviderChannelAllSvc(params)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listPaychannelResp)
	}

	return c.JSON(http.StatusOK, listPaychannelResp)
}

func (ctrl *Controller) JackDisbursementCallbackCtrl(c echo.Context) error {
	var req dto.CreateDisbursementRequestResponseData

	// convert response to struct
	err := c.Bind(&req)
	if err != nil {
		slog.Infof("JACK http-request /payOutCallback [end] [error] invalid request body (%v)", err.Error())
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  http.StatusUnprocessableEntity,
			"message": "invalid request body",
		})
	}

	amountInt := converter.FromStringToIntAmount(req.Destination.Amount)
	if amountInt < 10000 || amountInt > 500000000 {
		slog.Infof("JACK %v http-request /payOutCallback [end] [error] invalid amount (%v)", req.ReferenceID, req.Destination.Amount)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "invalid amount",
		})
	}

	_, err = ctrl.transactionService.JackDisbursementCallbackHandlingSvc(req)
	if err != nil {
		slog.Infof("JACK %v http-request /payOutCallback [end] [error] internal server error (%v)", req.ReferenceID, err.Error())
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  http.StatusUnprocessableEntity,
			"message": "interval server error",
		})
	}

	slog.Infof("JACK http-request /payOutCallback [end] [success]")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "ok",
	})
}
