package fondy

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/zhashkevych/creatly-backend/pkg/logger"

	"github.com/fatih/structs"
	"github.com/zhashkevych/creatly-backend/pkg/payment"
)

// Documentation
// https://docs.fondy.eu/ru/docs/page/1/

// Testing credentials
// success card - 4444555566661111
// failure card - 4444111166665555

const (
	// UserAgent is a value for user-agent header sent in Fondy's requests. Used to validate request.
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64; Twisted) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"

	checkoutUrl   = "https://pay.fondy.eu/api/checkout/url/"
	languageRU    = "ru"
	statusSuccess = "success"
)

type apiRequest struct {
	Request interface{} `json:"request"`
}

type apiResponse struct {
	Response interimResponse `json:"response"`
}

type checkoutRequest struct {
	OrderId           string `json:"order_id"`
	MerchantId        string `json:"merchant_id"`
	OrderDesc         string `json:"order_desc"`
	Signature         string `json:"signature"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	ResponseURL       string `json:"response_url,omitempty"`
	ServerCallbackURL string `json:"server_callback_url,omitempty"`
	SenderEmail       string `json:"sender_email,omitempty"`
	Language          string `json:"lang,omitempty"`
	ProductId         string `json:"product_id,omitempty"`
}

type interimResponse struct {
	Status       string `json:"response_status"`
	CheckoutURL  string `json:"checkout_url"`
	PaymentId    string `json:"payment_id"`
	ErrorMessage string `json:"error_message"`
}

type Callback struct {
	OrderId                 string      `json:"order_id"`
	MerchantId              int         `json:"merchant_id"`
	Amount                  string      `json:"amount"`
	Currency                string      `json:"currency"`
	OrderStatus             string      `json:"order_status"`    // created; processing; declined; approved; expired; reversed;
	ResponseStatus          string      `json:"response_status"` // 1) success; 2) failure
	Signature               string      `json:"signature"`
	TranType                string      `json:"tran_type"`
	SenderCellPhone         string      `json:"sender_cell_phone"`
	SenderAccount           string      `json:"sender_account"`
	CardBin                 int         `json:"card_bin"`
	MaskedCard              string      `json:"masked_card"`
	CardType                string      `json:"card_type"`
	RRN                     string      `json:"rrn"`
	ApprovalCode            string      `json:"approval_code"`
	ResponseCode            interface{} `json:"response_code"`
	ResponseDescription     string      `json:"response_description"`
	ReversalAmount          string      `json:"reversal_amount"`
	SettlementAmount        string      `json:"settlement_amount"`
	SettlementCurrency      string      `json:"settlement_currency"`
	OrderTime               string      `json:"order_time"`
	SettlementDate          string      `json:"settlement_date"`
	ECI                     string      `json:"eci"`
	Fee                     string      `json:"fee"`
	PaymentSystem           string      `json:"payment_system"`
	SenderEmail             string      `json:"sender_email"`
	PaymentId               int         `json:"payment_id"`
	ActualAmount            string      `json:"actual_amount"`
	ActualCurrency          string      `json:"actual_currency"`
	MerchantData            string      `json:"merchant_data"`
	VerificationStatus      string      `json:"verification_status"`
	Rectoken                string      `json:"rectoken"`
	RectokenLifetime        string      `json:"rectoken_lifetime"`
	ProductId               string      `json:"product_id"`
	AdditionalInfo          string      `json:"additional_info"`
	ResponseSignatureString string      `json:"response_signature_string"`
}

func (c Callback) Success() bool {
	return c.ResponseStatus == "success"
}

func (c Callback) PaymentApproved() bool {
	return c.OrderStatus == "approved"
}

func (r *checkoutRequest) setSignature(password string) {
	params := structs.Map(r)
	r.Signature = generateSignature(params, password)
}

//nolint:unused
func (c *Callback) validateSignature(password string) bool {
	params := structs.Map(c)

	logger.Debugf("[FONDY] callback: %+v", c)

	delete(params, "Signature")
	delete(params, "ResponseSignatureString")

	signature := generateSignature(params, password)

	logger.Debugf("[FONDY] generated signature: %s", signature)

	return c.Signature == signature
}

func generateSignature(params map[string]interface{}, password string) string {
	keys := make([]string, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	values := []string{}

	for _, key := range keys {
		value, ok := params[key].(string)
		if !ok && params[key] != nil {
			values = append(values, fmt.Sprintf("%v", params[key]))

			continue
		}

		if value == "" {
			continue
		}

		values = append(values, value)
	}

	newValues := []string{password}
	newValues = append(newValues, values...)

	signatureString := strings.Join(newValues, "|")

	logger.Infof("[FONDY] generated signatureString %s", signatureString)

	hash := sha1.New()
	hash.Write([]byte(signatureString)) //nolint:errcheck

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Client is a fondy payment provider API client.
type Client struct {
	merchantID       string
	merchantPassword string
}

func NewFondyClient(merchantID string, merchantPassword string) *Client {
	return &Client{merchantID: merchantID, merchantPassword: merchantPassword}
}

// GeneratePaymentLink returns payment URL for provided order info.
func (c *Client) GeneratePaymentLink(input payment.GeneratePaymentLinkInput) (string, error) {
	checkoutReq := &checkoutRequest{
		OrderId:           input.OrderId,
		MerchantId:        c.merchantID,
		OrderDesc:         input.OrderDesc,
		Amount:            fmt.Sprintf("%d", input.Amount),
		Currency:          input.Currency,
		ServerCallbackURL: input.CallbackURL,
		ResponseURL:       input.RedirectURL,
		Language:          languageRU,
	}

	checkoutReq.setSignature(c.merchantPassword)

	request := apiRequest{Request: checkoutReq}
	requestBody, _ := json.Marshal(request)

	resp, err := http.Post(checkoutUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	apiResp := apiResponse{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return "", err
	}

	if apiResp.Response.Status == statusSuccess {
		return apiResp.Response.CheckoutURL, nil
	}

	return "", errors.New(apiResp.Response.ErrorMessage)
}

func (c *Client) ValidateCallback(input interface{}) error {
	_, ok := input.(Callback)
	if !ok {
		return errors.New("invalid callback data")
	}

	// if !callback.validateSignature(c.merchantPassword) {
	//	return errors.New("invalid signature")
	// }

	return nil
}
