package service

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/hypay-id/backend-dashboard-hypay/internal"
	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
)

type Merchant struct {
	merchantRepoReads     internal.MerchantReadsRepositoryItf
	merchantRepoWrites    internal.MerchantWritesRepositoryItf
	userRepoReads         internal.UserReadsRepositoryItf
	merchantCallbackAdptr internal.MerchantCallbackItf
	transactionRepoReads  internal.TransactionsReadsRepositoryItf
	providerRepoReads     internal.ProviderReadsRepositoryItf
}

func NewMerchant(
	merchantRepoReads internal.MerchantReadsRepositoryItf,
	merchantRepoWrites internal.MerchantWritesRepositoryItf,
	userRepoReads internal.UserReadsRepositoryItf,
	adapterMerchantCallback internal.MerchantCallbackItf,
	transactionRepoReads internal.TransactionsReadsRepositoryItf,
	providerRepoReads internal.ProviderReadsRepositoryItf,
) *Merchant {
	return &Merchant{
		merchantRepoReads:     merchantRepoReads,
		merchantRepoWrites:    merchantRepoWrites,
		userRepoReads:         userRepoReads,
		merchantCallbackAdptr: adapterMerchantCallback,
		transactionRepoReads:  transactionRepoReads,
		providerRepoReads:     providerRepoReads,
	}
}

func (mr *Merchant) GetCallbackAttemptsByPaymentId(paymentId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listCallback, err := mr.merchantRepoReads.GetListMerchantCallback(paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listCallback) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Data not found",
			Data:            listCallback,
		}
		return resp, nil
	}

	modifiedListCallback := make([]entity.MerchantCallback, len(listCallback))
	copy(modifiedListCallback, listCallback)

	for i := range modifiedListCallback {
		modifiedListCallback[i].CallbackAt = modifiedListCallback[i].CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00")
		modifiedListCallback[i].CallbackRequestResp = modifiedListCallback[i].CallbackRequest
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            modifiedListCallback,
	}

	return resp, nil
}

func (mr *Merchant) GetListMerchantCallback(params dto.QueryParamsMerchantCallback) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listCallback, pagination, err := mr.merchantRepoReads.GetListMerchantCallbackWithFilter(params)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listCallback) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Data not found",
			Data:            listCallback,
		}
		return resp, nil
	}

	for i := range listCallback {
		merchantId := listCallback[i].MerchantId
		merchantData, err := mr.merchantRepoReads.GetMerchantDataByMerchantId(merchantId)
		if err != nil {
			slog.Infof("get list merchant callback for merchant data failed: %v", err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "failed get merchant data",
			}
			return resp, err
		}

		merchantDetailDataRes := dto.MerchantDataDtoRes{
			Id:           merchantData.Id,
			MerchantId:   merchantData.MerchantId,
			MerchantName: merchantData.MerchantName,
		}
		listCallback[i].MerchantDetailData = merchantDetailDataRes
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve callback data",
		Data:            listCallback,
		Pagination:      pagination,
	}

	return resp, nil
}

func (mr *Merchant) GetListCallbackMerchantSvc(params dto.QueryParamsMerchantCallback) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	user, err := mr.userRepoReads.GetUserByUsername(params.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	params.MerchantName = fmt.Sprintf("[%v]", *user.MerchantName)
	listCallback, pagination, err := mr.merchantRepoReads.GetListMerchantCallbackWithFilter(params)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listCallback) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Data not found",
			Data:            listCallback,
			Pagination:      pagination,
		}
		return resp, nil
	}

	for i := range listCallback {
		merchantId := listCallback[i].MerchantId
		merchantData, err := mr.merchantRepoReads.GetMerchantDataByMerchantId(merchantId)
		if err != nil {
			slog.Infof("get list merchant callback for merchant data failed: %v", err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "failed get merchant data",
			}
			return resp, err
		}

		merchantDetailDataRes := dto.MerchantDataDtoRes{
			Id:           merchantData.Id,
			MerchantId:   merchantData.MerchantId,
			MerchantName: merchantData.MerchantName,
		}
		listCallback[i].MerchantDetailData = merchantDetailDataRes
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve callback data",
		Data:            listCallback,
		Pagination:      pagination,
	}

	return resp, nil
}

func (mr *Merchant) GetLatestMerchantCallback(paymentId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listCallback, err := mr.merchantRepoReads.GetListMerchantCallback(paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listCallback) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Data not found",
			Data:            listCallback,
		}
		return resp, nil
	}

	latestCallback := listCallback[0]
	latestCallback.RetriedAt = listCallback[0].CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00")
	latestCallback.StartedAt = listCallback[len(listCallback)-1].CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00")
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success get latest callback",
		Data:            latestCallback,
	}

	return resp, nil
}

func (mr *Merchant) GetListManualPaymentWithFilter(params dto.QueryParamsManualPayment) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listManualPayment, pagination, err := mr.merchantRepoReads.GetListManualPayment(params)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listManualPayment) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Data not found",
			Data:            listManualPayment,
			Pagination:      pagination,
		}
		return resp, nil
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve manual payment data",
		Data:            listManualPayment,
		Pagination:      pagination,
	}

	return resp, nil
}

func (mr *Merchant) TopUpMerchantSvc(payload dto.AdjustBalanceReqPayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// user data
	user, err := mr.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(payload.Pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	merchantAccountBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.MerchantId)
	if err != nil {
		slog.Infof("top-up mechant id %v got failed: %v", payload.MerchantId, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "failed retrive data maybe wrong merchant id",
		}
		return resp, err
	}

	if merchantAccountBalance.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong merchant id",
		}
		return resp, errors.New("wrong merchant id")
	}

	topUpBalance := merchantAccountBalance.SettledBalance + float64(payload.Amount)
	balanceCapitalFlow := merchantAccountBalance.BalanceCapitalFlow + float64(payload.Amount)
	formattedTopUpBalance := helper.FormatFloat64(topUpBalance)
	formattedBalanceCapitalFlow := helper.FormatFloat64(balanceCapitalFlow)

	// update merchant account balance
	err = mr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedTopUpBalance, formattedBalanceCapitalFlow, payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// create merchant capital flow
	formattedPayloadAmount := helper.FormatFloat64(float64(payload.Amount))
	randomStr := helper.GenerateRandomString(30)
	id := "top_up-" + randomStr
	payloadMerchantCapitalFlowTopUp := dto.CreateMerchantCapitalFlowPayload{
		PaymentId:         id,
		MerchantAccountId: merchantAccountBalance.Id,
		TempBalance:       formattedBalanceCapitalFlow,
		ReasonId:          constant.ReasonIdTopUp,
		Status:            constant.StatusSuccess,
		CreateBy:          payload.Username,
		Amount:            formattedPayloadAmount,
		Notes:             payload.Notes,
		CapitalType:       constant.CapitalTypeCredit,
	}
	_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowTopUp)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	merchantAccountAfterTopUp, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	msg := fmt.Sprintf("success top up for merchant id: %v", merchantAccountBalance.MerchantId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
		Data:            merchantAccountAfterTopUp,
	}

	return resp, nil
}

func (mr *Merchant) HoldBalanceSvc(payload dto.AdjustBalanceReqPayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	balanceSettleOrNotSettleFlagging := constant.SettleBalance

	// user data
	user, err := mr.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(payload.Pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	merchantAccountBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if merchantAccountBalance.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong merchant id",
		}
		return resp, errors.New("wrong merchant id")
	}

	settleOrNotSettleBalance := merchantAccountBalance.SettledBalance
	if settleOrNotSettleBalance < float64(payload.Amount) {
		settleOrNotSettleBalance = merchantAccountBalance.NotSettledBalance
		balanceSettleOrNotSettleFlagging = constant.NotSettledBalance
	}

	adjustedSettleOrNotSettleBalance := settleOrNotSettleBalance - float64(payload.Amount)
	holdBalance := merchantAccountBalance.HoldBalance + float64(payload.Amount)
	formattedHoldBalance := helper.FormatFloat64(holdBalance)
	formattedAdjustMerchantBalance := helper.FormatFloat64(adjustedSettleOrNotSettleBalance)

	// if using settle balance to hold (updated)
	if balanceSettleOrNotSettleFlagging == constant.SettleBalance {
		err = mr.merchantRepoWrites.UpdateMerchantHoldBalanceAndSettleBalance(formattedAdjustMerchantBalance, formattedHoldBalance, payload.MerchantId)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}
	}

	// if using not settle balance to hold (updated)
	if balanceSettleOrNotSettleFlagging == constant.NotSettledBalance {
		err = mr.merchantRepoWrites.UpdateMerchantHoldBalanceAndNotSettleBalance(formattedAdjustMerchantBalance, formattedHoldBalance, payload.MerchantId)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}
	}

	// create merchant capital flow
	formattedPayloadAmount := helper.FormatFloat64(float64(payload.Amount))
	randomStr := helper.GenerateRandomString(30)
	id := "hld_blnc-" + randomStr
	tempBalance := helper.FormatFloat64(merchantAccountBalance.BalanceCapitalFlow)
	payloadMerchantCapitalFlowHoldBalance := dto.CreateMerchantCapitalFlowPayload{
		PaymentId:         id,
		MerchantAccountId: merchantAccountBalance.Id,
		TempBalance:       tempBalance,
		ReasonId:          constant.ReasonIdHoldBalance,
		Status:            constant.StatusSuccess,
		CreateBy:          payload.Username,
		Amount:            formattedPayloadAmount,
		Notes:             payload.Notes,
		CapitalType:       constant.CapitalTypeNotDebitNotCredit,
	}
	_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowHoldBalance)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	merchantAccountAfterHoldBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	msg := fmt.Sprintf("success hold balance for merchant id: %v", merchantAccountBalance.MerchantId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
		Data:            merchantAccountAfterHoldBalance,
	}

	return resp, nil
}

