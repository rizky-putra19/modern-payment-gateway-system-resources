package service

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal"
	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
)

type Transaction struct {
	transactionRepoReads  internal.TransactionsReadsRepositoryItf
	transactionRepoWrites internal.TransactionsWritesRepositoryItf
	merchantRepoReads     internal.MerchantReadsRepositoryItf
	merchantRepoWrites    internal.MerchantWritesRepositoryItf
	providerRepoReads     internal.ProviderReadsRepositoryItf
	providerRepoWrites    internal.ProviderWritesRepositoryItf
	userRepoReads         internal.UserReadsRepositoryItf
	configApp             config.App
	jackProvider          internal.JackProviderItf
	regex                 *regexp.Regexp
}

func NewTransaction(
	transactionRepoReads internal.TransactionsReadsRepositoryItf,
	transactionRepoWrites internal.TransactionsWritesRepositoryItf,
	userRepoReads internal.UserReadsRepositoryItf,
	merchantRepoReads internal.MerchantReadsRepositoryItf,
	merchantRepoWrites internal.MerchantWritesRepositoryItf,
	configApp config.App,
	jackProvider internal.JackProviderItf,
	providerRepoReads internal.ProviderReadsRepositoryItf,
	providerRepoWrites internal.ProviderWritesRepositoryItf,
) *Transaction {
	// regex only allow string
	reg, _ := regexp.Compile("[^a-zA-Z]+")
	return &Transaction{
		transactionRepoReads:  transactionRepoReads,
		transactionRepoWrites: transactionRepoWrites,
		merchantRepoReads:     merchantRepoReads,
		merchantRepoWrites:    merchantRepoWrites,
		userRepoReads:         userRepoReads,
		configApp:             configApp,
		jackProvider:          jackProvider,
		providerRepoReads:     providerRepoReads,
		providerRepoWrites:    providerRepoWrites,
		regex:                 reg,
	}
}

func (tr *Transaction) GetStatusChangeLog(paymentId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	statusChangeLogsData, err := tr.transactionRepoReads.GetStatusChangeLogData(paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(statusChangeLogsData) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "does not have status change log yet",
			Data:            statusChangeLogsData,
		}
		return resp, nil
	}

	// Create a new slice to hold the modified data
	modifiedStatusChangeLogsData := make([]entity.TransactionStatusLogs, len(statusChangeLogsData))
	copy(modifiedStatusChangeLogsData, statusChangeLogsData)

	// Change response with real message for operations dashboard only
	for i := range modifiedStatusChangeLogsData {
		modifiedStatusChangeLogsData[i].Notes = modifiedStatusChangeLogsData[i].RealNotes
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve status change logs",
		Data:            modifiedStatusChangeLogsData,
	}

	return resp, nil
}

func (tr *Transaction) GetPaymentDetailConfirmData(paymentId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	detailConfirmData, err := tr.transactionRepoReads.GetPaymentDetailPrvConfirmDetail(paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(detailConfirmData) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "does not have confirmation data yet from provider",
			Data:            detailConfirmData,
		}
		return resp, nil
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve provider confirmation detail data",
		Data:            detailConfirmData,
	}

	return resp, nil
}

func (tr *Transaction) GetTransactionList(params dto.QueryParams) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if params.MinDate == "" {
		params.MinDate = helper.GenerateTime(0)
	}

	if params.MaxDate == "" {
		params.MaxDate = helper.GenerateTime(24)
	}

	// get transaction list to repository
	transactionList, pagination, err := tr.transactionRepoReads.GetTransactionList(params)
	if err != nil {
		slog.Infof("username %v got failed get transaction list with err: %v", params.Username, err.Error())
		return dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}, err
	}

	if len(transactionList) < 1 {
		return dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "transaction not found",
			Data:            transactionList,
		}, nil
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "successfully get transaction",
		Data:            transactionList,
		Pagination:      pagination,
	}

	return resp, nil
}

func (tr *Transaction) GetTransactionInListSvc(params dto.QueryParams) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	users, err := tr.userRepoReads.GetUserByUsername(params.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if params.MinDate == "" {
		params.MinDate = helper.GenerateTime(0)
	}

	if params.MaxDate == "" {
		params.MaxDate = helper.GenerateTime(24)
	}

	params.MerchantId = *users.MerchantID
	transactionInList, pagination, err := tr.transactionRepoReads.GetTransactionInListRepo(params)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(transactionInList) < 1 {
		return dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "transaction not found",
			Data:            transactionInList,
		}, nil
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "successfully get transaction",
		Data:            transactionInList,
		Pagination:      pagination,
	}

	return resp, nil
}

func (tr *Transaction) GetTransactionOutListSvc(params dto.QueryParams) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	users, err := tr.userRepoReads.GetUserByUsername(params.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if params.MinDate == "" {
		params.MinDate = helper.GenerateTime(0)
	}

	if params.MaxDate == "" {
		params.MaxDate = helper.GenerateTime(24)
	}

	params.MerchantId = *users.MerchantID
	transactionInList, pagination, err := tr.transactionRepoReads.GetTransactionOutListRepo(params)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(transactionInList) < 1 {
		return dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "transaction not found",
			Data:            transactionInList,
		}, nil
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "successfully get transaction",
		Data:            transactionInList,
		Pagination:      pagination,
	}

	return resp, nil
}

func (tr *Transaction) GetPaymentDetailCapitalFlow(paymentId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var transactionDataFlows []dto.TransactionsCapitalFlow

	// get capital flow data from database
	capitalFlows, err := tr.transactionRepoReads.GetCapitalFlowsWitPaymentId(paymentId)
	if err != nil {
		return dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}, nil
	}

	if len(capitalFlows) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "capital flows didn't retrieve maybe wrong payment id or status transaction not yet SUCCESS or FAILED",
		}
		return resp, nil
	}

	for _, capitalData := range capitalFlows {
		transactionsData := dto.TransactionsCapitalFlow{
			TransactionType: capitalData.ReasonName,
			MerchantAccount: capitalData.MerchantId,
			Amount:          capitalData.Amount,
			Status:          capitalData.Status,
			CapitalType:     capitalData.CapitalType,
			CreatedAt:       capitalData.CreatedAt,
		}
		transactionDataFlows = append(transactionDataFlows, transactionsData)
	}

	msg := fmt.Sprintf("success get capital flows for payment id %v", paymentId)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
		Data:            transactionDataFlows,
	}

	return resp, nil
}

func (tr *Transaction) GetPaymentDetailMerchantProvider(req dto.GetPaymentDetailsRequest) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// get payment detail to database
	paymentData, err := tr.transactionRepoReads.GetPaymentDetailProviderMerchant(req.PaymentId)
	if err != nil {
		return dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}, nil
	}

	// Handle potential nil values with default values
	segment := ""
	if paymentData.Segment != nil {
		segment = *paymentData.Segment
	}

	interfaceSetting := ""
	if paymentData.InterfaceSetting != nil {
		interfaceSetting = *paymentData.InterfaceSetting
	}

	merchantData := dto.DetailData{
		Name:            paymentData.MerchantName,
		PaymentId:       paymentData.MerchantRefNumber,
		Channel:         paymentData.MerchantPaychannelCode,
		Segment:         segment,
		Fee:             paymentData.MerchantFee,
		FeeType:         paymentData.MerchantFeeType,
		ClientIpAddress: paymentData.ClientIPAddress,
		RequestMethod:   paymentData.RequestMethod,
		RequestedAmount: paymentData.TransactionAmount,
	}

	providerData := dto.DetailData{
		Name:             paymentData.ProviderName,
		PaymentId:        paymentData.ProviderRefNumber,
		Channel:          paymentData.PaychannelName,
		Fee:              paymentData.ProviderFee,
		FeeType:          paymentData.ProviderFeeType,
		InterfaceSetting: interfaceSetting,
	}

	paymentDetailDataResp := dto.PaymentDetailProviderMerchantResponse{
		MerchantData: merchantData,
		ProviderData: providerData,
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "successfully get payment detail provider and merchant",
		Data:            paymentDetailDataResp,
	}

	return resp, nil
}

