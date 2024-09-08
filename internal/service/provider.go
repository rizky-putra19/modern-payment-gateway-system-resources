package service

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal"
	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
)

type Provider struct {
	transactionRepoReads internal.TransactionsReadsRepositoryItf
	merchantRepoReads    internal.MerchantReadsRepositoryItf
	providerRepoReads    internal.ProviderReadsRepositoryItf
	providerRepoWrites   internal.ProviderWritesRepositoryItf
	userRepoReads        internal.UserReadsRepositoryItf
	config               config.App
}

func NewProvider(
	transactionRepoReads internal.TransactionsReadsRepositoryItf,
	merchantRepoReads internal.MerchantReadsRepositoryItf,
	providerRepoReads internal.ProviderReadsRepositoryItf,
	providerRepoWrites internal.ProviderWritesRepositoryItf,
	userRepoReads internal.UserReadsRepositoryItf,
	config config.App,
) *Provider {
	return &Provider{
		transactionRepoReads: transactionRepoReads,
		merchantRepoReads:    merchantRepoReads,
		providerRepoReads:    providerRepoReads,
		providerRepoWrites:   providerRepoWrites,
		userRepoReads:        userRepoReads,
		config:               config,
	}
}

func (pr *Provider) GetListProvidersSvc(paymentMethod string, search string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listProvider, err := pr.providerRepoReads.GetListProvidersWithFilterRepo(paymentMethod, search)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	for i := range listProvider {
		activePaychannels, paymentMethodProviders, err := pr.providerRepoReads.CountInterfacesProviderPaychannelByIdProvider(listProvider[i].Id)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}

		interfaceStr := converter.ToString(activePaychannels) + "/" + converter.ToString(paymentMethodProviders)
		listProvider[i].Interfaces = interfaceStr
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrive data",
		Data:            listProvider,
	}

	return resp, nil
}

func (pr *Provider) GetProviderAnalyticsSvc(payload dto.GetProviderAnalyticsDtoReq) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if payload.MinDate == "" {
		payload.MinDate = helper.GenerateTime(0)
	}

	if payload.MaxDate == "" {
		payload.MaxDate = helper.GenerateTime(24)
	}

	transactionData, err := pr.transactionRepoReads.GetTransactionDataForProviderAnalyticsRepo(payload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	data := supportProviderAnalyticsSvc(transactionData)
	listProviderInterfaces, err := pr.providerRepoReads.GetProviderInterfacesRepoById(payload.ProviderId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	data.ProviderInterfaces = listProviderInterfaces

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success retrieve data",
		Data:            data,
	}

	return resp, nil
}

func (pr *Provider) GetListProviderPaychannelSvc(providerInterfaceId int) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listPaychannel, err := pr.providerRepoReads.GetListProviderPaychannelById(providerInterfaceId)
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
		Data:            listPaychannel,
	}

	return resp, nil
}

func (pr *Provider) GetListProviderChannelAllSvc(params dto.QueryParams) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listPaychannel, err := pr.providerRepoReads.GetListProviderChannelAllRepo(params)
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
		Data:            listPaychannel,
	}

	return resp, nil
}

func (pr *Provider) GetProviderChannelAnalyticsSvc(payload dto.GetProviderAnalyticsDtoReq) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	if payload.MinDate == "" {
		payload.MinDate = helper.GenerateTime(0)
	}

	if payload.MaxDate == "" {
		payload.MaxDate = helper.GenerateTime(24)
	}

	transactionData, err := pr.transactionRepoReads.GetTransactionDataByProviderChannelRepo(payload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	dataAnalytics := supportMerchantAnalyticsByMerchantPaychannelSvc(transactionData)

	providerChannelData, err := pr.providerRepoReads.GetDetailProviderChannelById(payload.ProviderChannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if providerChannelData.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "data not found",
		}
		return resp, errors.New("data not found")
	}

	analyticsResp := dto.ProviderChannelAnalyticsResDto{
		AnalyticsData:         dataAnalytics,
		ProviderChannelDetail: providerChannelData,
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            analyticsResp,
	}

	return resp, nil
}

