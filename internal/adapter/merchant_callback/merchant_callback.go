package merchantcallback

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
)

type merchantCallback struct {
	configApp  config.App
	httpClient *http.Client
}

func New(cfg config.App) *merchantCallback {
	return &merchantCallback{
		configApp: cfg,
		httpClient: &http.Client{
			Timeout: constant.NinetySecond,
		},
	}
}

func (mc *merchantCallback) SendCallbackAdptr(url string, transactionEntity entity.PaymentDetailMerchantProvider, transactionStatusLogLatest entity.TransactionStatusLogs, merchantSecret string) (interface{}, error) {
	var merchantResponse interface{}
	amountFormatted := helper.FormatFloat64(transactionEntity.TransactionAmount)

	// request data to merchant
	requestData := dto.MerchantCallbackDto{
		TransactionId:         transactionEntity.PaymentID,
		MerchantTransactionId: transactionEntity.MerchantRefNumber,
		Status:                transactionEntity.Status,
		Amount:                amountFormatted,
		TransactionType:       transactionEntity.PaymentMethodName,
		TransactionCreatedAt:  transactionEntity.TransactionCreatedAt,
		TransactionUpdatedAt:  transactionEntity.TransactionUpdatedAt,
	}

	if transactionEntity.Status == constant.StatusFailed {
		requestData.FailedReason = *transactionStatusLogLatest.Notes
	}

	// payload json
	payloadJson, _ := json.Marshal(requestData)

	signature := helper.StringToSignatureSymmetric(string(payloadJson), merchantSecret)

	r, err := http.NewRequest(http.MethodPost, transactionEntity.MerchantCallbackURL, bytes.NewBuffer(payloadJson))
	if err != nil {
		return merchantResponse, err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-signature", signature)
	r.Close = true

	response, err := mc.httpClient.Do(r)
	if err != nil {
		merchantResponse = "500:" + err.Error()
		return merchantResponse, err
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Println("failed to close response body, could lead to memory leak")
		}
	}()

	// read response body
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		merchantResponse = "500:" + string(contents)
		return merchantResponse, err
	}

	if response.StatusCode != http.StatusOK {
		merchantResponse = converter.ToString(response.StatusCode) + ":" + string(contents)
		return merchantResponse, errors.New("status not ok")
	}

	merchantResponse = converter.ToString(response.StatusCode) + ":" + string(contents)

	return merchantResponse, nil
}