func (tr *Transaction) GetPaymentDetailAccountInformation(req dto.GetPaymentDetailsRequest) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var sendingBankInformation entity.AccountData
	var receiveBankInformation entity.AccountData

	accountData, err := tr.transactionRepoReads.GetPaymentDetailAccountInformation(req.PaymentId)
	if err != nil {
		return dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}, nil
	}

	if len(accountData) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadGateway,
			ResponseMessage: "account information didn't retrieve maybe wrong payment id",
		}
		return resp, nil
	}

	// Iterate through accountData to find and assign the correct account information
	for _, account := range accountData {
		if account.AccountType != nil && *account.AccountType == constant.AccountTypeDebitor {
			sendingBankInformation = account
		} else if account.AccountType != nil && *account.AccountType == constant.AccountTypeCreditor {
			receiveBankInformation = account
		}
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success get account information",
		Data: map[string]interface{}{
			"sendingBankInformation": sendingBankInformation,
			"receiveBankInformation": receiveBankInformation,
		},
	}

	return resp, nil
}

func (tr *Transaction) GetPaymentDetailFee(paymentId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var feeResp dto.FeeResponse

	// get payment detail to repository
	paymentData, err := tr.transactionRepoReads.GetPaymentDetailProviderMerchant(paymentId)
	if err != nil {
		return dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "unprocessable entity",
		}, nil
	}

	// calculate fee for percentage
	if paymentData.MerchantFeeType == constant.FeeTypePercentage {
		calculatedFeeMerchant := math.Ceil(paymentData.TransactionAmount * (paymentData.MerchantFee / 100))
		calculatedFeeProvider := math.Ceil(paymentData.TransactionAmount * (paymentData.ProviderFee / 100))

		merchantFee := dto.FeeData{
			ConfiguredFee: paymentData.MerchantFee,
			CalculatedFee: calculatedFeeMerchant,
			ChargedFee:    calculatedFeeMerchant,
		}

		providerFee := dto.FeeData{
			ConfiguredFee: paymentData.ProviderFee,
			CalculatedFee: calculatedFeeProvider,
		}

		calculatedProfit := calculatedFeeMerchant - calculatedFeeProvider

		feeResp = dto.FeeResponse{
			MerchantFee:      merchantFee,
			ProviderFee:      providerFee,
			CalculatedProfit: calculatedProfit,
			FeeType:          paymentData.MerchantFeeType,
		}

		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "success calculate fee trasaction",
			Data:            feeResp,
		}
	}

	// calculate fee for fixed fee
	if paymentData.MerchantFeeType == constant.FeeTypeFixedFee {
		merchantFee := dto.FeeData{
			ConfiguredFee: paymentData.MerchantFee,
			CalculatedFee: paymentData.MerchantFee,
			ChargedFee:    paymentData.MerchantFee,
		}

		providerFee := dto.FeeData{
			ConfiguredFee: paymentData.ProviderFee,
			CalculatedFee: paymentData.ProviderFee,
		}

		calculatedProfit := paymentData.MerchantFee - paymentData.ProviderFee

		feeResp = dto.FeeResponse{
			MerchantFee:      merchantFee,
			ProviderFee:      providerFee,
			CalculatedProfit: calculatedProfit,
			FeeType:          paymentData.MerchantFeeType,
		}

		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "success calculate fee trasaction",
			Data:            feeResp,
		}
	}

	return resp, nil
}