func (mr *Merchant) SettlementBalanceSvc(payload dto.AdjustBalanceReqPayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// user data
	user, err := mr.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(payload.Pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	merchantAccountBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if merchantAccountBalance.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "wrong merchant id",
		}
		return resp, errors.New("wrong merchant id")
	}

	if merchantAccountBalance.NotSettledBalance < float64(payload.Amount) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "not balance enough",
		}

		return resp, errors.New("not balance enough")
	}

	notSettleBalance := merchantAccountBalance.NotSettledBalance
	settleBalance := merchantAccountBalance.SettledBalance

	notSettleBalanceDebited := notSettleBalance - float64(payload.Amount)
	settleBalanceCredited := settleBalance + float64(payload.Amount)
	formattedNoSettleBalance := helper.FormatFloat64(notSettleBalanceDebited)
	formattedSettleBalance := helper.FormatFloat64(settleBalanceCredited)

	err = mr.merchantRepoWrites.UpdateMerchantSettlement(formattedSettleBalance, formattedNoSettleBalance, merchantAccountBalance.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// create merchant capital flow
	formattedPayloadAmount := helper.FormatFloat64(float64(payload.Amount))
	randomStr := helper.GenerateRandomString(30)
	id := "settle_blnc-" + randomStr
	tempBalance := helper.FormatFloat64(merchantAccountBalance.BalanceCapitalFlow)
	payloadMerchantCapitalFlowSettlementBalance := dto.CreateMerchantCapitalFlowPayload{
		PaymentId:         id,
		MerchantAccountId: merchantAccountBalance.Id,
		TempBalance:       tempBalance,
		ReasonId:          constant.ReasonIdSettlement,
		Status:            constant.StatusSuccess,
		CreateBy:          payload.Username,
		Amount:            formattedPayloadAmount,
		Notes:             payload.Notes,
		CapitalType:       constant.CapitalTypeNotDebitNotCredit,
	}
	_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowSettlementBalance)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	merchantAccountAfterSettlement, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	msg := fmt.Sprintf("success settled for merchant id: %v", merchantAccountBalance.MerchantId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
		Data:            merchantAccountAfterSettlement,
	}

	return resp, nil
}

func (mr *Merchant) BalanceTransferSvc(payload dto.BalanceTrfReqPayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// user data
	user, err := mr.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(payload.Pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	// adjust balance merchant account from first
	merchantAccountBalanceFrom, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.AccountFrom.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if merchantAccountBalanceFrom.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong merchant id for account from",
		}
		return resp, errors.New("wrong merchant id")
	}

	settleBalanceFrom := merchantAccountBalanceFrom.SettledBalance
	balanceCapitalFrom := merchantAccountBalanceFrom.BalanceCapitalFlow
	adjustSettleBalanceFrom := settleBalanceFrom - float64(payload.Amount)
	adjustBalanceCapitalFrom := balanceCapitalFrom - float64(payload.Amount)
	formattedSettleBalanceFrom := helper.FormatFloat64(adjustSettleBalanceFrom)
	formattedBalanceCapitalFrom := helper.FormatFloat64(adjustBalanceCapitalFrom)

	// update first merchant balance account from
	err = mr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalanceFrom, formattedBalanceCapitalFrom, merchantAccountBalanceFrom.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// create merchant capital flow for merchant account from
	formattedPayloadAmount := helper.FormatFloat64(float64(payload.Amount))
	randomStr := helper.GenerateRandomString(30)
	id := "blnc_trf-" + randomStr

	payloadMerchantCapitalFlowAccountFrom := dto.CreateMerchantCapitalFlowPayload{
		PaymentId:         id,
		MerchantAccountId: merchantAccountBalanceFrom.Id,
		TempBalance:       formattedBalanceCapitalFrom,
		ReasonId:          constant.ReasonIdBalanceTransfer,
		Status:            constant.StatusSuccess,
		CreateBy:          payload.Username,
		Amount:            formattedPayloadAmount,
		Notes:             payload.Notes,
		CapitalType:       constant.CapitalTypeDebit,
	}
	_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowAccountFrom)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// adjust merchant account to add
	// get merchant account beneficiary
	merchantAccountBalanceTo, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.AccountTo.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if merchantAccountBalanceTo.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong merchant id for beneficiary",
		}
		return resp, errors.New("wrong merchant id")
	}

	settleBalanceTo := merchantAccountBalanceTo.SettledBalance
	balanceCapitalTo := merchantAccountBalanceTo.BalanceCapitalFlow
	adjustSettleBalanceTo := settleBalanceTo + float64(payload.Amount)
	adjustBalanceCapitalTo := balanceCapitalTo + float64(payload.Amount)
	formattedSettleBalanceTo := helper.FormatFloat64(adjustSettleBalanceTo)
	formattedBalanceCapitalTo := helper.FormatFloat64(adjustBalanceCapitalTo)

	// update merchant balance account to
	err = mr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalanceTo, formattedBalanceCapitalTo, merchantAccountBalanceTo.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// create merchant capital flow for merchant account to
	payloadMerchantCapitalFlowAccountTo := dto.CreateMerchantCapitalFlowPayload{
		PaymentId:         id,
		MerchantAccountId: merchantAccountBalanceTo.Id,
		TempBalance:       formattedBalanceCapitalTo,
		ReasonId:          constant.ReasonIdBalanceTransfer,
		Status:            constant.StatusSuccess,
		CreateBy:          payload.Username,
		Amount:            formattedPayloadAmount,
		Notes:             payload.Notes,
		CapitalType:       constant.CapitalTypeCredit,
	}
	_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowAccountTo)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	msg := fmt.Sprintf("success transfer balance from %v to merchant account %v", merchantAccountBalanceFrom.MerchantId, merchantAccountBalanceTo.MerchantId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
	}

	return resp, nil
}

func (mr *Merchant) SendCallbackSvc(paymentId string, username string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	transactionDetail, err := mr.transactionRepoReads.GetPaymentDetailProviderMerchant(paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if transactionDetail.TransactionID == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "data didn't found maybe wrong transaction id",
		}

		return resp, errors.New("wrong transaction id")
	}

	merchantData, err := mr.merchantRepoReads.GetMerchantDataByMerchantId(transactionDetail.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	transactionStatusLatest, err := mr.transactionRepoReads.GetStatusChangeLogData(transactionDetail.PaymentID)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// send callback http
	callbackStatus := constant.StatusSuccess
	merchantResponse, err := mr.merchantCallbackAdptr.SendCallbackAdptr(transactionDetail.MerchantCallbackURL, transactionDetail, transactionStatusLatest[0], merchantData.MerchantSecret)
	if err != nil {
		callbackStatus = constant.StatusFailed
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "merchant didn't receive with status ok",
		}

		if err.Error() == "status not ok" {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusOK,
				ResponseMessage: "merchant didn't receive with status ok",
			}

			_, err := mr.merchantRepoWrites.CreateMerchantCallback(transactionDetail.PaymentID, callbackStatus, transactionDetail.Status, merchantResponse.(string), username)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
			}

			return resp, nil
		}

		_, err := mr.merchantRepoWrites.CreateMerchantCallback(transactionDetail.PaymentID, callbackStatus, transactionDetail.Status, merchantResponse.(string), username)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
		}

		return resp, err
	}

	// merchant callback
	_, err = mr.merchantRepoWrites.CreateMerchantCallback(transactionDetail.PaymentID, callbackStatus, transactionDetail.Status, merchantResponse.(string), username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	msg := fmt.Sprintf("success send callback for transaction id: %v", paymentId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
	}

	return resp, nil
}

