package jack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
)

type jack struct {
	configApp  config.App
	httpClient *http.Client
}

func New(cfg config.App) *jack {
	return &jack{
		configApp: cfg,
		httpClient: &http.Client{
			Timeout: constant.NinetySecond,
		},
	}
}

func (jk *jack) InquiryAccount(payload dto.MerchantDisbursement, credentials []entity.ProviderCredentialsEntity, bankCode string) (dto.InquiryAccountResponse, error) {
	var resp dto.InquiryAccountResponse
	bankName := constant.TransformJackBank[bankCode].BankName
	var cfg dto.JackCredentialsDto

	for _, cred := range credentials {
		if cred.Key == constant.JackApiKeyCred {
			cfg.ApiKey = cred.Value
		}

		if cred.Key == constant.JackInquiryKeyUrlCred {
			cfg.InquiryUrl = cred.Value
		}
	}

	slog.Infof("JACK %v [inquiry-account] with account number: %v, bank name: %v", payload.Username, payload.BankAccountNumber, payload.BankName)

	// http request
	r, err := http.NewRequest(http.MethodGet, cfg.InquiryUrl, nil)
	if err != nil {
		return resp, errors.New("failed to create request")
	}

	// set header
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", cfg.ApiKey)

	// set query params
	q := r.URL.Query()
	q.Add("bank_name", bankName)
	q.Add("account_number", payload.BankAccountNumber)
	r.URL.RawQuery = q.Encode()

	r.Close = true
	response, err := jk.httpClient.Do(r)
	if err != nil {
		return resp, errors.New("failed to send request")
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
		slog.Infof("JACK %v [inquiry-account] got error failed to read response %v", payload.Username, string(contents))
		return resp, errors.New("failed to read response body")
	}

	if response.StatusCode != http.StatusOK {
		slog.Infof("JACK %v [inquiry-account] got error response status not ok: %v", payload.Username, string(contents))
		return resp, errors.New(string(contents))
	}

	// conver response to struct
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		slog.Infof("JACK %v [inquiry-account] error failed to unmarshall response %v", payload.Username, string(contents))
		return resp, errors.New("failed to unmarshall response to struct")
	}

	// validate status on response payload
	if resp.Status != constant.JackStatusOk {
		if resp.Status == constant.JackStatusInvalid {
			msg := resp.Data.Errors
			return resp, errors.New(msg)
		}

		slog.Infof("JACK %v [inquiry-account] got error response status not ok: %v", payload.Username, converter.ToString(resp))
		return resp, errors.New(converter.ToString(resp))
	}

	slog.Infof("JACK %v [inquiry-account] response data: %v", payload.Username, converter.ToString(resp))

	return resp, nil
}

func (jk *jack) GetBalance(username string, credentials []entity.ProviderCredentialsEntity) (int, error) {
	var cfg dto.JackCredentialsDto
	var resp dto.JackGetBalanceResponse

	for _, cred := range credentials {
		if cred.Key == constant.JackApiKeyCred {
			cfg.ApiKey = cred.Value
		}

		if cred.Key == constant.JackGetBalanceKeyUrlCred {
			cfg.GetBalanceUrl = cred.Value
		}
	}

	// create request
	r, err := http.NewRequest(http.MethodGet, cfg.GetBalanceUrl, nil)
	if err != nil {
		return 0, errors.New("failed to create request")
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", cfg.ApiKey)
	r.Close = true

	response, err := jk.httpClient.Do(r)
	if err != nil {
		return 0, errors.New("failed to send request")
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
		slog.Infof("JACK %v [get-balance] got error failed to read response %v", username, string(contents))
		return 0, errors.New("failed to read response body")
	}

	if response.StatusCode != http.StatusOK {
		slog.Infof("JACK %v [get-balance] got error response status not ok", username)
		return 0, errors.New(string(contents))
	}

	// convert response to struct
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		slog.Infof("JACK %v [get-balance] error failed to unmarshall response %v", username, string(contents))
		return 0, errors.New("failed to unmarshall response to struct")
	}

	if resp.Status != constant.JackStatusOk {
		slog.Infof("JACK %v [get-balance] got error response status not ok", username)
		return 0, errors.New(converter.ToString(resp))
	}

	slog.Infof("JACK %v [get-balance]: %v", username, converter.ToString(resp))

	return resp.Data.Balances[2].Balance, nil
}