func (tr *Transaction) UpdateStatusTransaction(paymentId string, status string, username string, notes string, pin string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// user data
	user, err := tr.userRepoReads.GetUserByUsername(username)
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

	transactionData, err := tr.transactionRepoReads.GetPaymentDetailProviderMerchant(paymentId)
	if err != nil {
		return dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}, err
	}

	if transactionData.Status == strings.ToUpper(status) {
		msg := fmt.Sprintf("status already %v", status)
		return dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: msg,
		}, nil
	}

	if transactionData.TransactionID == 0 {
		msg := fmt.Sprintf("transaction id %v didn't found", paymentId)
		return dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: msg,
		}, nil
	}

	if transactionData.PayType == constant.PayTypePayin && strings.ToUpper(status) == constant.StatusSuccess {
		// validate if transaction have been success - failed, and would like to change into success
		transactionStatusLogs, err := tr.transactionRepoReads.GetStatusChangeLogData(transactionData.PaymentID)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}

		if transactionStatusLogs[1].StatusLog == constant.StatusLogSuccess {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusBadRequest,
				ResponseMessage: "can't updated into success due to have been success before",
			}
			return resp, nil
		}
	}

	// update transaction status
	err = tr.transactionRepoWrites.UpdateStatus(strings.ToUpper(status), paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// handle for status SUCCESS
	if strings.ToUpper(status) == constant.StatusSuccess {
		// update status change log
		_, err := tr.transactionRepoWrites.CreateTransactionStatusLog(paymentId, constant.StatusLogSuccess, username, notes, "")
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}

		// handle for payin and SUCCESS
		if transactionData.PayType == constant.PayTypePayin {
			// get merchant account
			merchantAccountData, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			balanceCapitalFlow := merchantAccountData.BalanceCapitalFlow
			notSettledBalance := merchantAccountData.NotSettledBalance

			// payin amount to balance
			balanceCapitalAddPayin := balanceCapitalFlow + transactionData.TransactionAmount
			notSettledBalanceAddPayin := notSettledBalance + transactionData.TransactionAmount
			formattedNotSettledBalanceAddPayin := helper.FormatFloat64(notSettledBalanceAddPayin)
			formattedBalanceAddPayin := helper.FormatFloat64(balanceCapitalAddPayin)
			// update merchant balance
			err = tr.merchantRepoWrites.UpdateMerchantCapitalAndNotSettleBalance(formattedNotSettledBalanceAddPayin, formattedBalanceAddPayin, transactionData.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// create merchant capital flow for transaction amount
			payloadMerchantCapitalFlowPayin := dto.CreateMerchantCapitalFlowPayload{
				PaymentId:         transactionData.PaymentID,
				MerchantAccountId: merchantAccountData.Id,
				TempBalance:       formattedBalanceAddPayin,
				ReasonId:          constant.ReasonIdPayin,
				Status:            strings.ToUpper(status),
				CreateBy:          constant.CreateBySystem,
				Amount:            helper.FormatFloat64(transactionData.TransactionAmount),
				CapitalType:       constant.CapitalTypeCredit,
			}
			_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowPayin)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// calculate fee for merchant capital flow
			var feeTransactions float64
			if transactionData.MerchantFeeType == constant.FeeTypePercentage {
				calculateFee := math.Ceil(transactionData.TransactionAmount * (transactionData.MerchantFee / 100))
				feeTransactions = helper.FormatFloat64(calculateFee)
			} else {
				feeTransactions = helper.FormatFloat64(transactionData.MerchantFee)
			}

			// get merchant account after add payin
			merchantAccountDataUpdated, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			balanceCapitalFlow = merchantAccountDataUpdated.BalanceCapitalFlow
			notSettledBalance = merchantAccountDataUpdated.NotSettledBalance

			// adjust balance minus fee
			balanceCapitalMinusFee := balanceCapitalFlow - feeTransactions
			notSettledBalanceMinusFee := notSettledBalance - feeTransactions
			formattedNotSettledBalanceMinusFee := helper.FormatFloat64(notSettledBalanceMinusFee)
			formattedBalanceMinusFee := helper.FormatFloat64(balanceCapitalMinusFee)

			err = tr.merchantRepoWrites.UpdateMerchantCapitalAndNotSettleBalance(formattedNotSettledBalanceMinusFee, formattedBalanceMinusFee, transactionData.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// create merchant capital flow for fee
			payloadMerchantCapitalFlowFee := dto.CreateMerchantCapitalFlowPayload{
				PaymentId:         transactionData.PaymentID,
				MerchantAccountId: merchantAccountData.Id,
				TempBalance:       formattedBalanceMinusFee,
				ReasonId:          constant.ReasonIdFee,
				Status:            strings.ToUpper(status),
				CreateBy:          constant.CreateBySystem,
				Amount:            feeTransactions,
				CapitalType:       constant.CapitalTypeDebit,
			}
			_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowFee)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}
		}

		// handle for payout and SUCCESS
		if transactionData.PayType == constant.PayTypePayout {
			// get merchant account
			merchantAccountData, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// check if transaction before is failed
			if transactionData.Status == constant.StatusFailed {
				balanceCapitalFlow := merchantAccountData.BalanceCapitalFlow
				settleBalance := merchantAccountData.SettledBalance

				// adjust settle balance minum transaction amount
				balanceCapitalMinusPayout := balanceCapitalFlow - float64(transactionData.TransactionAmount)
				settleBalanceMinusPayout := settleBalance - float64(transactionData.TransactionAmount)
				formattedBalanceCapitalMinusPayout := helper.FormatFloat64(balanceCapitalMinusPayout)
				formattedSettleBalanceMinusPayout := helper.FormatFloat64(settleBalanceMinusPayout)

				// update merchant balance
				err = tr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalanceMinusPayout, formattedBalanceCapitalMinusPayout, transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// create merchant capital flow for transaction amount
				payloadMerchantCapitalFlowPayout := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       formattedBalanceCapitalMinusPayout,
					ReasonId:          constant.ReasonIdPayout,
					Status:            strings.ToUpper(status),
					CreateBy:          constant.CreateBySystem,
					Amount:            helper.FormatFloat64(transactionData.TransactionAmount),
					CapitalType:       constant.CapitalTypeDebit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowPayout)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// get merchant account after minus payout
				merchantAccountDataUpdated, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				balanceCapitalFlow = merchantAccountDataUpdated.BalanceCapitalFlow
				settleBalance = merchantAccountDataUpdated.SettledBalance

				// fee transaction
				feeTransaction := helper.FormatFloat64(transactionData.MerchantFee)

				// adjust balance minus fee
				balanceCapitalMinusFee := balanceCapitalFlow - feeTransaction
				settleBalanceMinusFee := settleBalance - feeTransaction
				formattedBalanceCapitalMinusFee := helper.FormatFloat64(balanceCapitalMinusFee)
				formattedSettleBalanceMinusFee := helper.FormatFloat64(settleBalanceMinusFee)

				err = tr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalanceMinusFee, formattedBalanceCapitalMinusFee, transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// create merchant capital flow for fee
				payloadMerchantCapitalFlowFee := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       formattedBalanceCapitalMinusFee,
					ReasonId:          constant.ReasonIdFee,
					Status:            strings.ToUpper(status),
					CreateBy:          constant.CreateBySystem,
					Amount:            feeTransaction,
					CapitalType:       constant.CapitalTypeDebit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowFee)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}
			}

			if transactionData.Status == constant.StatusProcessing {
				balanceCapitalFlow := merchantAccountData.BalanceCapitalFlow
				pendingOutBalance := merchantAccountData.PendingTransactionOut

				// payout amount to balance
				balanceCapitalMinusPayout := balanceCapitalFlow - transactionData.TransactionAmount
				pendingOutBalanceMinusPayout := pendingOutBalance - transactionData.TransactionAmount
				formattedPendingOutBalanceMinusPayout := helper.FormatFloat64(pendingOutBalanceMinusPayout)
				formattedBalanceMinusPayout := helper.FormatFloat64(balanceCapitalMinusPayout)

				// update merchant balance
				err = tr.merchantRepoWrites.UpdateMerchantCapitalPendingOut(formattedPendingOutBalanceMinusPayout, formattedBalanceMinusPayout, transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// create merchant capital flow for transaction amount
				payloadMerchantCapitalFlowPayout := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       formattedBalanceMinusPayout,
					ReasonId:          constant.ReasonIdPayout,
					Status:            strings.ToUpper(status),
					CreateBy:          constant.CreateBySystem,
					Amount:            helper.FormatFloat64(transactionData.TransactionAmount),
					CapitalType:       constant.CapitalTypeDebit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowPayout)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// fee payout merchant
				feeTransaction := helper.FormatFloat64(transactionData.MerchantFee)

				// get merchant account after minus payout
				merchantAccountDataUpdated, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				balanceCapitalFlow = merchantAccountDataUpdated.BalanceCapitalFlow
				pendingOutBalance = merchantAccountDataUpdated.PendingTransactionOut

				// adjust balance minus fee
				balanceCapitalMinusFee := balanceCapitalFlow - feeTransaction
				pendingPayoutMinusFee := pendingOutBalance - feeTransaction
				formattedPendingPayoutBalance := helper.FormatFloat64(pendingPayoutMinusFee)
				formattedBalanceMinusFee := helper.FormatFloat64(balanceCapitalMinusFee)

				err = tr.merchantRepoWrites.UpdateMerchantCapitalPendingOut(formattedPendingPayoutBalance, formattedBalanceMinusFee, transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// create merchant capital flow for fee
				payloadMerchantCapitalFlowFee := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       formattedBalanceMinusFee,
					ReasonId:          constant.ReasonIdFee,
					Status:            strings.ToUpper(status),
					CreateBy:          constant.CreateBySystem,
					Amount:            feeTransaction,
					CapitalType:       constant.CapitalTypeDebit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowFee)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}
			}
		}
	}
	// end of blocked status handler SUCCESS

	// start blocked code for handle for status FAILED
	if strings.ToUpper(status) == constant.StatusFailed {

		// update status change log
		_, err := tr.transactionRepoWrites.CreateTransactionStatusLog(paymentId, constant.StatusLogFailed, username, notes, "")
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}

		// handle for payin and FAILED
		if transactionData.PayType == constant.PayTypePayin {
			// get merchant account
			merchantAccountData, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// check if transaction have SUCCESS status before
			if transactionData.Status == constant.StatusSuccess {
				balanceUse := constant.SettleBalance

				balanceCapitalFlow := merchantAccountData.BalanceCapitalFlow
				settleOrNotSettleBalance := merchantAccountData.SettledBalance
				if settleOrNotSettleBalance < transactionData.TransactionAmount {
					balanceUse = constant.NotSettledBalance
					settleOrNotSettleBalance = merchantAccountData.NotSettledBalance
				}

				// calculate fee for merchant capital flow
				var feeTransactions float64
				if transactionData.MerchantFeeType == constant.FeeTypePercentage {
					calculateFee := math.Ceil(transactionData.TransactionAmount * (transactionData.MerchantFee / 100))
					feeTransactions = helper.FormatFloat64(calculateFee)
				} else {
					feeTransactions = helper.FormatFloat64(transactionData.MerchantFee)
				}

				// transaction amount will minus fee first and take out from balance
				transactionAmountMinusFee := transactionData.TransactionAmount - feeTransactions

				balanceCapitalMinusFee := balanceCapitalFlow - transactionAmountMinusFee
				settleOrNotSettleBalanceMinusFee := settleOrNotSettleBalance - transactionAmountMinusFee
				formattedNotSettledOrSettleBalanceMinusFee := helper.FormatFloat64(settleOrNotSettleBalanceMinusFee)
				formattedBalanceMinusFee := helper.FormatFloat64(balanceCapitalMinusFee)

				if balanceUse == constant.NotSettledBalance {
					err = tr.merchantRepoWrites.UpdateMerchantCapitalAndNotSettleBalance(formattedNotSettledOrSettleBalanceMinusFee, formattedBalanceMinusFee, transactionData.MerchantId)
					if err != nil {
						resp = dto.ResponseDto{
							ResponseCode:    http.StatusUnprocessableEntity,
							ResponseMessage: err.Error(),
						}
						return resp, err
					}
				}

				if balanceUse == constant.SettleBalance {
					// update merchant balance
					err = tr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedNotSettledOrSettleBalanceMinusFee, formattedBalanceMinusFee, transactionData.MerchantId)
					if err != nil {
						resp = dto.ResponseDto{
							ResponseCode:    http.StatusUnprocessableEntity,
							ResponseMessage: err.Error(),
						}
						return resp, err
					}
				}

				// create merchant capital flow for transaction amount that already minus on capital balance
				payloadMerchantCapitalFlowPayin := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       formattedBalanceMinusFee,
					ReasonId:          constant.ReasonIdPayin,
					Status:            constant.StatusReversed,
					CreateBy:          constant.CreateBySystem,
					Amount:            helper.FormatFloat64(transactionAmountMinusFee),
					Notes:             "reversed balance due to manual change into failed",
					CapitalType:       constant.CapitalTypeDebit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowPayin)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}
			}
			// end of blocked code transaction have SUCCESS status before

			// start of blocked code transaction have PROCESSING status change into FAILED
			if transactionData.Status == constant.StatusProcessing {
				// create merchant capital flow for transaction amount
				payloadMerchantCapitalFlowPayin := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       helper.FormatFloat64(merchantAccountData.BalanceCapitalFlow),
					ReasonId:          constant.ReasonIdPayin,
					Status:            strings.ToUpper(status),
					CreateBy:          constant.CreateBySystem,
					Amount:            helper.FormatFloat64(transactionData.TransactionAmount),
					CapitalType:       constant.CapitalTypeNotDebitNotCredit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowPayin)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// calculate fee for merchant capital flow
				var feeTransactions float64
				if transactionData.MerchantFeeType == constant.FeeTypePercentage {
					calculateFee := math.Ceil(transactionData.TransactionAmount * (transactionData.MerchantFee / 100))
					feeTransactions = helper.FormatFloat64(calculateFee)
				} else {
					feeTransactions = helper.FormatFloat64(transactionData.MerchantFee)
				}

				// create merchant capital flow for fee
				payloadMerchantCapitalFlowFee := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       helper.FormatFloat64(merchantAccountData.BalanceCapitalFlow),
					ReasonId:          constant.ReasonIdFee,
					Status:            strings.ToUpper(status),
					CreateBy:          constant.CreateBySystem,
					Amount:            feeTransactions,
					CapitalType:       constant.CapitalTypeNotDebitNotCredit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowFee)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}
			}
			// end of blocked code transaction have PROCESSING status change into FAILED
		}

		// handle for payout and FAILED
		if transactionData.PayType == constant.PayTypePayout {
			// get merchant account
			merchantAccountData, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// check if transaction have SUCCESS status before
			if transactionData.Status == constant.StatusSuccess {
				balanceCapitalFlow := merchantAccountData.BalanceCapitalFlow
				settledBalance := merchantAccountData.SettledBalance

				// calculate fee for merchant capital flow
				feeTransactions := helper.FormatFloat64(transactionData.MerchantFee)

				balanceCapitalAddPayout := balanceCapitalFlow + transactionData.TransactionAmount
				settleBalanceAddPayout := settledBalance + transactionData.TransactionAmount
				formattedSettleBalanceAddPayout := helper.FormatFloat64(settleBalanceAddPayout)
				formattedBalanceAddPayout := helper.FormatFloat64(balanceCapitalAddPayout)

				// update merchant balance
				err = tr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalanceAddPayout, formattedBalanceAddPayout, transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// create merchant capital flow for transaction amount that already add on capital balance
				payloadMerchantCapitalFlowPayin := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       formattedBalanceAddPayout,
					ReasonId:          constant.ReasonIdPayout,
					Status:            constant.StatusReversed,
					CreateBy:          constant.CreateBySystem,
					Amount:            helper.FormatFloat64(transactionData.TransactionAmount),
					Notes:             "reversed balance due to manual change into failed",
					CapitalType:       constant.CapitalTypeCredit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowPayin)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// get merchant account after add payout
				merchantAccountDataUpdated, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				balanceCapitalFlow = merchantAccountDataUpdated.BalanceCapitalFlow
				settledBalance = merchantAccountDataUpdated.SettledBalance

				// adjust balance add fee
				balanceCapitalAddFee := balanceCapitalFlow + feeTransactions
				settledBalanceAddFee := settledBalance + feeTransactions
				formattedSettledBalanceAddFee := helper.FormatFloat64(settledBalanceAddFee)
				formattedBalanceAddFee := helper.FormatFloat64(balanceCapitalAddFee)

				err = tr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettledBalanceAddFee, formattedBalanceAddFee, transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// create merchant capital flow for fee
				payloadMerchantCapitalFlowFee := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       formattedBalanceAddFee,
					ReasonId:          constant.ReasonIdFee,
					Status:            constant.StatusReversed,
					CreateBy:          constant.CreateBySystem,
					Amount:            feeTransactions,
					Notes:             "reverse fee to balance due to manual failed",
					CapitalType:       constant.CapitalTypeCredit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowFee)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}
			}
			// end of blocked code transaction have SUCCESS status before

			// start of blocked code transaction have PROCESSING status change into FAILED
			if transactionData.Status == constant.StatusProcessing {
				feeTransactions := helper.FormatFloat64(transactionData.MerchantFee)

				// get merchant account
				merchantAccountData, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				transactionPlusFee := transactionData.TransactionAmount + feeTransactions
				pendingOutBalance := merchantAccountData.PendingTransactionOut

				// payout amount to balance
				pendingOutBalanceMinusFeeAndPayout := pendingOutBalance - transactionPlusFee
				formattedPendingOutBalanceMinusPayout := helper.FormatFloat64(pendingOutBalanceMinusFeeAndPayout)
				formattedBalanceCapital := helper.FormatFloat64(merchantAccountData.BalanceCapitalFlow)

				// update merchant balance
				err = tr.merchantRepoWrites.UpdateMerchantCapitalPendingOut(formattedPendingOutBalanceMinusPayout, formattedBalanceCapital, transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// get merchant account after minus payout
				merchantAccountDataUpdated, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(transactionData.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// adding settle balance from pending payout
				settleBalance := merchantAccountDataUpdated.SettledBalance
				settleBalanceAddFeeAndPayout := settleBalance + transactionPlusFee
				formattedSettleBalance := helper.FormatFloat64(settleBalanceAddFeeAndPayout)
				formattedBalanceCapitalUpdated := helper.FormatFloat64(merchantAccountDataUpdated.BalanceCapitalFlow)

				// update merchant balance for settle balanc
				err = tr.merchantRepoWrites.UpdateMerchantCapitalAndSettleBalance(formattedSettleBalance, formattedBalanceCapitalUpdated, merchantAccountDataUpdated.MerchantId)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// create merchant capital flow for transaction amount
				payloadMerchantCapitalFlowPayout := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       helper.FormatFloat64(merchantAccountDataUpdated.BalanceCapitalFlow),
					ReasonId:          constant.ReasonIdPayout,
					Status:            strings.ToUpper(status),
					CreateBy:          constant.CreateBySystem,
					Amount:            helper.FormatFloat64(transactionData.TransactionAmount),
					CapitalType:       constant.CapitalTypeNotDebitNotCredit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowPayout)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}

				// create merchant capital flow for fee
				payloadMerchantCapitalFlowFee := dto.CreateMerchantCapitalFlowPayload{
					PaymentId:         transactionData.PaymentID,
					MerchantAccountId: merchantAccountData.Id,
					TempBalance:       helper.FormatFloat64(merchantAccountDataUpdated.BalanceCapitalFlow),
					ReasonId:          constant.ReasonIdFee,
					Status:            strings.ToUpper(status),
					CreateBy:          constant.CreateBySystem,
					Amount:            feeTransactions,
					CapitalType:       constant.CapitalTypeNotDebitNotCredit,
				}
				_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowFee)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}
			}
		}
	}
	// end blocked code for handle FAILED status

	msg := fmt.Sprintf("success updated transaction status for payment id %v", transactionData.PaymentID)
	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: msg,
	}

	return resp, nil
}