func (mr *Merchant) PayoutSettlementSvc(payload dto.AdjustBalanceReqPayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// user data
	user, err := mr.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(payload.Pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	merchantAccountBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "failed retrieve data maybe wrong merchant id",
		}
		return resp, err
	}

	if merchantAccountBalance.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong merchant id",
		}
		return resp, errors.New("wrong merchant id")
	}

	outBalance := merchantAccountBalance.SettledBalance - float64(payload.Amount)
	balanceCapitalFlowMinusPayoutSettlement := merchantAccountBalance.BalanceCapitalFlow - float64(payload.Amount)
	formattedOutBalance := helper.FormatFloat64(outBalance)
	formattedBalanceCapitalFlowMinusPayoutSettlement := helper.FormatFloat64(balanceCapitalFlowMinusPayoutSettlement)

	// update merchant account balance
	err = mr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedOutBalance, formattedBalanceCapitalFlowMinusPayoutSettlement, payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// create merchant capital flow
	formattedPayloadAmount := helper.FormatFloat64(float64(payload.Amount))
	randomStr := helper.GenerateRandomString(30)
	id := "out_sttlmnt-" + randomStr
	payloadMerchantCapitalFlowTopUp := dto.CreateMerchantCapitalFlowPayload{
		PaymentId:         id,
		MerchantAccountId: merchantAccountBalance.Id,
		TempBalance:       formattedBalanceCapitalFlowMinusPayoutSettlement,
		ReasonId:          constant.ReasonIdOutSettlement,
		Status:            constant.StatusSuccess,
		CreateBy:          payload.Username,
		Amount:            formattedPayloadAmount,
		Notes:             payload.Notes,
		CapitalType:       constant.CapitalTypeDebit,
	}
	_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowTopUp)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	merchantAccountAfterOutSettlement, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(payload.MerchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	msg := fmt.Sprintf("success out settlement for merchant id: %v", merchantAccountBalance.MerchantId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
		Data:            merchantAccountAfterOutSettlement,
	}

	return resp, nil
}

func (mr *Merchant) GetDetailManualPaymentSvc(paymentId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var manualPaymentDetailRespDto dto.ManualPaymentDetailDto

	manualPaymentDetail, err := mr.merchantRepoReads.GetDetailManualPayment(paymentId)
	if err != nil {
		slog.Infof("manual payment detail got failed due to: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "data not found maybe wrong payment id",
		}

		return resp, err
	}

	if len(manualPaymentDetail) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "data not found maybe wrong payment id",
		}

		return resp, nil
	}

	var notes string
	if manualPaymentDetail[0].Notes != nil {
		notes = *manualPaymentDetail[0].Notes
	} else {
		notes = "" // or handle the case where notes should have a default value
	}
	if len(manualPaymentDetail) > 1 {
		// handle if get manual payment with reason id payin payout and fee
		if manualPaymentDetail[0].ReasonName == constant.ReasonNameInTransaction ||
			manualPaymentDetail[0].ReasonName == constant.ReasonNameOutTransaction ||
			manualPaymentDetail[0].ReasonName == constant.ReasonNameFeeTransaction {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "this is payment id for transaction",
			}

			return resp, errors.New("payment id for transaction")
		}

		// handle for balance transfer
		if manualPaymentDetail[0].CapitalType == constant.CapitalTypeCredit {
			manualPaymentDetailRespDto = dto.ManualPaymentDetailDto{
				CreditedAccount: manualPaymentDetail[0].MerchantId,
				DebitedAccount:  manualPaymentDetail[1].MerchantId,
				Reason:          manualPaymentDetail[0].ReasonName,
				Notes:           notes,
			}
		} else {
			manualPaymentDetailRespDto = dto.ManualPaymentDetailDto{
				CreditedAccount: manualPaymentDetail[1].MerchantId,
				DebitedAccount:  manualPaymentDetail[0].MerchantId,
				Reason:          manualPaymentDetail[0].ReasonName,
				Notes:           notes,
			}
		}
	}

	if len(manualPaymentDetail) == 1 {
		if manualPaymentDetail[0].CapitalType == constant.CapitalTypeCredit {
			manualPaymentDetailRespDto = dto.ManualPaymentDetailDto{
				CreditedAccount: manualPaymentDetail[0].MerchantId,
				Reason:          manualPaymentDetail[0].ReasonName,
				Notes:           notes,
			}
		} else if manualPaymentDetail[0].CapitalType == constant.CapitalTypeNotDebitNotCredit && manualPaymentDetail[0].ReasonName == constant.ReasonNameSettlement {
			manualPaymentDetailRespDto = dto.ManualPaymentDetailDto{
				CreditedAccount: manualPaymentDetail[0].MerchantId,
				Reason:          manualPaymentDetail[0].ReasonName,
				Notes:           notes,
			}
		} else {
			manualPaymentDetailRespDto = dto.ManualPaymentDetailDto{
				DebitedAccount: manualPaymentDetail[0].MerchantId,
				Reason:         manualPaymentDetail[0].ReasonName,
				Notes:          notes,
			}
		}
	}

	msg := fmt.Sprintf("succces get detail manual payment id: %v", paymentId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
		Data:            manualPaymentDetailRespDto,
	}

	return resp, nil
}

func (mr *Merchant) GetListOtherTransactionsSvc(params dto.QueryParamsManualPayment) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if params.MinDate == "" {
		params.MinDate = helper.GenerateTime(0)
	}

	if params.MaxDate == "" {
		params.MaxDate = helper.GenerateTime(24)
	}

	user, err := mr.userRepoReads.GetUserByUsername(params.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	params.MerchantName = fmt.Sprintf("[%v]", *user.MerchantName)
	listOtherTransactions, pagination, err := mr.merchantRepoReads.GetListManualPayment(params)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listOtherTransactions) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Data not found",
			Data:            listOtherTransactions,
			Pagination:      pagination,
		}
		return resp, nil
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve manual payment data",
		Data:            listOtherTransactions,
		Pagination:      pagination,
	}

	return resp, nil
}

func (mr *Merchant) GetListMerchantWithFilterSvc(params dto.QueryParams) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listMerchants, err := mr.merchantRepoReads.GetListMerchantWithFilterRepo(params)
	if err != nil {
		slog.Infof("get merchant list got err: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "data not found maybe wrong params",
		}

		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success get list merchant",
		Data:            listMerchants,
	}

	return resp, nil
}

