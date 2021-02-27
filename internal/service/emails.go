package service

import (
	"fmt"
	"github.com/zhashkevych/courses-backend/pkg/email"
)

const (
	nameField              = "name"
	verificationLinklField = "verificationLink"
	verificationLinkTmpl   = "%s/verification?code=%s" // <frontend_url>/verification?code=<verification_code>
)

type EmailService struct {
	provider    email.Provider
	listId      string
	frontendUrl string
}

func NewEmailsService(provider email.Provider, listId, frontendUrl string) *EmailService {
	return &EmailService{provider: provider, listId: listId, frontendUrl: frontendUrl}
}

func (s *EmailService) AddToList(input AddToListInput) error {
	return s.provider.AddEmailToList(email.AddEmailInput{
		Email:  input.Email,
		ListID: s.listId,
		Variables: map[string]string{
			nameField:              input.Name,
			verificationLinklField: s.createVerificationLink(input.VerificationCode),
		},
	})
}

func (s *EmailService) createVerificationLink(code string) string {
	return fmt.Sprintf(verificationLinkTmpl, s.frontendUrl, code)
}
