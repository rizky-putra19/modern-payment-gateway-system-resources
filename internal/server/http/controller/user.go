package controller

import (
	"fmt"
	"net/http"

	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/labstack/echo/v4"
)

func (ctrl *Controller) Authentication(c echo.Context) error {
	var payload dto.LoginPayload
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: err.Error(),
		})
	}

	if payload.Username == "" || payload.Password == "" {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "username or password can't be empty",
		})
	}

	authResponse, err := ctrl.userService.Login(payload)
	if err != nil {
		if err.Error() == "invalid password" || err.Error() == "user not found" {
			return c.JSON(http.StatusBadRequest, authResponse)
		}

		return c.JSON(http.StatusUnprocessableEntity, authResponse)
	}

	return c.JSON(http.StatusOK, authResponse)
}

func (ctrl *Controller) GetListUserMerchants(c echo.Context) error {
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
			ResponseMessage: "merchant id on query param is mandatory",
		})
	}

	listUser, err := ctrl.userService.GetListUserMerchantSvc(merchantId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listUser)
	}

	return c.JSON(http.StatusOK, listUser)
}

func (ctrl *Controller) InviteUserMerchantCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	roleName := c.Get("roleName").(string)
	var payload dto.InviteMerchantUserDto

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
			ResponseMessage: "only admin can invite merchant user",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.MerchantId == "" ||
		payload.Email == "" ||
		payload.RolesId == 0 {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant id, email, and roles id is mandatory",
		})
	}

	if !helper.IsValidEmail(payload.Email) {
		return c.JSON(http.StatusUnprocessableEntity, dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: fmt.Sprintf("%v it's not valid email", payload.Email),
		})
	}

	inviteUserResp, err := ctrl.userService.InviteUserMerchantSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, inviteUserResp)
	}

	return c.JSON(http.StatusOK, inviteUserResp)
}

func (ctrl *Controller) GetUserInformationCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	userInformationRes, err := ctrl.userService.GetUserInformationsSvc(username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, userInformationRes)
	}

	return c.JSON(http.StatusOK, userInformationRes)
}

func (ctrl *Controller) UpdatePinPasswordCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)
	var payload dto.UpdatePassOrPinDto

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid request payload",
		})
	}

	if payload.Password == nil && payload.Pin == nil {
		return c.JSON(http.StatusBadRequest, dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "need to input password or pin",
		})
	}

	payload.Username = username
	updateRes, err := ctrl.userService.UpdatePasswordOrPinSvc(payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, updateRes)
	}

	return c.JSON(http.StatusOK, updateRes)
}

func (ctrl *Controller) GetMerchantRolesCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	rolesList, err := ctrl.userService.GetListRolesSvc()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, rolesList)
	}

	return c.JSON(http.StatusOK, rolesList)
}

func (ctrl *Controller) GetMerchantListUserCtrl(c echo.Context) error {
	userType := c.Get("userType").(string)
	username := c.Get("username").(string)

	// blocked merchant user for further access
	if userType != constant.UserMerchant {
		return c.JSON(http.StatusBadGateway, dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "only merchant can access this endpoint",
		})
	}

	listUser, err := ctrl.userService.GetListMerchantUserSvc(username)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, listUser)
	}

	return c.JSON(http.StatusOK, listUser)
}