func (mr *Merchant) ReverseManualPaymentSvc(payload dto.UpdateStatusTransaction, username string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// user data
	user, err := mr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(payload.Pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	manualPaymentData, err := mr.merchantRepoReads.GetDetailManualPayment(payload.PaymentId)
	if err != nil {
		slog.Infof("get manual payment data error: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "oops there is something error please ask customer support for more detail",
		}

		return resp, err
	}

	if len(manualPaymentData) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "data not found maybe wrong transaction id",
		}

		return resp, nil
	}

	// validate transaction ID
	if manualPaymentData[0].ReasonName == constant.ReasonNameInTransaction || manualPaymentData[0].ReasonName == constant.ReasonNameOutTransaction || manualPaymentData[0].ReasonName == constant.ReasonNameFeeTransaction {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "this is transaction data not manual payment",
		}

		return resp, errors.New("this is transaction data not manual payment")
	}

	// validate reverse status
	if manualPaymentData[0].Status == constant.StatusReversed {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "can't reverse twice",
		}
		return resp, errors.New("double reverse")
	}

	manualPaymentDataReverse, err := mr.merchantRepoReads.CheckReverseStatusRepo(payload.PaymentId)
	if err != nil {
		slog.Infof("%v check reverse got failed: %v", payload.PaymentId, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "can't reverse twice",
		}
		return resp, errors.New("double reverse")
	}

	if len(manualPaymentDataReverse) > 0 {
		slog.Infof("%v check reverse got failed: double reverse", payload.PaymentId)
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "can't reverse twice",
		}
		return resp, errors.New("double reverse")
	}

	// adjust for balance trf
	if len(manualPaymentData) > 1 {
		var balanceTrfCreditor entity.ManualPayment
		var balanceTrfDebitor entity.ManualPayment

		for _, manualPayment := range manualPaymentData {
			if manualPayment.CapitalType == constant.CapitalTypeCredit {
				balanceTrfCreditor = manualPayment
			}

			if manualPayment.CapitalType == constant.CapitalTypeDebit {
				balanceTrfDebitor = manualPayment
			}
		}

		creditorMerchantAccount, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(balanceTrfCreditor.MerchantId)
		if err != nil {
			slog.Infof("creditor merchant account got error: %v", err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "oops there is something error please ask customer support for more detail",
			}
			return resp, err
		}

		// adjust balance creditor
		creditorBalanceCapitalMinusBalanceTrf := creditorMerchantAccount.BalanceCapitalFlow - balanceTrfCreditor.Amount
		creditorSettleBalanceMinusBalanceTrf := creditorMerchantAccount.SettledBalance - balanceTrfCreditor.Amount
		formattedBalanceCapitalCreditor := helper.FormatFloat64(creditorBalanceCapitalMinusBalanceTrf)
		formattedSettleBalanceCreditor := helper.FormatFloat64(creditorSettleBalanceMinusBalanceTrf)
		formattedCreditorAmount := helper.FormatFloat64(balanceTrfCreditor.Amount)

		err = mr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalanceCreditor, formattedBalanceCapitalCreditor, balanceTrfCreditor.MerchantId)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}

		// create merchant capital flow for adjust creditor
		reverseNotes := payload.Notes + fmt.Sprintf(" (reverse transaction id %v)", balanceTrfCreditor.PaymentId)
		randomStr := helper.GenerateRandomString(30)
		id := "rvrs_balance-" + randomStr
		payloadMerchantCapitalCreditor := dto.CreateMerchantCapitalFlowPayload{
			PaymentId:         id,
			MerchantAccountId: creditorMerchantAccount.Id,
			TempBalance:       formattedBalanceCapitalCreditor,
			ReasonId:          balanceTrfCreditor.ReasonId,
			Status:            constant.StatusReversed,
			CreateBy:          username,
			Amount:            formattedCreditorAmount,
			Notes:             reverseNotes,
			CapitalType:       constant.CapitalTypeDebit,
			ReverseFrom:       balanceTrfCreditor.PaymentId,
		}
		_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalCreditor)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}

		debitorMerchantAccount, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(balanceTrfDebitor.MerchantId)
		if err != nil {
			slog.Infof("creditor merchant account got error: %v", err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "oops there is something error please ask customer support for more detail",
			}
			return resp, err
		}

		// adjust balance debitor
		debitorBalanceCapitalAddBalanceTrf := debitorMerchantAccount.BalanceCapitalFlow + balanceTrfDebitor.Amount
		debitorSettleBalanceAddBalanceTrf := debitorMerchantAccount.SettledBalance + balanceTrfDebitor.Amount
		formattedDebitorAmount := helper.FormatFloat64(balanceTrfDebitor.Amount)
		formattedBalanceCapitalDebitor := helper.FormatFloat64(debitorBalanceCapitalAddBalanceTrf)
		formatedSettleBalanceDebitor := helper.FormatFloat64(debitorSettleBalanceAddBalanceTrf)

		// update merchant balance debitor
		err = mr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formatedSettleBalanceDebitor, formattedBalanceCapitalDebitor, balanceTrfDebitor.MerchantId)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}

		// create merchant capital flow for adjust debitor
		payloadMerchantCapitalDebitor := dto.CreateMerchantCapitalFlowPayload{
			PaymentId:         id,
			MerchantAccountId: debitorMerchantAccount.Id,
			TempBalance:       formattedBalanceCapitalDebitor,
			ReasonId:          balanceTrfDebitor.ReasonId,
			Status:            constant.StatusReversed,
			CreateBy:          username,
			Amount:            formattedDebitorAmount,
			Notes:             reverseNotes,
			CapitalType:       constant.CapitalTypeCredit,
			ReverseFrom:       balanceTrfCreditor.PaymentId,
		}
		_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalDebitor)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}

		msg := fmt.Sprintf("success reversed with id %v for balance trf id %v", id, balanceTrfCreditor.PaymentId)
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: msg,
		}
	}

	if len(manualPaymentData) == 1 {
		if manualPaymentData[0].ReasonId == constant.ReasonIdTopUp {
			merchantAccountBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(manualPaymentData[0].MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// reverse balance
			settleBalanceMinusTopup := merchantAccountBalance.SettledBalance - manualPaymentData[0].Amount
			balanceCapitalMinusTopup := merchantAccountBalance.BalanceCapitalFlow - manualPaymentData[0].Amount
			formattedTopupAmount := helper.FormatFloat64(manualPaymentData[0].Amount)
			formattedSettleBalanceMinusTopup := helper.FormatFloat64(settleBalanceMinusTopup)
			formattedBalanceCapitalMinusTopup := helper.FormatFloat64(balanceCapitalMinusTopup)

			err = mr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalanceMinusTopup, formattedBalanceCapitalMinusTopup, merchantAccountBalance.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// create merchant capital flow for reverse top up
			reverseNotes := payload.Notes + fmt.Sprintf(" (reverse transaction id %v)", manualPaymentData[0].PaymentId)
			randomStr := helper.GenerateRandomString(30)
			id := "rvrs_balance-" + randomStr
			payloadMerchantCapitalTopup := dto.CreateMerchantCapitalFlowPayload{
				PaymentId:         id,
				MerchantAccountId: merchantAccountBalance.Id,
				TempBalance:       formattedBalanceCapitalMinusTopup,
				ReasonId:          manualPaymentData[0].ReasonId,
				Status:            constant.StatusReversed,
				CreateBy:          username,
				Amount:            formattedTopupAmount,
				Notes:             reverseNotes,
				CapitalType:       constant.CapitalTypeDebit,
				ReverseFrom:       manualPaymentData[0].PaymentId,
			}
			_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalTopup)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			msg := fmt.Sprintf("success reverse with id %v, for top up id %v", id, manualPaymentData[0].PaymentId)
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusOK,
				ResponseMessage: msg,
			}
		}

		if manualPaymentData[0].ReasonId == constant.ReasonIdHoldBalance {
			merchantAccountBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(manualPaymentData[0].MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// reverse balance for hold balance
			formattedBalanceCapitalHoldBalance := helper.FormatFloat64(merchantAccountBalance.BalanceCapitalFlow)
			holdBalanceMinusAmount := merchantAccountBalance.HoldBalance - manualPaymentData[0].Amount
			settleBalanceAddAmount := merchantAccountBalance.SettledBalance + manualPaymentData[0].Amount
			formattedHoldBalanceMinusAmount := helper.FormatFloat64(holdBalanceMinusAmount)
			formattedSettleBalanceAddAmount := helper.FormatFloat64(settleBalanceAddAmount)
			formattedAmountHoldBalance := helper.FormatFloat64(manualPaymentData[0].Amount)

			err = mr.merchantRepoWrites.UpdateMerchantHoldBalanceAndSettleBalance(formattedSettleBalanceAddAmount, formattedHoldBalanceMinusAmount, merchantAccountBalance.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// create merchant capital flow for reverse hold balance
			reverseNotes := payload.Notes + fmt.Sprintf(" (reverse transaction id %v)", manualPaymentData[0].PaymentId)
			randomStr := helper.GenerateRandomString(30)
			id := "rvrs_balance-" + randomStr
			payloadMerchantCapitalHoldBalance := dto.CreateMerchantCapitalFlowPayload{
				PaymentId:         id,
				MerchantAccountId: merchantAccountBalance.Id,
				TempBalance:       formattedBalanceCapitalHoldBalance,
				ReasonId:          manualPaymentData[0].ReasonId,
				Status:            constant.StatusReversed,
				CreateBy:          username,
				Amount:            formattedAmountHoldBalance,
				Notes:             reverseNotes,
				CapitalType:       constant.CapitalTypeNotDebitNotCredit,
				ReverseFrom:       manualPaymentData[0].PaymentId,
			}
			_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalHoldBalance)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			msg := fmt.Sprintf("success reverse with id %v, for hold balance id %v", id, manualPaymentData[0].PaymentId)
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusOK,
				ResponseMessage: msg,
			}
		}

		if manualPaymentData[0].ReasonId == constant.ReasonIdSettlement {
			merchantAccountBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(manualPaymentData[0].MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// reverse balance
			formattedBalanceCapitalMerchant := helper.FormatFloat64(merchantAccountBalance.BalanceCapitalFlow)
			settleBalanceMinusAmount := merchantAccountBalance.SettledBalance - manualPaymentData[0].Amount
			notSettleBalanceAddAmount := merchantAccountBalance.NotSettledBalance + manualPaymentData[0].Amount
			formattedSettleBalanceMinusAmount := helper.FormatFloat64(settleBalanceMinusAmount)
			formattedNotSettleBalanceAddAmount := helper.FormatFloat64(notSettleBalanceAddAmount)
			formattedSettlementAmount := helper.FormatFloat64(manualPaymentData[0].Amount)

			// update settle balance
			err = mr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalanceMinusAmount, formattedBalanceCapitalMerchant, manualPaymentData[0].MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// update not settle balance
			err = mr.merchantRepoWrites.UpdateMerchantCapitalAndNotSettleBalance(formattedNotSettleBalanceAddAmount, formattedBalanceCapitalMerchant, manualPaymentData[0].MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			reverseNotes := payload.Notes + fmt.Sprintf(" (reverse transaction id %v)", manualPaymentData[0].PaymentId)
			randomStr := helper.GenerateRandomString(30)
			id := "rvrs_balance-" + randomStr
			payloadMerchantCapitalSettlement := dto.CreateMerchantCapitalFlowPayload{
				PaymentId:         id,
				MerchantAccountId: merchantAccountBalance.Id,
				TempBalance:       formattedBalanceCapitalMerchant,
				ReasonId:          manualPaymentData[0].ReasonId,
				Status:            constant.StatusReversed,
				CreateBy:          username,
				Amount:            formattedSettlementAmount,
				Notes:             reverseNotes,
				CapitalType:       constant.CapitalTypeNotDebitNotCredit,
				ReverseFrom:       manualPaymentData[0].PaymentId,
			}
			_, err = mr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalSettlement)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			msg := fmt.Sprintf("success reverse with id %v, for settlement id %v", id, manualPaymentData[0].PaymentId)
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusOK,
				ResponseMessage: msg,
			}
		}

		if manualPaymentData[0].ReasonId == constant.ReasonIdOutSettlement {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "can't reverse for out settlement",
			}

			return resp, errors.New("can't reverse out settlement")
		}
	}

	return resp, nil
}