func (pr *Provider) UpdateFeeLimitInterfaceProviderChannelSvc(payload dto.AdjustLimitOrFeeProviderPayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	providerPaychannelData, err := pr.providerRepoReads.GetDetailProviderChannelById(payload.ProviderChannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// Use existing values if the new values are not provided
	if payload.MaxAmount == nil {
		payload.MaxAmount = &providerPaychannelData.MaxAmount
	}

	if payload.MinAmount == nil {
		payload.MinAmount = &providerPaychannelData.MinAmount
	}

	if payload.MaxDailyLimit == nil {
		payload.MaxDailyLimit = &providerPaychannelData.MaxDailyLimit
	}

	if payload.Fee == nil {
		payload.Fee = &providerPaychannelData.Fee
	}

	if payload.FeeType == nil {
		payload.FeeType = &providerPaychannelData.FeeType
	}

	if payload.InterfaceSetting == nil {
		payload.InterfaceSetting = &providerPaychannelData.InterfaceSetting
	}

	// update provider channel
	err = pr.providerRepoWrites.UpdateProviderPaychannelByIdRepo(payload)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	providerChannelUpdated, err := pr.providerRepoReads.GetDetailProviderChannelById(payload.ProviderChannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: fmt.Sprintf("Successfully updated merchant pay channel with id: %v", payload.ProviderChannelId),
		Data:            providerChannelUpdated,
	}

	return resp, nil
}

func (pr *Provider) GetProviderChannelOperatorSvc(providerChannelId int) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	bankList, err := pr.providerRepoReads.GetBankListProviderMethodRepo(providerChannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	bankNameListPaychannel, err := pr.providerRepoReads.GetBankListProviderChannelRepo(providerChannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	for j := range bankList {
		for i := range bankNameListPaychannel {
			if bankNameListPaychannel[i].BankName == bankList[j].BankName {
				bankList[j].CheckListFlagging = true
			}
		}
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            bankList,
	}

	return resp, nil
}

func (pr *Provider) GetListRoutedProviderChannelSvc(providerChannelId int) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listRoutedProviderChannel, err := pr.providerRepoReads.GetListRoutedProviderChannelRepo(providerChannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listRoutedProviderChannel) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "data not found",
		}
		return resp, nil
	}

	for i := range listRoutedProviderChannel {
		availableChannel, err := pr.providerRepoReads.CountProviderChannelByPaymentMethodRepo(listRoutedProviderChannel[i].PaymentMethodName)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}
		availableChannelStr := converter.ToString(availableChannel)

		activeChannel, err := pr.providerRepoReads.CountActiveProviderChannelRepo(listRoutedProviderChannel[i].Id)
		if err != nil {
			resp = dto.ResponseDto{
				ResponseCode:    http.StatusUnprocessableEntity,
				ResponseMessage: err.Error(),
			}
			return resp, err
		}
		activeChannelStr := converter.ToString(activeChannel)
		activeAvailableChannel := activeChannelStr + "/" + availableChannelStr
		listRoutedProviderChannel[i].ActiveAvailableChannels = activeAvailableChannel
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            listRoutedProviderChannel,
	}

	return resp, nil
}

func (pr *Provider) AddOperatorProviderChannelSvc(payload []dto.AddOperatorProviderChannelPayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	for _, load := range payload {
		if load.CheckListFlagging {
			bankDetailData, err := pr.transactionRepoReads.GetBankDataDetailByBankCodeRepo(load.BankCode)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			if bankDetailData.Id == 0 {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: "data not found, maybe wrong bank code",
				}
				return resp, errors.New("data not found, maybe wrong bank code")
			}

			// checking if it exist to avoid double input
			operatorRes, err := pr.providerRepoReads.GetProviderBankListChannelRepo(load.ProviderChannelId, bankDetailData.Id)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			if operatorRes.Id == 0 {
				//create provider channel operator
				_, err = pr.providerRepoWrites.AddOperatorProviderChannelRepo(load.ProviderChannelId, bankDetailData.Id)
				if err != nil {
					resp = dto.ResponseDto{
						ResponseCode:    http.StatusUnprocessableEntity,
						ResponseMessage: err.Error(),
					}
					return resp, err
				}
			}
		} else {
			bankDetailData, err := pr.transactionRepoReads.GetBankDataDetailByBankCodeRepo(load.BankCode)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			// checking if it exist to avoid double input
			operatorRes, err := pr.providerRepoReads.GetProviderBankListChannelRepo(load.ProviderChannelId, bankDetailData.Id)
			if err != nil {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: err.Error(),
				}
				return resp, err
			}

			if operatorRes.Id != 0 {
				// deleted if false
				err = pr.providerRepoWrites.DeleteOperatorProviderChannelRepo(load.ProviderChannelId, bankDetailData.Id)
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

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success add/remove bank channel",
	}

	return resp, nil
}

