package service

import "github.com/zhashkevych/courses-backend/pkg/email"

const (
	nameField             = "name"
	registerSourceField   = "registerSource"
	verificationCodeField = "verificationCode"
)

type EmailService struct {
	provider email.Provider
	listId   string
}

func NewEmailsService(provider email.Provider, listId string) *EmailService {
	return &EmailService{provider: provider, listId: listId}
}

func (s *EmailService) AddToList(input AddToListInput) error {
	return s.provider.AddEmailToList(email.AddEmailInput{
		Email:  input.Email,
		ListID: s.listId,
		Variables: map[string]string{
			nameField:             input.Name,
			registerSourceField:   input.RegisterSource,
			verificationCodeField: input.VerificationCode,
		},
	})
}