func (mr *Merchant) GetMerchantBalanceSvc(merchantId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	merchantBalanceData, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(merchantId)
	if err != nil {
		slog.Infof("got error GetMerchantAccountByMerchantId: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "oopss got error, please ask customer support for more detail",
		}

		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrive merchant balance data",
		Data:            merchantBalanceData,
	}

	return resp, nil
}

func (mr *Merchant) CreateMerchantSvc(payload dto.CreateMerchantDtoReq) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	randomStrMerchantId := helper.GenerateRandomString(10)
	randomStrMerchantScret := helper.GenerateRandomString(30)
	merchantId := "hypy_mrchnt-" + randomStrMerchantId
	merchantSecret := "secret_key-" + randomStrMerchantScret

	createMerchantId, err := mr.merchantRepoWrites.CreateMerchantRepo(payload.MerchantName, merchantId, merchantSecret)
	if err != nil {
		msg := fmt.Sprintf("failed to create merchant with err: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: msg,
		}
		return resp, err
	}

	_, err = mr.merchantRepoWrites.CreateMerchantAccountsRepo(merchantId)
	if err != nil {
		msg := fmt.Sprintf("failed to create merchant account with err: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: msg,
		}
		return resp, err
	}

	if len(payload.PaymentMethod) > 0 {
		// create payment method and paychannel
		for _, method := range payload.PaymentMethod {
			// create payment method
			paymentMethodId, err := mr.merchantRepoWrites.CreateMerchantPaymentMethodRepo(createMerchantId, method.PaymentMethodId)
			if err != nil {
				msg := fmt.Sprintf("failed to create payment method with err: %v", err.Error())
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: msg,
				}
				return resp, err
			}

			// create merchant paychannel
			segment := "MAIN"
			methodStr := constant.TransformPaymentMethodNameIntoCode[method.PaymentMethodName]
			randomStrMerchantChannelCode := helper.GenerateRandomString(15)
			merchantPaychannelCode := methodStr + "-" + strings.ToUpper(segment) + "-" + randomStrMerchantChannelCode
			minAmount := method.MinAmountPerTransaction
			maxAmount := method.MaxAmountPerTransaction
			dailyLimit := method.DailyLimit
			fee := method.Fee
			feeType := method.FeeType
			if feeType == "" {
				feeType = constant.FeeTypeFixedFee
			}

			_, err = mr.merchantRepoWrites.CreateMerchantPaychannelRepo(paymentMethodId, segment, fee, feeType, minAmount, maxAmount, dailyLimit, merchantPaychannelCode)
			if err != nil {
				msg := fmt.Sprintf("failed to create merchant paychannel with err: %v", err.Error())
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: msg,
				}
				return resp, err
			}
		}

		// get merchant paychannel
		merchantPaychannelList, err := mr.merchantRepoReads.GetMerchantPaychannelByMerchantId(merchantId)
		if err != nil {
			slog.Infof("create merchant name %v got failed get merchant paychannel list: %v", payload.MerchantName, err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "failed get channel list but merchant already created",
			}
			return resp, err
		}

		if len(merchantPaychannelList) == 0 {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "failed get channel list but merchant already created",
			}
			return resp, errors.New("failed create merchant paychannel")
		}

		active := "0"

		for i := range merchantPaychannelList {
			availableChannel, err := mr.providerRepoReads.CountProviderChannelByPaymentMethodRepo(merchantPaychannelList[i].PaymentMethodChannel)
			if err != nil {
				slog.Infof("create merchant name %v got failed get merchant paychannel list: %v", payload.MerchantName, err.Error())
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: "failed get channel list but merchant already created",
				}
				return resp, err
			}
			availableChannelStr := converter.ToString(availableChannel)
			activeAvailableChannel := active + "/" + availableChannelStr
			merchantPaychannelList[i].ActiveAvailableChannel = activeAvailableChannel
		}

		msg := fmt.Sprintf("success create merchant with id: %v", createMerchantId)
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: msg,
			Data:            merchantPaychannelList,
		}

		return resp, nil
	}

	msg := fmt.Sprintf("success create merchant with id: %v", createMerchantId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
		Data:            "no merchant paychannel created",
	}

	return resp, nil
}

func (mr *Merchant) GetMerchantAnalyticsSvc(payload dto.GetMerchantAnalyticsDtoReq) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if payload.MinDate == "" {
		payload.MinDate = helper.GenerateTime(0)
	}

	if payload.MaxDate == "" {
		payload.MaxDate = helper.GenerateTime(24)
	}

	transactionData, err := mr.transactionRepoReads.GetTransactionAnalyticsRepo(payload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	data := supportMerchantAnalyticsSvc(transactionData)

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve",
		Data:            data,
	}

	return resp, nil
}

func (mr *Merchant) MerchantUpdateStatusSvc(payload dto.AccountData) (dto.ResponseDto, error) {
	var res dto.ResponseDto

	merchantData, err := mr.merchantRepoReads.GetMerchantDataByMerchantId(payload.MerchantId)
	if err != nil {
		msg := fmt.Sprintf("failed get merchant data for merchantId: %v", payload.MerchantId)
		slog.Infof("failed get merchant data %v with err: %v", payload.MerchantId, err.Error())
		res = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: msg,
		}
		return res, err
	}

	// check if merchant already have merchant paychannel
	merchantPaychannel, err := mr.merchantRepoReads.GetMerchantPaychannelByMerchantId(payload.MerchantId)
	if err != nil {
		msg := fmt.Sprintf("failed get merchant paychannel for merchantId: %v", payload.MerchantId)
		slog.Infof("failed get merchant paychannel %v with err: %v", payload.MerchantId, err.Error())
		res = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: msg,
		}
		return res, err
	}

	if len(merchantPaychannel) == 0 {
		msg := fmt.Sprintf("need to create paychannel for merchantId: %v", payload.MerchantId)
		res = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: msg,
		}
		return res, errors.New("failed to update status")
	}

	merchantStatus := merchantData.Status

	var msgRes string
	if merchantStatus == constant.StatusActive {
		// updated to inactive
		msgRes = fmt.Sprintf("Merchant %v has been successfully deactived", payload.MerchantId)
		err := mr.merchantRepoWrites.UpdateMerchantStatusRepo(payload.MerchantId, constant.StatusInactive)
		if err != nil {
			res = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return res, err
		}
	} else {
		// update to active
		msgRes = fmt.Sprintf("Merchant %v has been successfully actived", payload.MerchantId)
		err := mr.merchantRepoWrites.UpdateMerchantStatusRepo(payload.MerchantId, constant.StatusActive)
		if err != nil {
			res = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return res, err
		}
	}

	res = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msgRes,
	}

	return res, nil
}

func (mr *Merchant) GetMerchantPaychannelSvc(merchantId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listMerchantPaychannel, err := mr.merchantRepoReads.GetMerchantPaychannelByMerchantId(merchantId)
	if err != nil {
		slog.Infof("merchant id %v got failed get list paychannel with err: %v", merchantId, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "data not found",
		}
		return resp, err
	}

	if len(listMerchantPaychannel) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Merchant doesn't have merchant paychannel yet",
			Data:            "",
		}
		return resp, nil
	}

	for i := range listMerchantPaychannel {
		availableChannel, err := mr.providerRepoReads.CountProviderChannelByPaymentMethodRepo(listMerchantPaychannel[i].PaymentMethodChannel)
		if err != nil {
			slog.Infof("count available channel %v got failed get merchant paychannel list: %v", merchantId, err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "failed get available channel",
			}
			return resp, err
		}
		availableChannelStr := converter.ToString(availableChannel)

		activeChannel, err := mr.providerRepoReads.CountActiveProviderChannelRepo(listMerchantPaychannel[i].Id)
		if err != nil {
			slog.Infof("count available channel %v got failed get merchant paychannel list: %v", merchantId, err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: "failed get available channel",
			}
			return resp, err
		}

		activeChannelStr := converter.ToString(activeChannel)
		activeAvailableChannel := activeChannelStr + "/" + availableChannelStr
		listMerchantPaychannel[i].ActiveAvailableChannel = activeAvailableChannel
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data list merchant paychannel",
		Data:            listMerchantPaychannel,
	}

	return resp, nil
}

