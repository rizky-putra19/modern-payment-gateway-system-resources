package controller

import (
	"net/http"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
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