func (jk *jack) ConfirmDisbursement(payload dto.ConfirmTransactionPayload, credentials []entity.ProviderCredentialsEntity) (dto.CreateDisbursementRequestResponse, error) {
	var resp dto.CreateDisbursementRequestResponse
	var cfg dto.JackCredentialsDto

	for _, cred := range credentials {
		if cred.Key == constant.JackApiKeyCred {
			cfg.ApiKey = cred.Value
		}

		if cred.Key == constant.JackDisbursementUrlCred {
			cfg.DisbursementUrl = cred.Value
		}
	}

	// URL path param join
	confirmTransactionURL, _ := url.JoinPath(cfg.DisbursementUrl, url.PathEscape(payload.ProviderID), "confirm")

	// http request
	r, err := http.NewRequest(http.MethodPost, confirmTransactionURL, nil)
	if err != nil {
		return resp, errors.New("failed to create request")
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", cfg.ApiKey)
	r.Close = true

	response, err := jk.httpClient.Do(r)
	if err != nil {
		return resp, errors.New("failed to send request")
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
		slog.Infof("JACK %v [confirm-disbursement] got error failed to read response %v", payload.Username, string(contents))
		return resp, errors.New("failed to read response body")
	}

	if response.StatusCode != http.StatusOK {
		slog.Infof("JACK %v [confirm-disbursement] got error response status not ok", payload.Username)
		return resp, errors.New(string(contents))
	}

	// convert response to struct
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		slog.Infof("JACK %v [confirm-disbursement] error failed to unmarshall response %v", payload.Username, string(contents))
		return resp, errors.New("failed to unmarshall response to struct")
	}

	if resp.Status != constant.JackStatusOk {
		slog.Infof("JACK %v [confirm-disbursement] response data with status not ok: %v", payload.Username, converter.ToString(resp))
		return resp, errors.New(resp.Data.ErrorMessage)
	}

	if resp.Data.State != constant.JackStateStatusConfirm {
		slog.Infof("JACK %v [confirm-disbursement] response data with state not confirmed: %v", payload.Username, converter.ToString(resp))
		return resp, errors.New("status not confirmed")
	}

	slog.Infof("JACK %v [confirm-disbursement] response data: %v", payload.Username, converter.ToString(resp))

	return resp, nil
}

func (jk *jack) CreateDisbursement(payload dto.MerchantDisbursement, credentials []entity.ProviderCredentialsEntity, bankCode string, paymentId string) (dto.CreateDisbursementRequestResponse, error) {
	var resp dto.CreateDisbursementRequestResponse
	var cfg dto.JackCredentialsDto

	for _, cred := range credentials {
		if cred.Key == constant.JackApiKeyCred {
			cfg.ApiKey = cred.Value
		}

		if cred.Key == constant.JackDisbursementUrlCred {
			cfg.DisbursementUrl = cred.Value
		}
	}

	// first and last beneficiary account name
	firstAccountName, lastAccountName, err := helper.CheckingFirstAndLastStr(payload.BankAccountName)
	if err != nil {
		slog.Infof("JACK %v [create-disbursement] got failed there is no account name", payload.Username)
		return resp, err
	}
	// create disbursement payload request
	senderData := dto.SenderData{
		FirstName:      constant.JackDisbursementFirstSenderName,
		LastName:       constant.JackDisbursementSecondSenderName,
		CountryIsoCode: constant.JackDisbursementCountryIsoName,
	}

	destinationData := dto.DestinationData{
		Amount:         converter.ToString(payload.Amount),
		Currency:       constant.JackDisbursementCurrency,
		CountryIsoCode: constant.JackDisbursementCountryIsoName,
	}

	sourceData := dto.SourceData{
		Currency:       constant.JackDisbursementCurrency,
		CountryIsoCode: constant.JackDisbursementCountryIsoName,
	}

	beneficiaryData := dto.BeneficiaryData{
		FirstName:      firstAccountName,
		LastName:       lastAccountName,
		CountryIsoCode: constant.JackDisbursementCountryIsoName,
		Account:        payload.BankAccountNumber,
	}

	notes := constant.JackDisbursementNotes + " - " + payload.BankAccountName

	requestData := dto.CreateDisbursementRequest{
		ReferenceID: paymentId,
		CallbackURL: jk.configApp.CallbackUrl,
		PayerID:     constant.TransformJackBank[bankCode].Id,
		Mode:        constant.JackDisbursementMode,
		Source:      sourceData,
		Sender:      senderData,
		Destination: destinationData,
		Beneficiary: beneficiaryData,
		Notes:       notes,
	}

	// transform to JSON
	payloadJSON, _ := json.Marshal(requestData)
	slog.Infof("JACK %v [create-disbursement] payload data: %v", payload.Username, string(payloadJSON))

	// http request
	r, err := http.NewRequest(http.MethodPost, cfg.DisbursementUrl, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return resp, errors.New("failed to create request")
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", cfg.ApiKey)
	r.Close = true

	response, err := jk.httpClient.Do(r)
	if err != nil {
		return resp, errors.New("failed to send request")
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
		slog.Infof("JACK %v [create-disbursement] got error failed to read response %v", payload.Username, string(contents))
		return resp, errors.New("failed to read response body")
	}

	if response.StatusCode != http.StatusOK {
		slog.Infof("JACK %v [create-disbursement] got error response status not ok: %v", payload.Username, string(contents))
		return resp, errors.New(string(contents))
	}

	// convert response to struct
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		slog.Infof("JACK %v [create-disbursement] error failed to unmarshall response %v", payload.Username, string(contents))
		return resp, errors.New("failed to unmarshall response to struct")
	}

	if resp.Status != constant.JackStatusOk {
		slog.Infof("JACK %v [create-disbursement] got error response status not ok: %v", payload.Username, converter.ToString(resp))
		return resp, errors.New(resp.Data.ErrorMessage)
	}

	slog.Infof("JACK %v [create-disbursement] response data: %v", payload.Username, converter.ToString(resp))

	return resp, nil
}
