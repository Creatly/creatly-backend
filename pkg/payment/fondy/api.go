package fondy

import (
	"crypto/sha1"
	"fmt"
	"github.com/fatih/structs"
	"sort"
	"strings"
)

// https://docs.fondy.eu/ru/docs/page/1/

// success card - 4444555566661111
// failure card - 4444111166665555

const (
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

type callbackResponse struct {
	OrderId        string `json:"order_id"`
	MerchantId     string `json:"merchant_id"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Signature      string `json:"signature"`
	OrderStatus    string `json:"order_status"`    // created; processing; declined; approved; expired; reversed;
	ResponseStatus string `json:"response_status"` // 1) success; 2) failure
	MaskedCard     string `json:"masked_card"`
	CardType       string `json:"card_type"`
	Fee            string `json:"fee"`
	PaymentSystem  string `json:"payment_system"`
	ProductId      string `json:"product_id"`
	AdditionalInfo string `json:"additional_info"`
}

func (r *checkoutRequest) setSignature(password string) {
	params := structs.Map(r)

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	values := []string{}

	for _, key := range keys {
		value := params[key].(string)
		if value == "" {
			continue
		}

		values = append(values, value)
	}

	r.Signature = generateSignature(values, password)
}

func generateSignature(values []string, password string) string {
	newValues := []string{password}
	newValues = append(newValues, values...)

	signatureString := strings.Join(newValues, "|")

	fmt.Println(signatureString)

	hash := sha1.New()
	hash.Write([]byte(signatureString))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