func (tr *Transaction) GetListFilterSvc() (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listFilter, err := tr.merchantRepoReads.GetListFilter()
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success get list filter",
		Data:            listFilter,
	}

	return resp, nil
}

func (tr *Transaction) CreateMerchantExportSvc(payload dto.CreateMerchantExportReqDto) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	generateRandomNumber := helper.GenerateRandomString(5)
	fileName := fmt.Sprintf("report%v-%v-%v.xlsx", generateRandomNumber, payload.MerchantId, constant.TransformExportType[payload.ExportType])

	go func() {
		_, err := tr.supportExportTypeCapitalFlow(payload, fileName)
		if err != nil {
			slog.Infof("Error in supportExportTypeCapitalFlow: %v", err)

			if err.Error() == "data empty" {
				// update report list with status Error
				err = tr.transactionRepoWrites.UpdateReportStoragesByFileName("no url", fileName, constant.ReportStatusNoData)
				if err != nil {
					return
				}
				return
			}

			// update report list with status Error
			err = tr.transactionRepoWrites.UpdateReportStoragesByFileName("no url", fileName, constant.ReportStatusError)
			if err != nil {
				return
			}
			return
		}
	}()

	// payload list storage
	minExtract, _ := helper.ExtractDate(payload.MinDate)
	maxExtract, _ := helper.ExtractDate(payload.MaxDate)
	period := fmt.Sprintf("%v - %v", minExtract, maxExtract)
	reportPayload := dto.CreateReportStorageDto{
		MerchantId:    payload.MerchantId,
		Period:        period,
		ExportType:    payload.ExportType,
		Status:        constant.ReportStatusPending,
		ReportUrl:     "",
		CreatedByUser: payload.UserType,
		FileName:      fileName,
	}

	// create list report storages
	id, err := tr.transactionRepoWrites.CreateListReportStoragesRepo(reportPayload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success",
		Data:            fmt.Sprintf("Success create list with id: %v", id),
	}

	return resp, nil
}

func (tr *Transaction) GetListMerchantExportSvc(params dto.GetListMerchantExportFilter) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if params.MinDate == "" {
		params.MinDate = helper.GenerateTime(0)
	}

	if params.MaxDate == "" {
		params.MaxDate = helper.GenerateTime(24)
	}

	listMerchantExport, err := tr.transactionRepoReads.GetListMerchantExportRepo(params)
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
		Data:            listMerchantExport,
	}

	return resp, nil
}

