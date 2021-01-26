package fondy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"io/ioutil"
	"net/http"
)

type Client struct {
	merchantId       string
	merchantPassword string
}

func NewClient(merchantId string, merchantPassword string) *Client {
	return &Client{merchantId: merchantId, merchantPassword: merchantPassword}
}

func (c *Client) GeneratePaymentLink(input payment.GeneratePaymentLinkInput) (string, error) {
	checkoutReq := &checkoutRequest{
		OrderId:           input.OrderId,
		MerchantId:        c.merchantId,
		OrderDesc:         input.OrderDesc,
		Amount:            fmt.Sprintf("%d", input.Amount),
		Currency:          input.Currency,
		ServerCallbackURL: input.CallbackURL,
		ResponseURL:       input.ResponseURL,
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

	return apiResp.Response.CheckoutURL, nil
}
