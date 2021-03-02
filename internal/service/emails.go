package service

import (
	"fmt"
	emailProvider "github.com/zhashkevych/courses-backend/pkg/email"
)

const (
	nameField            = "name"
	verificationLinkTmpl = "%s/verification?code=%s" // <frontend_url>/verification?code=<verification_code>

	// TODO: manage email templates. Store in config
	verificationEmailSubject = "Спасибо за регистрацию, %s!"
	verificationEmailTmpl    = `<h1>Спасибо за регистрацию!</h1><br>Чтобы подтвердить свой аккаунт, <a href="%s">переходи по ссылке</a>`
)

type EmailService struct {
	provider    emailProvider.Provider
	sender      emailProvider.Sender
	listId      string
	frontendUrl string
}

func NewEmailsService(provider emailProvider.Provider, sender emailProvider.Sender, listId, frontendUrl string) *EmailService {
	return &EmailService{provider: provider, sender: sender, listId: listId, frontendUrl: frontendUrl}
}

func (s *EmailService) AddToList(name, email string) error {
	return s.provider.AddEmailToList(emailProvider.AddEmailInput{
		Email:  email,
		ListID: s.listId,
		Variables: map[string]string{
			nameField: name,
		},
	})
}

func (s *EmailService) SendVerificationEmail(input SendVerificationEmailInput) error {
	subject := fmt.Sprintf(verificationEmailSubject, input.Name)
	body := fmt.Sprintf(verificationEmailTmpl, s.createVerificationLink(input.VerificationCode))

	return s.sender.Send(emailProvider.SendEmailInput{
		To:      input.Email,
		Subject: subject,
		Body:    body,
	})
}

func (s *EmailService) createVerificationLink(code string) string {
	return fmt.Sprintf(verificationLinkTmpl, s.frontendUrl, code)
}