func (tr *Transaction) GetListInternalExportSvc(params dto.GetListMerchantExportFilter) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var listInternalExport []dto.InternalExportDto

	if params.MinDate == "" {
		params.MinDate = helper.GenerateTime(0)
	}

	if params.MaxDate == "" {
		params.MaxDate = helper.GenerateTime(24)
	}

	params.Merchants = constant.InternalExport
	listMerchantExport, err := tr.transactionRepoReads.GetListMerchantExportRepo(params)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	for i := range listMerchantExport {
		listReport := dto.InternalExportDto{
			Id:         listMerchantExport[i].Id,
			CreatedAt:  listMerchantExport[i].CreatedAt,
			Currency:   "IDR",
			ExportType: listMerchantExport[i].ExportType,
			Period:     listMerchantExport[i].Period,
			Status:     listMerchantExport[i].Status,
			ReportUrl:  listMerchantExport[i].ReportUrl,
		}
		listInternalExport = append(listInternalExport, listReport)
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            listInternalExport,
	}

	return resp, nil
}

func (tr *Transaction) GetListFilterExportSvc() (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var listExportType []dto.ExportData
	var listExportStatus []dto.ExportStatusData

	for i := range constant.ExportStatus {
		exportStatus := dto.ExportStatusData{
			Id:               i + 1,
			ExportStatusName: constant.ExportStatus[i],
		}
		listExportStatus = append(listExportStatus, exportStatus)
	}

	for i := range constant.ExportType {
		exportType := dto.ExportData{
			Id:             i + 1,
			ExportTypeName: constant.ExportType[i],
		}
		listExportType = append(listExportType, exportType)
	}

	listFilterDtoResp := dto.ListFilterExportDto{
		ExportType:   listExportType,
		ExportStatus: listExportStatus,
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve filter export",
		Data:            listFilterDtoResp,
	}

	return resp, nil
}

func (tr *Transaction) CreateInternalExportSvc(payload dto.CreateMerchantExportReqDto) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	generateRandomNumber := helper.GenerateRandomString(5)
	fileName := fmt.Sprintf("internal-report%v-%v.xlsx", generateRandomNumber, constant.TransformExportType[payload.ExportType])

	go func() {
		_, err := tr.supportExportTypeCapitalFlow(payload, fileName)
		if err != nil {
			slog.Infof("Error in supportExportTypeCapitalFlow: %v", err)

			if err.Error() == "data empty" {
				// update report list with status Error
				err = tr.transactionRepoWrites.UpdateReportStoragesByFileName("no url", fileName, constant.ReportStatusNoData)
				if err != nil {
					return
				}
				return
			}

			// update report list with status Error
			err = tr.transactionRepoWrites.UpdateReportStoragesByFileName("no url", fileName, constant.ReportStatusError)
			if err != nil {
				return
			}
			return
		}
	}()

	// payload list storage
	minExtract, _ := helper.ExtractDate(payload.MinDate)
	maxExtract, _ := helper.ExtractDate(payload.MaxDate)
	period := fmt.Sprintf("%v - %v", minExtract, maxExtract)
	reportPayload := dto.CreateReportStorageDto{
		MerchantId:    "INTERNAL-EXPORT",
		Period:        period,
		ExportType:    payload.ExportType,
		Status:        constant.ReportStatusPending,
		ReportUrl:     "",
		CreatedByUser: payload.UserType,
		FileName:      fileName,
	}

	// create list report storages
	id, err := tr.transactionRepoWrites.CreateListReportStoragesRepo(reportPayload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success",
		Data:            fmt.Sprintf("Success create internal export with id: %v", id),
	}

	return resp, nil
}