func (pr *Provider) ActiveOrDeactivateProviderPaychannelIdSvc(payload dto.UpdateStatusProviderPaychannelDto) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	status := constant.StatusActive

	providerPaychannelDetail, err := pr.providerRepoReads.GetDetailProviderChannelById(payload.ProviderChannelId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if providerPaychannelDetail.Id == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "wrong provider channel id",
		}
		return resp, errors.New("wrong provider channel id")
	}

	if providerPaychannelDetail.Status == constant.StatusActive {
		status = constant.StatusInactive
	}

	// check if fee already set
	if providerPaychannelDetail.Fee == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "can't activate need have to set fee first",
		}
		return resp, errors.New("need to set fee")
	}

	// update status provider paychannel
	err = pr.providerRepoWrites.UpdateStatusProviderPaychannelRepo(payload.ProviderChannelId, status)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: fmt.Sprintf("success update status for provider paychannel id %v", payload.ProviderChannelId),
	}

	return resp, nil
}

func (pr *Provider) GetListInterfaceProviderSvc(params dto.QueryParams) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// get list interface
	listInterface, err := pr.providerRepoReads.GetProviderInterfaceWithFilterRepo(params)
	if err != nil {
		slog.Errorw("listInterface failed", "stack_trace", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "failed to get interface list",
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success get list interface code",
		Data:            listInterface,
	}

	return resp, nil
}

func (pr *Provider) GetListPaymentOperatorCreateChannelProviderSvc(providerPaymentMethodId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	interfaceIdInt := converter.ToInt(providerPaymentMethodId)

	// operator list
	operatorList, err := pr.providerRepoReads.GetBankListProviderInterfaceRepo(interfaceIdInt)
	if err != nil {
		slog.Errorw("failed get operator list", "stack_trace", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "failed get operator list",
		}
		return resp, err
	}

	if len(operatorList) == 0 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "this provider interface doesn't have payment operator",
		}
		return resp, errors.New("data not found")
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success get operator list",
		Data:            operatorList,
	}

	return resp, nil
}

func (pr *Provider) CreateProviderChannelSvc(payload dto.CreateProviderChannelDto) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	// Check and provide default value if nil
	if payload.InterfaceSetting == nil {
		defaultSetting := "MAIN"
		payload.InterfaceSetting = &defaultSetting
	}

	if payload.MinAmount == nil {
		defaultMinAmount := 0.0
		payload.MinAmount = &defaultMinAmount
	}

	if payload.MaxAmount == nil {
		defaultMaxAmount := 0.0
		payload.MaxAmount = &defaultMaxAmount
	}

	if payload.DailyLimit == nil {
		defaultDailyLimit := 0.0
		payload.DailyLimit = &defaultDailyLimit
	}

	if payload.Fee == nil {
		defaultFee := 0.0
		payload.Fee = &defaultFee
	}

	if payload.FeeType == nil {
		defaultFeeType := constant.FeeTypeFixedFee
		payload.FeeType = &defaultFeeType
	}

	// create provider paychannel
	paychannelId, err := pr.providerRepoWrites.CreateProviderPaychannelRepo(payload)
	if err != nil {
		slog.Errorw("create paychannel failed", "stack_trace", err.Error())
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: "failed create paychannel",
		}
		return resp, err
	}

	if len(payload.BankOperator) > 0 {
		for _, bank := range payload.BankOperator {
			bankDetailData, err := pr.transactionRepoReads.GetBankDataDetailByBankCodeRepo(bank.BankCode)
			if err != nil {
				slog.Errorw("bankDetailData", "stack_trace", err.Error())
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: constant.GeneralErrMsg,
				}
				return resp, err
			}

			if bankDetailData.Id == 0 {
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: "data not found, maybe wrong bank code",
				}
				return resp, errors.New("data not found, maybe wrong bank code")
			}

			_, err = pr.providerRepoWrites.AddOperatorProviderChannelRepo(paychannelId, bankDetailData.Id)
			if err != nil {
				slog.Errorw("AddOperatorProviderChannelRepo", "stack_trace", err.Error())
				resp = dto.ResponseDto{
					ResponseCode:    http.StatusUnprocessableEntity,
					ResponseMessage: constant.GeneralErrMsg,
				}
				return resp, err
			}
		}
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: fmt.Sprintf("success create provider paychannel with id %v", paychannelId),
	}

	return resp, nil
}

func supportProviderAnalyticsSvc(payload []entity.PaymentDetailMerchantProvider) dto.AnalyticsProviderRespDto {
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

	providerAnalyticsDataRes := dto.AnalyticsProviderRespDto{
		TransactionIn:  inAnalyticsData,
		TransactionOut: outAnalyticsData,
	}

	return providerAnalyticsDataRes
}
