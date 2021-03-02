package email

type SendEmailInput struct {
	To      string
	Subject string
	Body    string
}

type Sender interface {
	Send(input SendEmailInput) error
}