func (tr *Transaction) GetTransactionDetailSvc(paymentId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var transactionData dto.TransactionMerchantData
	var failedInformation dto.FailureInformationDto
	var CallbackData dto.CallbackMerchantResp

	paymentData, err := tr.transactionRepoReads.GetPaymentDetailProviderMerchant(paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// transaction detail data
	if paymentData.PaymentMethodName == constant.VirtualAccountPaymentMethod {
		transactionNetAmount := paymentData.TransactionAmount - paymentData.MerchantFee
		accountData, _ := tr.transactionRepoReads.GetAccountInformationByPaymentIdAccountType(paymentId, constant.AccountTypeCreditor)

		transactionData = dto.TransactionMerchantData{
			Name:                  paymentData.MerchantName,
			MerchantTransactionId: paymentData.MerchantRefNumber,
			TransactionFee:        paymentData.MerchantFee,
			TransactionFeeType:    paymentData.MerchantFeeType,
			TransactionAmount:     paymentData.TransactionAmount,
			TransactionNetAmount:  transactionNetAmount,
			TransactionMethod:     paymentData.PaymentMethodName,
			AccountName:           accountData.AccountName,
			AccountNumber:         accountData.AccountNumber,
			BankName:              accountData.BankName,
			IpAddress:             paymentData.ClientIPAddress,
		}
	}

	if paymentData.PaymentMethodName == constant.QrisPaymentMethod || paymentData.PaymentMethodName == constant.EwalletPaymentMethod {
		transactionFee := math.Ceil(paymentData.TransactionAmount * (paymentData.MerchantFee / 100))
		transactionFeeFormatted := helper.FormatFloat64(transactionFee)
		transactionNetAmount := paymentData.TransactionAmount - transactionFeeFormatted

		transactionData = dto.TransactionMerchantData{
			Name:                  paymentData.MerchantName,
			MerchantTransactionId: paymentData.MerchantRefNumber,
			TransactionFee:        paymentData.MerchantFee,
			TransactionFeeType:    paymentData.MerchantFeeType,
			TransactionAmount:     paymentData.TransactionAmount,
			TransactionNetAmount:  transactionNetAmount,
			TransactionMethod:     paymentData.PaymentMethodName,
			AccountName:           nil,
			AccountNumber:         nil,
			BankName:              nil,
			IpAddress:             paymentData.ClientIPAddress,
		}
	}

	if paymentData.PaymentMethodName == constant.DisbursementPaymentMethod {
		transactionNetAmount := paymentData.TransactionAmount + paymentData.MerchantFee
		accountData, _ := tr.transactionRepoReads.GetAccountInformationByPaymentIdAccountType(paymentId, constant.AccountTypeCreditor)

		transactionData = dto.TransactionMerchantData{
			Name:                  paymentData.MerchantName,
			MerchantTransactionId: paymentData.MerchantRefNumber,
			TransactionFee:        paymentData.MerchantFee,
			TransactionFeeType:    paymentData.MerchantFeeType,
			TransactionAmount:     paymentData.TransactionAmount,
			TransactionNetAmount:  transactionNetAmount,
			TransactionMethod:     paymentData.PaymentMethodName,
			AccountName:           accountData.AccountName,
			AccountNumber:         accountData.AccountNumber,
			BankName:              accountData.BankName,
			IpAddress:             paymentData.ClientIPAddress,
		}
	}

	listCallback, err := tr.merchantRepoReads.GetListMerchantCallback(paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listCallback) > 0 {
		latestCallback := listCallback[0]
		latestCallback.RetriedAt = listCallback[0].CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00")
		latestCallback.StartedAt = listCallback[len(listCallback)-1].CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00")

		// callback data
		CallbackData = dto.CallbackMerchantResp{
			StatusCallback:              latestCallback.CallbackStatus,
			TransactionStatusInCallback: latestCallback.PaymentStatusInCallback,
			BeginsAt:                    latestCallback.StartedAt,
			LatestAt:                    latestCallback.RetriedAt,
			MerchantResponse:            latestCallback.CallbackResult,
		}
	}

	statusChangeLogsData, err := tr.transactionRepoReads.GetStatusChangeLogData(paymentId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if statusChangeLogsData[0].StatusLog == constant.StatusLogFailed {
		failedAtStr := statusChangeLogsData[0].CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00")
		failedInformation = dto.FailureInformationDto{
			FailedAt: failedAtStr,
			Message:  *statusChangeLogsData[0].Notes,
		}
	}

	detailResp := dto.TransactionMerchantDetailDto{
		TransactionData:    transactionData,
		MerchantCallback:   CallbackData,
		FailureInformation: failedInformation,
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            detailResp,
	}

	return resp, nil
}

func (tr *Transaction) MerchantDisbursementSvc(payload dto.MerchantDisbursement) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var disburseMerchantChannel entity.MerchantPaychannel

	// business validation
	user, err := tr.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		slog.Infof("username: %v, failed get user data, err: %v", payload.Username, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
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

	merchantData, err := tr.merchantRepoReads.GetMerchantDataByMerchantId(*user.MerchantID)
	if err != nil {
		slog.Infof("username: %v, failed get merchant data, err: %v", payload.Username, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	if merchantData.Status == constant.StatusInactive {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant status is inactive",
		}
		return resp, errors.New("insufficient")
	}

	accountBalance, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(*user.MerchantID)
	if err != nil {
		slog.Infof("username: %v, failed get account balance, err: %v", payload.Username, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	if float64(payload.Amount) > accountBalance.SettledBalance {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "not enough balance for disbursement",
		}
		return resp, errors.New("insufficient")
	}

	listMerchantPaychannel, err := tr.merchantRepoReads.GetMerchantPaychannelByMerchantId(*user.MerchantID)
	if err != nil {
		slog.Infof("username: %v, failed get merchant channel, err: %v", payload.Username, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	if len(listMerchantPaychannel) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "this merchant not routed for disbursement",
		}
		return resp, errors.New("insufficient")
	}

	for _, channel := range listMerchantPaychannel {
		if channel.PaymentMethodChannel == constant.DisbursementPaymentMethod && channel.Segment == constant.MainType {
			disburseMerchantChannel = channel
		}
	}

	if disburseMerchantChannel.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "this merchant not routed for disbursement",
		}
		return resp, errors.New("insufficient")
	}

	getRoutedChannel, err := tr.merchantRepoReads.GetListRoutedPaychannelByIdMerchantPaychannelRepo(disburseMerchantChannel.Id)
	if err != nil {
		slog.Infof("username: %v, getRoutedChannel got err: %v", payload.Username, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	if len(getRoutedChannel) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "this merchant not routed for disbursement",
		}
		return resp, errors.New("insufficient")
	}

	if getRoutedChannel[0].Status == constant.StatusInactive {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "merchant status is inactive",
		}
		return resp, errors.New("insufficient")
	}

	if disburseMerchantChannel.MinTransaction > 0 || disburseMerchantChannel.MaxTransaction > 0 {
		if float64(payload.Amount) < disburseMerchantChannel.MinTransaction || float64(payload.Amount) > disburseMerchantChannel.MaxTransaction {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusBadRequest,
				ResponseMessage: "amount limit",
			}
			return resp, errors.New("insufficient")
		}
	}

	if getRoutedChannel[0].MinTransaction > 0 || getRoutedChannel[0].MaxTransaction > 0 {
		if float64(payload.Amount) < getRoutedChannel[0].MinTransaction || float64(payload.Amount) > getRoutedChannel[0].MaxTransaction {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusBadRequest,
				ResponseMessage: "amount limit",
			}
			return resp, errors.New("insufficient")
		}
	}

	providerId := getRoutedChannel[0].ProviderId
	interfaceSetting := getRoutedChannel[0].InterfaceSetting
	credentials, err := tr.providerRepoReads.GetAllCredentialsRepo(providerId, interfaceSetting)
	if err != nil {
		slog.Infof("username: %v, get credentials got err: %v", payload.Username, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	bankData, err := tr.transactionRepoReads.GetBankDataDetailRepo(payload.BankName)
	if err != nil {
		slog.Infof("username: %v, got err: %v", payload.Username, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	channelIdCodePayload := dto.ChannelIdCodeDisbursement{
		MerchantPaychanneId:  disburseMerchantChannel.Id,
		ProviderPaychannelId: getRoutedChannel[0].ProviderPaychannelId,
		BankCode:             bankData.BankCode,
	}

	_, err = tr.disbursementSupport(credentials, payload, *user.MerchantID, disburseMerchantChannel.Fee, channelIdCodePayload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: constant.GeneralErrMsg,
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: fmt.Sprintf("success dibursement with bank name: %v, account number: %v, account name: %v", payload.BankName, payload.BankAccountNumber, payload.BankAccountName),
	}

	return resp, nil
}

func (tr *Transaction) GetBankListDisbursementSvc(username string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var disburseMerchantChannel entity.MerchantPaychannel
	var paychannel string

	user, err := tr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	listMerchantPaychannel, err := tr.merchantRepoReads.GetMerchantPaychannelByMerchantId(*user.MerchantID)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listMerchantPaychannel) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "bank list not found, this merchant not routed for disbursement",
		}
		return resp, nil
	}

	for _, channel := range listMerchantPaychannel {
		if channel.PaymentMethodChannel == constant.DisbursementPaymentMethod && channel.Segment == constant.MainType {
			disburseMerchantChannel = channel
		}
	}

	if disburseMerchantChannel.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "bank list not found, this merchant not routed for disbursement",
		}
		return resp, nil
	}

	getRoutedChannel, err := tr.merchantRepoReads.GetListRoutedPaychannelByIdMerchantPaychannelRepo(disburseMerchantChannel.Id)
	if err != nil {
		slog.Infof("getRoutedChannel got err: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "please ask admin for more info",
		}
		return resp, err
	}

	if len(getRoutedChannel) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "bank list not found, this merchant not routed for disbursement",
		}
		return resp, nil
	}

	if len(getRoutedChannel) > 1 {
		for _, routedChan := range getRoutedChannel {
			paychannel += routedChan.ProviderPaychannelName + ","
		}
	}

	if len(getRoutedChannel) == 1 {
		paychannel += getRoutedChannel[0].ProviderPaychannelName
	}

	routedChannelName := fmt.Sprintf("[%v]", paychannel)
	getBankList, err := tr.merchantRepoReads.GetBankListForDisbursementRepo(routedChannelName)
	if err != nil {
		slog.Infof("getBankList got err: %v", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "please ask admin for more info",
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            getBankList,
	}

	return resp, nil
}

func (tr *Transaction) CountDisbursementTotalAmountSvc(payload dto.CountDisbursementTotalAmountDto) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	var disburseMerchantChannel entity.MerchantPaychannel

	user, err := tr.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	listMerchantPaychannel, err := tr.merchantRepoReads.GetMerchantPaychannelByMerchantId(*user.MerchantID)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listMerchantPaychannel) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "this merchant not routed for disbursement",
		}
		return resp, nil
	}

	for _, channel := range listMerchantPaychannel {
		if channel.PaymentMethodChannel == constant.DisbursementPaymentMethod && channel.Segment == constant.MainType {
			disburseMerchantChannel = channel
		}
	}

	if disburseMerchantChannel.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "this merchant not routed for disbursement",
		}
		return resp, nil
	}

	countTotalAmount := float64(payload.Amount) + disburseMerchantChannel.Fee

	countResp := dto.CountDisbursementRespDto{
		FeeAmount:   disburseMerchantChannel.Fee,
		TotalAmount: countTotalAmount,
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success counted",
		Data:            countResp,
	}

	return resp, nil
}

