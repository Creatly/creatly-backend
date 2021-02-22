package payment

type GeneratePaymentLinkInput struct {
	OrderId     string
	Amount      uint
	Currency    string
	OrderDesc   string
	CallbackURL string
	ResponseURL string
}

type Provider interface {
	GeneratePaymentLink(input GeneratePaymentLinkInput) (string, error)
	ValidateCallback(input interface{}) error
}