func (mr *Merchant) GetListCapitalFlowTransactionSvc(params dto.QueryParams) (dto.ResponseDto, error) {
	var res dto.ResponseDto

	if params.MinDate == "" {
		params.MinDate = helper.GenerateTime(0)
	}

	if params.MaxDate == "" {
		params.MaxDate = helper.GenerateTime(24)
	}

	listTransactionCapital, pagination, err := mr.transactionRepoReads.GetTransactionCapitalFlowRepo(params)
	if err != nil {
		slog.Infof("merchant id %v list transaction capital got failed: %v", params.MerchantId, err.Error())
		res = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "Data not found",
		}
		return res, err
	}

	if len(listTransactionCapital) == 0 {
		res = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Data not found",
		}
		return res, nil
	}

	res = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve list transaction capital flow",
		Data:            listTransactionCapital,
		Pagination:      pagination,
	}

	return res, nil
}

func (mr *Merchant) GetListTierPaychannelSvc(paychannelId int) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	paychannelDetail, err := mr.merchantRepoReads.GetMerchantPaychannelDetailById(paychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	listPaychannel, err := mr.merchantRepoReads.GetMerchantPaychannelByPaymentMethodId(paychannelDetail.MerchantPaymentMethodId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	listTier := helper.GenerateTiers(listPaychannel)

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success get tier list",
		Data:            listTier,
	}

	return resp, nil
}

func (mr *Merchant) GetListMerchantAccountSvc(params dto.QueryParams) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listAccount, err := mr.merchantRepoReads.GetListMerchantAccountRepo(params)
	if err != nil {
		slog.Infof("get list merchant account got failed: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "Data not found",
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            listAccount,
	}

	return resp, nil
}

func (mr *Merchant) GetMerchantPaychannelAnalyticsSvc(payload dto.GetMerchantAnalyticsDtoReq) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if payload.MinDate == "" {
		payload.MinDate = helper.GenerateTime(0)
	}

	if payload.MaxDate == "" {
		payload.MaxDate = helper.GenerateTime(24)
	}

	transactionData, err := mr.transactionRepoReads.GetTransactionListByMerchantPaychannelRepo(payload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	dataAnalytics := supportMerchantAnalyticsByMerchantPaychannelSvc(transactionData)

	detailMerchantPaychannel, err := mr.merchantRepoReads.GetMerchantPaychannelDetailById(payload.MerchantPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	data := dto.MerchantPaychannelAnalyticsRspDto{
		AnalyticsData:            dataAnalytics,
		MerchantPaychannelDetail: detailMerchantPaychannel,
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            data,
	}

	return resp, nil
}

func (mr *Merchant) GetRoutedPaychannelSvc(merchantPaychannelId int) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	routedPaychannelList, err := mr.merchantRepoReads.GetListRoutedPaychannelByIdMerchantPaychannelRepo(merchantPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, nil
	}

	// convert payload for response
	for i := range routedPaychannelList {
		if routedPaychannelList[i].FeeType == constant.FeeTypePercentage {
			feeDb := converter.ToString(routedPaychannelList[i].FeesDb)
			feeRes := feeDb + "%"
			routedPaychannelList[i].FeeResp = feeRes
		} else {
			feeDb := converter.ToString(routedPaychannelList[i].FeesDb)
			routedPaychannelList[i].FeeResp = feeDb
		}
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            routedPaychannelList,
	}

	return resp, nil
}

func (mr *Merchant) GetPaymentOperatorsSvc(routedPaychannelId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	bankList, err := mr.merchantRepoReads.GetBankListProviderPaymentMethodRepo(routedPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, nil
	}

	bankNameListPaychannel, err := mr.merchantRepoReads.GetBankListFromProviderPaychannelRepo(routedPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, nil
	}

	for j := range bankList {
		for i := range bankNameListPaychannel {
			if bankNameListPaychannel[i] == bankList[j].BankName {
				bankList[j].CheckListFlagging = true
			}
		}
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            bankList,
	}

	return resp, nil
}

func (mr *Merchant) UpdateLimitOrFeeMerchantPaychannelSvc(payload dto.AdjustLimitOrFeePayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// Fetch existing merchant pay channel details
	merchantPaychannelData, err := mr.merchantRepoReads.GetMerchantPaychannelDetailById(payload.MerchantPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: fmt.Sprintf("Failed to retrieve merchant pay channel details: %v", err),
		}
		return resp, err
	}

	// Use existing values if the new values are not provided
	if payload.MaxAmount == nil {
		payload.MaxAmount = &merchantPaychannelData.MaxTransaction
	}

	if payload.MinAmount == nil {
		payload.MinAmount = &merchantPaychannelData.MinTransaction
	}

	if payload.MaxDailyLimit == nil {
		payload.MaxDailyLimit = &merchantPaychannelData.MaxDailyTransaction
	}

	if payload.Fee == nil {
		payload.Fee = &merchantPaychannelData.Fee
	}

	if payload.FeeType == nil {
		payload.FeeType = &merchantPaychannelData.FeeType
	}

	// Update merchant pay channel
	err = mr.merchantRepoWrites.UpdateMerchantPaychannelByIdRepo(payload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: fmt.Sprintf("Failed to update merchant pay channel: %v", err),
		}
		return resp, err
	}

	// Fetch updated merchant pay channel details
	merchantPaychannelDataUpdated, err := mr.merchantRepoReads.GetMerchantPaychannelDetailById(payload.MerchantPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: fmt.Sprintf("Failed to retrieve updated merchant pay channel details: %v", err),
		}
		return resp, err
	}

	msg := fmt.Sprintf("Successfully updated merchant pay channel with id: %v", payload.MerchantPaychannelId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
		Data:            merchantPaychannelDataUpdated,
	}

	return resp, nil
}

func (mr *Merchant) GetAggregatedPaychannelSvc(paychannelId int) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	aggregatedData, err := mr.merchantRepoReads.GetAggregatedPaychannelByIdRepo(paychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            aggregatedData,
	}

	return resp, nil
}

func (mr *Merchant) AddSegmentMerchantPaychannelSvc(payload dto.AddSegmentDtoReq) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	paychannel, err := mr.merchantRepoReads.GetMerchantPaychannelDetailById(payload.MerchantPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, nil
	}

	// create merchant paychannel
	method := constant.TransformPaymentMethodNameIntoCode[paychannel.PaymentMethodChannel]
	randomStrMerchantChannelCode := helper.GenerateRandomString(15)
	merchantPaychannelCode := method + "-" + strings.ToUpper(payload.TierName) + "-" + randomStrMerchantChannelCode
	minAmount := 0.00
	maxAmount := 0.00
	dailyLimit := 0.00
	fee := 0.00
	feeType := constant.FeeTypeFixedFee

	_, err = mr.merchantRepoWrites.CreateMerchantPaychannelRepo(paychannel.MerchantPaymentMethodId, payload.TierName, fee, feeType, minAmount, maxAmount, dailyLimit, merchantPaychannelCode)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	msg := fmt.Sprintf("Success created merchant paychannel with tier: %v", payload.TierName)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
	}

	return resp, nil
}

func (mr *Merchant) ActivateOrDeactivateMerchantPaychannelSvc(payload dto.AccountData) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	status := constant.StatusActive

	merchanPaychannelData, err := mr.merchantRepoReads.GetMerchantPaychannelDetailById(payload.MerchantPaychannelId)
	if err != nil {
		slog.Infof("get merchanPaychannelData with ID %v failed: %v", payload.MerchantPaychannelId, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "Data not found",
		}
		return resp, err
	}

	if merchanPaychannelData.Status == constant.StatusActive {
		status = constant.StatusInactive
	}

	// checking routed paychanne
	routedList, err := mr.merchantRepoReads.GetListRoutedPaychannelByIdMerchantPaychannelRepo(merchanPaychannelData.Id)
	if err != nil {
		slog.Infof("routedList got error: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "Can't activate due to didn't have routed paychannel",
		}
		return resp, err
	}

	if len(routedList) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "Can't activate due to didn't have routed paychannel",
		}
		return resp, errors.New("must have routed paychannel first")
	}

	// update merchant paychannel status
	err = mr.merchantRepoWrites.UpdateStatusMerchantPaychannelById(payload.MerchantPaychannelId, status)
	if err != nil {
		slog.Infof("merchant paychannel ID %v got err when update status %v", payload.MerchantPaychannelId, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "Please contact support for more info",
		}
		return resp, err
	}

	msg := fmt.Sprintf("Merchant paychannel id %v, success to %v", payload.MerchantPaychannelId, status)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
	}

	return resp, nil
}

