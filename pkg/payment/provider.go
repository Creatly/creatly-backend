package payment

import "github.com/zhashkevych/courses-backend/pkg/logger"

type GeneratePaymentLinkInput struct {
	OrderId string
	Amount int
	Currency string
	OrderDesc string
	CallbackURL string
	ResponseURL string
}

type Provider interface {
	GeneratePaymentLink(input GeneratePaymentLinkInput) (string, error)
}

type MockProvider struct {

}

func (p MockProvider) GeneratePaymentLink(input GeneratePaymentLinkInput) (string, error) {
	logger.Info(input)
	return "paymentLink", nil
}