func (tr *Transaction) JackDisbursementCallbackHandlingSvc(payload dto.CreateDisbursementRequestResponseData) (string, error) {
	amountInt := converter.FromStringToIntAmount(payload.Destination.Amount)

	slog.Infof("Jack %v callback payload: %v", payload.ReferenceID, payload)

	if payload.State == constant.JackStateStatusDeclined || payload.State == constant.JackStateStatusCanceled {
		// update status into failed
		err := tr.transactionRepoWrites.UpdateStatus(constant.StatusFailed, payload.ReferenceID)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// update merchant balance
		detailTransaction, err := tr.transactionRepoReads.GetPaymentDetailProviderMerchant(payload.ReferenceID)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		merchantBalance, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(detailTransaction.MerchantId)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		settleBalance := merchantBalance.SettledBalance
		pendingPayout := merchantBalance.PendingTransactionOut

		settleBalancePlusOut := settleBalance + float64(amountInt) + detailTransaction.MerchantFee
		pendingPayoutMinusOut := pendingPayout - float64(amountInt) - detailTransaction.MerchantFee
		err = tr.merchantRepoWrites.UpdateMerchantBalanceSettleAndPendingOutBalanceRepo(settleBalancePlusOut, pendingPayoutMinusOut, detailTransaction.MerchantId)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// create merchant capital flow for transaction amount
		payloadMerchantCapitalFlowPayout := dto.CreateMerchantCapitalFlowPayload{
			PaymentId:         payload.ReferenceID,
			MerchantAccountId: merchantBalance.Id,
			TempBalance:       helper.FormatFloat64(merchantBalance.BalanceCapitalFlow),
			ReasonId:          constant.ReasonIdPayout,
			Status:            strings.ToUpper(constant.StatusFailed),
			CreateBy:          constant.CreateBySystem,
			Amount:            helper.FormatFloat64(float64(amountInt)),
			CapitalType:       constant.CapitalTypeNotDebitNotCredit,
		}
		_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowPayout)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// create merchant capital flow for fee
		payloadMerchantCapitalFlowFee := dto.CreateMerchantCapitalFlowPayload{
			PaymentId:         payload.ReferenceID,
			MerchantAccountId: merchantBalance.Id,
			TempBalance:       helper.FormatFloat64(merchantBalance.BalanceCapitalFlow),
			ReasonId:          constant.ReasonIdFee,
			Status:            strings.ToUpper(constant.StatusFailed),
			CreateBy:          constant.CreateBySystem,
			Amount:            detailTransaction.MerchantFee,
			CapitalType:       constant.CapitalTypeNotDebitNotCredit,
		}
		_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowFee)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// update status transaction log
		_, err = tr.transactionRepoWrites.CreateTransactionStatusLog(payload.ReferenceID, constant.StatusLogFailed, constant.CreateBySystem, constant.GeneralErrMsg, payload.ErrorMessage)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// create provider confirmation detail
		_, err = tr.providerRepoWrites.CreateProviderConfirmationDetail(constant.SourceCallback, payload.ReferenceID, constant.StatusFailed)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}
	}

	if payload.State == constant.JackStateStatusCompleted {
		// update status into success
		err := tr.transactionRepoWrites.UpdateStatus(constant.StatusSuccess, payload.ReferenceID)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// update merchant balance
		detailTransaction, err := tr.transactionRepoReads.GetPaymentDetailProviderMerchant(payload.ReferenceID)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		merchantBalance, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(detailTransaction.MerchantId)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		pendingPayout := merchantBalance.PendingTransactionOut
		BalanceCapital := merchantBalance.BalanceCapitalFlow

		pendingPayoutMinusOut := pendingPayout - float64(amountInt)
		balanceCapitalMinusOut := BalanceCapital - float64(amountInt)
		err = tr.merchantRepoWrites.UpdateMerchantCapitalPendingOut(pendingPayoutMinusOut, balanceCapitalMinusOut, detailTransaction.MerchantId)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// create merchant capital flow out
		payloadMerchantCapitalFlowOut := dto.CreateMerchantCapitalFlowPayload{
			PaymentId:         payload.ReferenceID,
			MerchantAccountId: merchantBalance.Id,
			TempBalance:       balanceCapitalMinusOut,
			ReasonId:          constant.ReasonIdPayout,
			Status:            strings.ToUpper(constant.StatusSuccess),
			CreateBy:          constant.CreateBySystem,
			Amount:            helper.FormatFloat64(float64(amountInt)),
			CapitalType:       constant.CapitalTypeDebit,
		}
		_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowOut)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		merchantBalanceUpdate, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(detailTransaction.MerchantId)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		pendingPayoutUpdated := merchantBalanceUpdate.PendingTransactionOut
		balanceCapitalUpdate := merchantBalanceUpdate.BalanceCapitalFlow

		pendingPayoutUpdatedMinusFee := pendingPayoutUpdated - detailTransaction.MerchantFee
		balanceCapitalUpdateMinusFee := balanceCapitalUpdate - detailTransaction.MerchantFee
		err = tr.merchantRepoWrites.UpdateMerchantCapitalPendingOut(pendingPayoutUpdatedMinusFee, balanceCapitalUpdateMinusFee, detailTransaction.MerchantId)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// create merchant capital fee
		payloadMerchantCapitalFlowFee := dto.CreateMerchantCapitalFlowPayload{
			PaymentId:         payload.ReferenceID,
			MerchantAccountId: merchantBalance.Id,
			TempBalance:       balanceCapitalUpdateMinusFee,
			ReasonId:          constant.ReasonIdFee,
			Status:            strings.ToUpper(constant.StatusSuccess),
			CreateBy:          constant.CreateBySystem,
			Amount:            detailTransaction.MerchantFee,
			CapitalType:       constant.CapitalTypeDebit,
		}
		_, err = tr.merchantRepoWrites.CreateMerchantCapitalFlow(payloadMerchantCapitalFlowFee)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// update transaction status log
		_, err = tr.transactionRepoWrites.CreateTransactionStatusLog(payload.ReferenceID, constant.StatusLogSuccess, constant.CreateBySystem, "", "")
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}

		// create provider confirmation detail
		_, err = tr.providerRepoWrites.CreateProviderConfirmationDetail(constant.SourceCallback, payload.ReferenceID, constant.StatusSuccess)
		if err != nil {
			slog.Infof("JackDisbursementHandlingSvc %v got err: %v", payload.ReferenceID, err.Error())
			return "", err
		}
	}

	return "ok", nil
}

func (tr *Transaction) GetReportListMerchantSvc(req dto.GetListMerchantExportFilter, username string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	user, err := tr.userRepoReads.GetUserByUsername(username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	listMerchantExport, err := tr.transactionRepoReads.GetListMerchantReportRepo(req, *user.MerchantID)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            listMerchantExport,
	}

	return resp, nil
}

func (tr *Transaction) CreateReportMerchantSvc(req dto.CreateReportMerchantReqDto) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	generateRandomNumber := helper.GenerateRandomString(5)
	fileName := fmt.Sprintf("merchant%v-%v.xlsx", generateRandomNumber, constant.TransformExportType[req.ExportType])

	user, err := tr.userRepoReads.GetUserByUsername(req.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	req.MerchantId = *user.MerchantID
	payloadReport := dto.CreateMerchantExportReqDto{
		ExportType: req.ExportType,
		MinDate:    req.MinDate,
		MaxDate:    req.MaxDate,
		MerchantId: req.MerchantId,
		UserType:   req.UserType,
	}

	go func() {
		_, err := tr.supportExportTypeCapitalFlow(payloadReport, fileName)
		if err != nil {
			slog.Infof("Error in supportExportTypeCapitalFlow: %v", err)

			if err.Error() == "data empty" {
				// update report list with status Error
				err = tr.transactionRepoWrites.UpdateReportStoragesByFileName("no url", fileName, constant.ReportStatusNoData)
				if err != nil {
					return
				}
				return
			}

			// update report list with status Error
			err = tr.transactionRepoWrites.UpdateReportStoragesByFileName("no url", fileName, constant.ReportStatusError)
			if err != nil {
				return
			}
			return
		}
	}()

	// payload list storage
	minExtract, _ := helper.ExtractDate(req.MinDate)
	maxExtract, _ := helper.ExtractDate(req.MaxDate)
	period := fmt.Sprintf("%v - %v", minExtract, maxExtract)
	reportPayload := dto.CreateReportStorageDto{
		MerchantId:    *user.MerchantID,
		Period:        period,
		ExportType:    req.ExportType,
		Status:        constant.ReportStatusPending,
		ReportUrl:     "",
		CreatedByUser: req.UserType,
		FileName:      fileName,
	}

	// create list report storages
	id, err := tr.transactionRepoWrites.CreateListReportStoragesRepo(reportPayload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success",
		Data:            fmt.Sprintf("Success create report with id: %v", id),
	}

	return resp, nil
}

func (tr *Transaction) GetListTransactionMerchantFlowSvc(params dto.QueryParams) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if params.MinDate == "" {
		params.MinDate = helper.GenerateTime(0)
	}

	if params.MaxDate == "" {
		params.MaxDate = helper.GenerateTime(24)
	}

	user, err := tr.userRepoReads.GetUserByUsername(params.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	params.MerchantId = *user.MerchantID
	listTransactionCapital, pagination, err := tr.transactionRepoReads.GetTransactionCapitalFlowRepo(params)
	if err != nil {
		slog.Infof("merchant id %v list transaction capital got failed: %v", params.MerchantId, err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "Data not found",
		}
		return resp, err
	}

	if len(listTransactionCapital) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "Data not found",
		}
		return resp, nil
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve list transaction capital flow",
		Data:            listTransactionCapital,
		Pagination:      pagination,
	}

	return resp, nil
}

