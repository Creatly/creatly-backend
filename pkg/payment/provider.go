package payment

type GeneratePaymentLinkInput struct {
	OrderId     string
	Amount      int
	Currency    string
	OrderDesc   string
	CallbackURL string
	ResponseURL string
}

type FondyProvider interface {
	GeneratePaymentLink(input GeneratePaymentLinkInput) (string, error)
	ValidateCallback(input Callback) error
}