func (mr *Merchant) AddChannelSvc(id int, payload []dto.PaymentMethodData) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if len(payload) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "There is no payment method choose",
		}
		return resp, errors.New("need choose at least one payment method")
	}

	for _, method := range payload {
		// create payment method
		paymentMethodId, err := mr.merchantRepoWrites.CreateMerchantPaymentMethodRepo(id, method.PaymentMethodId)

		if err != nil {
			msg := fmt.Sprintf("failed to create payment method with err: %v", err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: msg,
			}
			return resp, err
		}

		// create merchant paychannel
		segment := "MAIN"
		methodStr := constant.TransformPaymentMethodNameIntoCode[method.PaymentMethodName]
		randomStrMerchantChannelCode := helper.GenerateRandomString(15)
		merchantPaychannelCode := methodStr + "-" + strings.ToUpper(segment) + "-" + randomStrMerchantChannelCode
		minAmount := method.MinAmountPerTransaction
		maxAmount := method.MaxAmountPerTransaction
		dailyLimit := method.DailyLimit
		fee := method.Fee
		feeType := method.FeeType
		if feeType == "" {
			feeType = constant.FeeTypeFixedFee
		}

		_, err = mr.merchantRepoWrites.CreateMerchantPaychannelRepo(paymentMethodId, segment, fee, feeType, minAmount, maxAmount, dailyLimit, merchantPaychannelCode)
		if err != nil {
			msg := fmt.Sprintf("failed to create merchant paychannel with err: %v", err.Error())
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: msg,
			}
			return resp, err
		}
	}

	msg := fmt.Sprintf("Success create merchant paychannel for id merchant: %v", id)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
	}

	return resp, nil
}

func (mr *Merchant) GetListMerchantPaymentMethodsSvc(idMerchant int) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var listPaymentMethodFilter []entity.PaymentMethods

	listMerchantPaymentMethods, listPaymentMethods, err := mr.merchantRepoReads.GetMerchantPaymentMethodByIdMerchantRepo(idMerchant)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	existingPaymentMethod := make(map[string]bool)
	for _, merchantMethod := range listMerchantPaymentMethods {
		existingPaymentMethod[merchantMethod.Name] = true
	}

	// filter list payment methods
	for _, paymentMethod := range listPaymentMethods {
		if !existingPaymentMethod[paymentMethod.Name] {
			listPaymentMethodFilter = append(listPaymentMethodFilter, paymentMethod)
		}
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success get list payment method",
		Data:            listPaymentMethodFilter,
	}

	return resp, nil
}

func (mr *Merchant) AddRoutingPaychannelSvc(merchantPaychannelId int, payload dto.AddPaychannelRouting) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// delete paychannel routing
	err := mr.merchantRepoWrites.DeleteRoutingPaychannelByMerchantPaychannelId(merchantPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// add paychannel routing
	for _, providerPaychannelId := range payload.ProviderPaychannelId {
		_, err := mr.merchantRepoWrites.AddRoutingPaychannelRepo(merchantPaychannelId, providerPaychannelId)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success updated routing paychannel",
	}

	return resp, nil
}

func (mr *Merchant) GetActiveAvailablePaychannelSvc(merchantPaychannelId int) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// get merchant paychannel detail
	merchantPaychannelDetail, err := mr.merchantRepoReads.GetMerchantPaychannelDetailById(merchantPaychannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, nil
	}

	activePaychannel, availableChannel, err := mr.merchantRepoReads.GetActiveAndAvailableChannelRepo(merchantPaychannelId, merchantPaychannelDetail.PaymentMethodChannel)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	var filterAvailableChannel []entity.ProviderPaychannelEntity

	exisitingChannelId := make(map[int]bool)
	for _, activeChannel := range activePaychannel {
		exisitingChannelId[activeChannel.Id] = true
	}

	for _, availChannel := range availableChannel {
		if !exisitingChannelId[availChannel.Id] {
			filterAvailableChannel = append(filterAvailableChannel, availChannel)
		}
	}

	activeAvailableChannelResp := dto.ActiveAvailableChannelRespDto{
		ActivePaychannel:    activePaychannel,
		AvailablePaychannel: filterAvailableChannel,
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve active and available paychannel",
		Data:            activeAvailableChannelResp,
	}

	return resp, nil
}

func (mr *Merchant) HomeAnalyticsSvc(payload dto.HomeAnalyticsDto) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if payload.MinDate == "" {
		payload.MinDate = helper.GenerateTime(0)
	}

	if payload.MaxDate == "" {
		payload.MaxDate = helper.GenerateTime(24)
	}

	users, err := mr.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	requestData := dto.GetMerchantAnalyticsDtoReq{
		MinDate:    payload.MinDate,
		MaxDate:    payload.MaxDate,
		MerchantId: *users.MerchantID,
	}

	transactionData, err := mr.transactionRepoReads.GetTransactionAnalyticsRepo(requestData)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	data := supportHomeAnalyticsSvc(transactionData)

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success analytics data merchant",
		Data:            data,
	}

	return resp, nil
}

func (mr *Merchant) GetMerchantInformationSvc(merchantId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	merchantData, err := mr.merchantRepoReads.GetMerchantDataByMerchantId(merchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            merchantData,
	}

	return resp, nil
}

func (mr *Merchant) DisplaySecretKeySvc(pin string, username string, merchantId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	user, err := mr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	secretKey, err := mr.merchantRepoReads.GetSecretKeyByMerchantIdRepo(merchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve secret key",
		Data:            secretKey,
	}

	return resp, nil
}

func (mr *Merchant) GenerateSecretKeySvc(pin string, username string, merchantId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	randomStrMerchantScret := helper.GenerateRandomString(30)
	merchantSecret := "secret_key-" + randomStrMerchantScret

	user, err := mr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	err = mr.merchantRepoWrites.UpdateMerchantSecretKeyRepo(merchantSecret, merchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success generated new secret key",
	}

	return resp, nil
}

func (mr *Merchant) GetMerchantAccountBalanceSvc(username string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	user, err := mr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	accountBalance, err := mr.merchantRepoReads.GetMerchantAccountByMerchantId(*user.MerchantID)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            accountBalance,
	}

	return resp, nil
}