func (tr *Transaction) supportExportTypeCapitalFlow(payload dto.CreateMerchantExportReqDto, fileName string) (string, error) {
	var headers []string

	// get data for create excel
	list, err := tr.transactionRepoWrites.CreateMerchantExportCapitalFlowRepo(payload)
	if err != nil {
		slog.Infof("supportExportTypeCapitalFlow got err: %v", err.Error())
		return "", err
	}

	if len(list) == 0 {
		return "", errors.New("data empty")
	}

	data := make([][]interface{}, len(list))

	if payload.UserType == constant.UserMerchant {
		headers = []string{
			"Payment ID",
			"Merchant ID",
			"Merchant Name",
			"Amount",
			"Reason Name",
			"Payment Method",
			"Merchant Fee",
			"Merchant Fee Type",
			"Merchant Balance",
			"Status",
			"Notes",
			"Capital Type",
			"Reverse From",
			"Created At",
		}

		// Prepare data for Excel
		for i, item := range list {
			data[i] = []interface{}{
				item.PaymentId,
				nullSafeString(item.MerchantId),
				nullSafeString(item.MerchantName),
				nullSafeFloat64(item.Amount),
				nullSafeString(item.ReasonName),
				nullSafeString(item.PaymentMethod),
				nullSafeFloat64(item.Fee),
				nullSafeString(item.FeeType),
				nullSafeFloat64(item.MerchantBalance),
				nullSafeString(item.Status),
				nullSafeString(item.Notes),
				nullSafeString(item.CapitalType),
				nullSafeString(item.ReverseFrom),
				item.CreatedAt,
			}
		}
	}

	if payload.UserType == constant.UserOperation {

		headers = []string{
			"Payment ID",
			"Merchant ID",
			"Merchant Name",
			"Amount",
			"Reason Name",
			"Payment Method",
			"Merchant Fee",
			"Merchant Fee Type",
			"Merchant Balance",
			"Provider",
			"Paychannel Routed",
			"Provider Fee",
			"Provider Fee Type",
			"Status",
			"Notes",
			"Capital Type",
			"Reverse From",
			"Created At",
		}

		// Prepare data for Excel
		for i, item := range list {
			data[i] = []interface{}{
				item.PaymentId,
				nullSafeString(item.MerchantId),
				nullSafeString(item.MerchantName),
				nullSafeFloat64(item.Amount),
				nullSafeString(item.ReasonName),
				nullSafeString(item.PaymentMethod),
				nullSafeFloat64(item.Fee),
				nullSafeString(item.FeeType),
				nullSafeFloat64(item.MerchantBalance),
				nullSafeString(item.Provider),
				nullSafeString(item.PaychannelRouted),
				nullSafeFloat64(item.ProviderFee),
				nullSafeString(item.ProviderFeeType),
				nullSafeString(item.Status),
				nullSafeString(item.Notes),
				nullSafeString(item.CapitalType),
				nullSafeString(item.ReverseFrom),
				item.CreatedAt,
			}
		}
	}

	err = helper.CreateExcelFile(headers, fileName, data)
	if err != nil {
		slog.Infof("%v", err.Error())
		return "", err
	}

	credentials := helper.GetSecret(tr.configApp)

	publicUrl, err := helper.UploadFile(constant.BucketName, fileName, fileName, credentials)
	if err != nil {
		return "", err
	}

	err = tr.transactionRepoWrites.UpdateReportStoragesByFileName(publicUrl, fileName, constant.ReportStatusFinished)
	if err != nil {
		return "", err
	}

	os.Remove(fileName)

	return "ok", nil
}

func (tr *Transaction) disbursementSupport(credentials []entity.ProviderCredentialsEntity, payload dto.MerchantDisbursement, merchantId string, merchantFee float64, channelCodeId dto.ChannelIdCodeDisbursement) (string, error) {
	providerId := credentials[0].ProviderId

	if providerId == constant.ProviderJack {
		_, err := tr.jackSupportDisbursement(credentials, payload, merchantId, merchantFee, channelCodeId)
		if err != nil {
			slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
			return "", err
		}
	}

	if !helper.StringInSlice(providerId, constant.ProviderListName) {
		return "", errors.New("this merchant not routed for disbursement")
	}

	return "ok", nil
}

func (tr *Transaction) jackSupportDisbursement(credentials []entity.ProviderCredentialsEntity, payload dto.MerchantDisbursement, merchantId string, merchantFee float64, channelCodeId dto.ChannelIdCodeDisbursement) (string, error) {
	randomStr := helper.GenerateRandomString(30)
	randomStrMerchantReferenceNumber := helper.GenerateRandomString(30)
	paymentId := "out_dsb-" + randomStr
	merchantReferenceNumber := merchantId + "-" + randomStrMerchantReferenceNumber
	currentBalance, err := tr.jackProvider.GetBalance(payload.Username, credentials)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	if payload.Amount > currentBalance {
		slog.Infof("username: %v, disbursement limit", payload.Username)
		return "", errors.New("amount limit")
	}

	inquiryData, err := tr.jackProvider.InquiryAccount(payload, credentials, channelCodeId.BankCode)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	accountHolder := inquiryData.Data.AccountName
	responseName := tr.regex.ReplaceAllString(accountHolder, "")
	requestName := tr.regex.ReplaceAllString(payload.BankAccountName, "")
	similarWord := helper.CompareTwoStrings(strings.ToLower(responseName), strings.ToLower(requestName))
	if similarWord < constant.BankAccountNameSimilarityMatchInPercent {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, "name validation not match")
		return "", errors.New("name validation not match")
	}

	createDisbursement, err := tr.jackProvider.CreateDisbursement(payload, credentials, channelCodeId.BankCode, paymentId)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	providerCreateId := converter.ToString(createDisbursement.Data.ID)
	confirmTransactionData := dto.ConfirmTransactionPayload{
		Username:   payload.Username,
		ProviderID: providerCreateId,
	}

	confirm, err := tr.jackProvider.ConfirmDisbursement(confirmTransactionData, credentials)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	// update merchant account
	merchantAccount, err := tr.merchantRepoReads.GetMerchantAccountByMerchantId(merchantId)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	// balance adjustment with fee
	settleBalance := merchantAccount.SettledBalance
	pendingOutBalance := merchantAccount.PendingTransactionOut
	settleBalanceMinusOutAndFee := settleBalance - float64(payload.Amount) - merchantFee
	pendingOutBalancePlusOutAndFee := pendingOutBalance + float64(payload.Amount) + merchantFee

	// updated merchant settle balance
	err = tr.merchantRepoWrites.UpdateMerchantBalanceSettleAndPendingOutBalanceRepo(settleBalanceMinusOutAndFee, pendingOutBalancePlusOutAndFee, merchantId)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	// create transaction
	createTransactionPayload := dto.CreateTransactionsDto{
		PaymentId:               paymentId,
		MerchantReferenceNumber: merchantReferenceNumber,
		ProviderReferenceNumber: converter.ToString(confirm.Data.ID),
		MerchantPaychanneId:     channelCodeId.MerchantPaychanneId,
		ProviderPaychannelId:    channelCodeId.ProviderPaychannelId,
		TransactionAmount:       float64(payload.Amount),
		BankCode:                channelCodeId.BankCode,
		Status:                  constant.StatusProcessing,
		RequestMethod:           "MERCHANT_DASHBOARD",
		IpAddress:               constant.IpAddressHypay,
		CallbackUrl:             constant.CallbackUrlHypay,
	}

	_, err = tr.transactionRepoWrites.CreateTransactionsRepo(createTransactionPayload)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	// create transaction logs
	_, err = tr.transactionRepoWrites.CreateTransactionStatusLog(paymentId, constant.StatusLogAcceptedByPlatform, constant.CreateBySystem, "", "")
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	_, err = tr.transactionRepoWrites.CreateTransactionStatusLog(paymentId, constant.StatusLogAcceptedByProvider, constant.CreateBySystem, "", "")
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	// create account information
	payloadAccountInformationCreditor := dto.CreateAccountInformationDto{
		PaymentId:       paymentId,
		AccountNumber:   payload.BankAccountNumber,
		AccountName:     payload.BankAccountName,
		BankName:        payload.BankName,
		BankCode:        channelCodeId.BankCode,
		ReferenceNumber: converter.ToString(confirm.Data.ID),
		AccountType:     constant.AccountTypeCreditor,
	}
	_, err = tr.transactionRepoWrites.CreateAccountInformationRepo(payloadAccountInformationCreditor)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	payloadAccountInformationDebitor := dto.CreateAccountInformationDto{
		PaymentId:   paymentId,
		AccountType: constant.AccountTypeDebitor,
	}
	_, err = tr.transactionRepoWrites.CreateAccountInformationRepo(payloadAccountInformationDebitor)
	if err != nil {
		slog.Infof("username: %v, disbursementSupport got err %v", payload.Username, err.Error())
		return "", err
	}

	return "ok", nil
}

func nullSafeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func nullSafeFloat64(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}
