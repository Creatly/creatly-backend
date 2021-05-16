package email

type AddEmailInput struct {
	Email     string
	ListID    string
	Variables map[string]string
}

type Provider interface {
	AddEmailToList(AddEmailInput) error
}