func (mr *Merchant) GetInformationMerchantSvc(username string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	user, err := mr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		slog.Errorw("failed get user data", "stack_trace", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	merchantData, err := mr.merchantRepoReads.GetMerchantDataByMerchantId(*user.MerchantID)
	if err != nil {
		slog.Errorw("failed get merchant data", "stack_trace", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            merchantData,
	}

	return resp, nil
}

func (mr *Merchant) DisplayMerchantKeySvc(username string, pin string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	user, err := mr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		slog.Errorw("failed get data user", "stack_trace", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	secretKey, err := mr.merchantRepoReads.GetSecretKeyByMerchantIdRepo(*user.MerchantID)
	if err != nil {
		slog.Errorw("failed get merchant key", "stack_trace", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve secret key",
		Data:            secretKey,
	}

	return resp, nil
}

func (mr *Merchant) GenerateMerchantKeySvc(pin string, username string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	randomStrMerchantScret := helper.GenerateRandomString(30)
	merchantSecret := "secret_key-" + randomStrMerchantScret

	user, err := mr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// check input pin
	if !comparePasswords(user.Pin, []byte(pin)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "wrong pin",
		}

		return resp, errors.New("wrong pin")
	}

	err = mr.merchantRepoWrites.UpdateMerchantSecretKeyRepo(merchantSecret, *user.MerchantID)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success generated new secret key",
	}

	return resp, nil
}

func supportMerchantAnalyticsSvc(payload []entity.PaymentDetailMerchantProvider) dto.AnalyticsMerchantRespDto {
	var totalVolumeSuccessIn float64
	var totalSuccessTransactionIn int
	var totalFailedTransactionIn int
	var totalProcessingTransactionIn int
	var totalTransactionIn int
	var successRateIn float64
	var totalDurationIn time.Duration
	var totalCompletedIn int

	var totalVolumeSuccessOut float64
	var totalSuccessTransactionOut int
	var totalFailedTransactionOut int
	var totalProcessingTransactionOut int
	var totalTransactionOut int
	var successRateOut float64
	var totalDurationOut time.Duration
	var totalCompletedOut int

	for _, transaction := range payload {
		if transaction.PayType == constant.PayTypePayin {
			totalTransactionIn++
			if transaction.Status == constant.StatusSuccess {
				totalVolumeSuccessIn += transaction.TransactionAmount
				totalSuccessTransactionIn++

				// completion rate
				var processingTime time.Time
				var successTime time.Time

				processingTime = transaction.TransactionCreatedAt
				successTime = transaction.TransactionUpdatedAt

				if !processingTime.IsZero() && !successTime.IsZero() {
					duration := successTime.Sub(processingTime)
					totalDurationIn += duration
					totalCompletedIn++
				}
			}

			if transaction.Status == constant.StatusFailed {
				totalFailedTransactionIn++
			}

			if transaction.Status == constant.StatusProcessing {
				totalProcessingTransactionIn++
			}
		}

		if transaction.PayType == constant.PayTypePayout {
			totalTransactionOut++
			if transaction.Status == constant.StatusSuccess {
				totalVolumeSuccessOut += transaction.TransactionAmount
				totalSuccessTransactionOut++

				// completion rate
				var processingTime time.Time
				var successTime time.Time

				processingTime = transaction.TransactionCreatedAt
				successTime = transaction.TransactionUpdatedAt

				if !processingTime.IsZero() && !successTime.IsZero() {
					duration := successTime.Sub(processingTime)
					totalDurationOut += duration
					totalCompletedOut++
				}
			}

			if transaction.Status == constant.StatusFailed {
				totalFailedTransactionOut++
			}

			if transaction.Status == constant.StatusProcessing {
				totalProcessingTransactionOut++
			}
		}
	}

	// Calculate the average time between `PROCESSING` and `SUCCESS` (in)
	var averageDuration time.Duration
	var formattedCompletionIn string
	var formattedCompletionOut string

	if totalCompletedIn > 0 {
		averageDuration = totalDurationIn / time.Duration(totalCompletedIn)
		formattedCompletionIn = converter.FormattedCompletionRate(averageDuration)
	}

	if totalCompletedOut > 0 {
		averageDuration = totalDurationOut / time.Duration(totalCompletedOut)
		formattedCompletionOut = converter.FormattedCompletionRate(averageDuration)
	}

	// Calculate success rate for out transactions
	if totalSuccessTransactionOut > 0 {
		successRateOut = math.Ceil((float64(totalSuccessTransactionOut) / float64(totalTransactionOut)) * 100)
	}

	// Calculate success rate for in transactions
	if totalSuccessTransactionIn > 0 {
		successRateIn = math.Ceil((float64(totalSuccessTransactionIn) / float64(totalTransactionIn)) * 100)
	}

	succesRateInFormatted := helper.FormattedUsingPercent(successRateIn)
	succesRateOutFormatted := helper.FormattedUsingPercent(successRateOut)

	// Prepare response data
	inAnalyticsData := dto.AnalyticsDataRespDto{
		TotalVolume:        totalVolumeSuccessIn,
		SuccessRate:        succesRateInFormatted,
		CompletionRate:     formattedCompletionIn,
		TransactionTotal:   totalTransactionIn,
		SuccessTransaction: totalSuccessTransactionIn,
		FailedTransaction:  totalFailedTransactionIn,
	}

	outAnalyticsData := dto.AnalyticsDataRespDto{
		TotalVolume:        totalVolumeSuccessOut,
		SuccessRate:        succesRateOutFormatted,
		CompletionRate:     formattedCompletionOut,
		TransactionTotal:   totalTransactionOut,
		SuccessTransaction: totalSuccessTransactionOut,
		FailedTransaction:  totalFailedTransactionOut,
	}

	merchantAnalyticsDataRes := dto.AnalyticsMerchantRespDto{
		TransactionIn:  inAnalyticsData,
		TransactionOut: outAnalyticsData,
	}

	return merchantAnalyticsDataRes
}

func supportMerchantAnalyticsByMerchantPaychannelSvc(payload []entity.PaymentDetailMerchantProvider) dto.AnalyticsDataRespDto {
	var totalVolumeSuccess float64
	var successRate float64
	var totalSuccess int
	var totalFailed int
	var totalProcessing int
	var totalTransaction int
	var totalDuration time.Duration
	var totalCompleted int

	for _, transaction := range payload {
		totalTransaction++
		if transaction.Status == constant.StatusSuccess {
			totalVolumeSuccess += transaction.TransactionAmount
			totalSuccess++

			// completion rate
			var processingTime time.Time
			var successTime time.Time

			processingTime = transaction.TransactionCreatedAt
			successTime = transaction.TransactionUpdatedAt

			if !processingTime.IsZero() && !successTime.IsZero() {
				duration := successTime.Sub(processingTime)
				totalDuration += duration
				totalCompleted++
			}
		}

		if transaction.Status == constant.StatusFailed {
			totalFailed++
		}

		if transaction.Status == constant.StatusProcessing {
			totalProcessing++
		}
	}

	var averageDuration time.Duration
	var formattedCompletion string

	if totalCompleted > 0 {
		averageDuration = totalDuration / time.Duration(totalCompleted)
		formattedCompletion = converter.FormattedCompletionRate(averageDuration)
	}

	// calculate success rate
	if totalSuccess > 0 {
		successRate = math.Ceil((float64(totalSuccess) / float64(totalTransaction)) * 100)
	}

	successRateFormatted := helper.FormattedUsingPercent(successRate)

	// response data
	analyticsDataPaychannel := dto.AnalyticsDataRespDto{
		TotalVolume:        totalVolumeSuccess,
		SuccessRate:        successRateFormatted,
		CompletionRate:     formattedCompletion,
		TransactionTotal:   totalTransaction,
		SuccessTransaction: totalSuccess,
		FailedTransaction:  totalFailed,
	}

	return analyticsDataPaychannel
}

func supportHomeAnalyticsSvc(payload []entity.PaymentDetailMerchantProvider) dto.HomeAnalyticsRespDto {
	var totalNumberTransactionIn int
	var totalAmountTransactionIn float64
	var totalNumberTransactionOut int
	var totalAmountTransactionOut float64

	// Qris
	var totalNumberQris int
	var totalAmountQris float64

	// ewallet
	var totalNumberEwallet int
	var totalAmountEwallet float64

	// virtual account
	var totalNumberVa int
	var totalAmountVa float64

	// disbursement
	var totalNumberDisbursement int
	var totalAmountDisbursement float64

	for _, transaction := range payload {
		if transaction.PayType == constant.PayTypePayin {
			totalNumberTransactionIn++
			totalAmountTransactionIn += transaction.TransactionAmount

			if transaction.PaymentMethodName == constant.QrisPaymentMethod {
				totalNumberQris++
				totalAmountQris += transaction.TransactionAmount
			}

			if transaction.PaymentMethodName == constant.VirtualAccountPaymentMethod {
				totalNumberVa++
				totalAmountVa += transaction.TransactionAmount
			}

			if transaction.PaymentMethodName == constant.EwalletPaymentMethod {
				totalNumberEwallet++
				totalAmountEwallet += transaction.TransactionAmount
			}
		}

		if transaction.PayType == constant.PayTypePayout {
			totalNumberTransactionOut++
			totalAmountTransactionOut += transaction.TransactionAmount

			if transaction.PaymentMethodName == constant.DisbursementPaymentMethod {
				totalNumberDisbursement++
				totalAmountDisbursement += transaction.TransactionAmount
			}
		}
	}

	qris := dto.HomeAnalyticsDataRespDto{
		TotalNumber: totalNumberQris,
		TotalAmount: totalAmountQris,
	}

	ewallet := dto.HomeAnalyticsDataRespDto{
		TotalNumber: totalNumberEwallet,
		TotalAmount: totalAmountEwallet,
	}

	virtualAccount := dto.HomeAnalyticsDataRespDto{
		TotalNumber: totalNumberVa,
		TotalAmount: totalAmountVa,
	}

	disbursement := dto.HomeAnalyticsDataRespDto{
		TotalNumber: totalNumberDisbursement,
		TotalAmount: totalAmountDisbursement,
	}

	transactionIn := dto.HomeAnalyticsDataRespDto{
		TotalNumber: totalNumberTransactionIn,
		TotalAmount: totalAmountTransactionIn,
	}

	transactionOut := dto.HomeAnalyticsDataRespDto{
		TotalNumber: totalNumberTransactionOut,
		TotalAmount: totalAmountTransactionOut,
	}

	resposeHome := dto.HomeAnalyticsRespDto{
		TransactionIn:  transactionIn,
		TransactionOut: transactionOut,
		TotalSuccessPayment: dto.SuccessPayment{
			Qris:           qris,
			Ewallet:        ewallet,
			VirtualAccount: virtualAccount,
			Disbursement:   disbursement,
		},
	}

	return resposeHome
}
